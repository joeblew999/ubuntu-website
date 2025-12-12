// Standalone Google OAuth helper for google-mcp-server
// This tool authenticates and saves the token file that google-mcp-server expects.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Println("Error: GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set")
		fmt.Println("")
		fmt.Println("Run: source .env")
		os.Exit(1)
	}

	scopes := []string{
		"https://www.googleapis.com/auth/calendar",
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/gmail.modify",
		"https://www.googleapis.com/auth/spreadsheets",
		"https://www.googleapis.com/auth/documents",
		"https://www.googleapis.com/auth/presentations",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}

	// Generate auth URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("Opening browser for Google authentication...")
	fmt.Println("")
	fmt.Println("If browser doesn't open, visit:")
	fmt.Println(authURL)
	fmt.Println("")

	// Open browser
	openBrowser(authURL)

	// Start callback server
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 10 * time.Second,
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errChan <- fmt.Errorf("no authorization code received")
			http.Error(w, "No authorization code", http.StatusBadRequest)
			return
		}
		codeChan <- code
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><body>
			<h1>✅ Authentication successful!</h1>
			<p>You can close this window and return to the terminal.</p>
		</body></html>`)
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	fmt.Println("Waiting for authorization (timeout: 5 minutes)...")

	// Wait for code
	var code string
	select {
	case code = <-codeChan:
		fmt.Println("Authorization code received!")
	case err := <-errChan:
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	case <-time.After(5 * time.Minute):
		fmt.Println("Timeout waiting for authorization")
		os.Exit(1)
	}

	// Shutdown server
	ctx := context.Background()
	server.Shutdown(ctx)

	// Exchange code for token
	fmt.Println("Exchanging code for token...")
	token, err := config.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("Failed to exchange code: %v\n", err)
		os.Exit(1)
	}

	// Save token to file
	homeDir, _ := os.UserHomeDir()
	tokenFile := filepath.Join(homeDir, ".google-mcp-token.json")

	file, err := os.OpenFile(tokenFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Failed to create token file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(token); err != nil {
		fmt.Printf("Failed to save token: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("✅ Authentication successful!")
	fmt.Printf("Token saved to: %s\n", tokenFile)
	fmt.Println("")
	fmt.Println("Next: Run 'task google-mcp:setup:claude' to configure Claude Code")
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}
