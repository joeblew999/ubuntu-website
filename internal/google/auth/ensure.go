package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// A credential is checked where it's used, with Google, not by looking at a file: a command that needs one tries it
// first, and when it's missing, expired or revoked, signs in again in the browser and carries on. Clicking Allow is
// the one thing left to a person.

// tokenURL is Google's token endpoint (a var so tests can point it at their own server).
var tokenURL = OAuthTokenURL

// EnsureGcloudCredentials makes sure the application-default credentials (Terraform's, and any Google SDK's) work:
// it refreshes them at Google, and when they're missing, expired or revoked (invalid_grant), signs in again in the
// browser and writes them.
func EnsureGcloudCredentials(w io.Writer) error {
	ok, credFile, err := CheckGcloudCredentials()
	if err == nil && ok {
		data, _ := os.ReadFile(credFile)
		var creds GcloudCredentials
		json.Unmarshal(data, &creds)
		switch err := refresh(creds.ClientID, creds.ClientSecret, creds.RefreshToken); {
		case err == nil:
			return nil
		case !errors.Is(err, ErrSignInAgain):
			return fmt.Errorf("application-default credentials: %w", err)
		}
		fmt.Fprintln(w, "The application-default credentials (Terraform's) have expired: sign in again in the browser.")
	} else {
		fmt.Fprintln(w, "No application-default credentials (Terraform's): sign in in the browser.")
	}
	result, err := RunOAuthFlowAssisted(FlowOptions{Timeout: 3 * time.Minute})
	if err != nil {
		return err
	}
	if err := WriteGcloudCredentials(result.RefreshToken); err != nil {
		return err
	}
	fmt.Fprintln(w, "✓ application-default credentials written")
	return nil
}

// ErrSignInAgain is Google refusing a refresh token (invalid_grant: expired or revoked): only signing in again helps.
var ErrSignInAgain = errors.New("the sign-in has expired or been revoked")

// refresh tries a refresh token at Google.
func refresh(clientID, clientSecret, refreshToken string) error {
	resp, err := http.PostForm(tokenURL, url.Values{
		"grant_type": {"refresh_token"}, "refresh_token": {refreshToken}, "client_id": {clientID}, "client_secret": {clientSecret},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	var body TokenResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Error == "invalid_grant" {
		return ErrSignInAgain
	}
	return fmt.Errorf("Google's token endpoint answered HTTP %d (%s)", resp.StatusCode, body.Error)
}
