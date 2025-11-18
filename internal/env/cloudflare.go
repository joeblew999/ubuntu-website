package env

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// CloudflareVerifyResponse represents the token verification API response
type CloudflareVerifyResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"result"`
}

// CloudflareTokenResponse represents a token details API response
type CloudflareTokenResponse struct {
	Success bool `json:"success"`
	Result  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

// CloudflareAccountResponse represents the account info API response
type CloudflareAccountResponse struct {
	Success bool `json:"success"`
	Result  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

// CloudflareAccountsResponse represents the accounts list API response
type CloudflareAccountsResponse struct {
	Success bool `json:"success"`
	Result  []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

// ValidateCloudflareToken validates a Cloudflare API token and returns the token name
func ValidateCloudflareToken(token string) (string, error) {
	if token == "" || token == PlaceholderToken {
		return "", fmt.Errorf("no token to validate")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Verify token
	req, err := http.NewRequest("GET", CloudflareAPITokenVerifyURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var verifyResp CloudflareVerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if !verifyResp.Success {
		if len(verifyResp.Errors) > 0 {
			return "", fmt.Errorf("invalid token: %s", verifyResp.Errors[0].Message)
		}
		return "", fmt.Errorf("token verification failed")
	}

	// Get token details to retrieve the name
	// This requires "User: API Tokens: Read" permission
	// If token lacks this permission, validation will still succeed but won't show the name
	tokenID := verifyResp.Result.ID
	tokenReq, err := http.NewRequest("GET", fmt.Sprintf(CloudflareAPITokenInfoURL, tokenID), nil)
	if err != nil {
		// Can't create request - return without name
		return "", nil
	}

	tokenReq.Header.Set("Authorization", "Bearer "+token)
	tokenReq.Header.Set("Content-Type", "application/json")

	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		// Can't fetch details - return without name
		return "", nil
	}
	defer tokenResp.Body.Close()

	tokenBody, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		// Can't read response - return without name
		return "", nil
	}

	var tokenDetails CloudflareTokenResponse
	if err := json.Unmarshal(tokenBody, &tokenDetails); err != nil {
		// Can't parse response - return without name
		return "", nil
	}

	if !tokenDetails.Success {
		// Token doesn't have permission to read its own details - return without name
		return "", nil
	}

	return tokenDetails.Result.Name, nil
}

// ValidateCloudflareAccount validates the account ID with the given token
func ValidateCloudflareAccount(token, accountID string) (string, error) {
	if accountID == "" {
		return "", fmt.Errorf("no account ID to validate")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf(CloudflareAPIAccountURL, accountID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to verify account: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var accountResp CloudflareAccountResponse
	if err := json.Unmarshal(body, &accountResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if !accountResp.Success {
		// Check for specific error in response
		if resp.StatusCode == 403 {
			return "", fmt.Errorf("token lacks permission to access account %s (need Account:Read permission)", accountID)
		}
		if resp.StatusCode == 404 {
			return "", fmt.Errorf("account ID %s not found or not accessible with this token", accountID)
		}
		return "", fmt.Errorf("account validation failed (status: %d)", resp.StatusCode)
	}

	return accountResp.Result.Name, nil
}

// ValidateCloudflareProjectName validates a Cloudflare Pages project name
// Project names must be lowercase alphanumeric with hyphens, 1-63 characters
func ValidateCloudflareProjectName(projectName string) error {
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Cloudflare project name requirements
	if len(projectName) > 63 {
		return fmt.Errorf("project name must be 63 characters or less")
	}

	// Must match: lowercase letters, numbers, hyphens only
	// Cannot start or end with hyphen
	matched, err := regexp.MatchString(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`, projectName)
	if err != nil {
		return fmt.Errorf("failed to validate project name format: %w", err)
	}

	if !matched {
		return fmt.Errorf("project name must contain only lowercase letters, numbers, and hyphens")
	}

	return nil
}

// GetCloudflareAccounts fetches all accounts accessible by the token
// Returns the first account ID and name if exactly one account is found
func GetCloudflareAccounts(token string) (accountID, accountName string, err error) {
	if token == "" || token == PlaceholderToken {
		return "", "", fmt.Errorf("no token provided")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", CloudflareAPIAccountsURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch accounts: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var accountsResp CloudflareAccountsResponse
	if err := json.Unmarshal(body, &accountsResp); err != nil {
		return "", "", fmt.Errorf("failed to parse response: %w", err)
	}

	if !accountsResp.Success {
		return "", "", fmt.Errorf("failed to fetch accounts (status: %d)", resp.StatusCode)
	}

	if len(accountsResp.Result) == 0 {
		return "", "", fmt.Errorf("no accounts found for this token")
	}

	// Return the first account (most tokens only have access to one account)
	return accountsResp.Result[0].ID, accountsResp.Result[0].Name, nil
}
