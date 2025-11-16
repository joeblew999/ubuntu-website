package env

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

// ValidateCloudflareToken validates a Cloudflare API token and returns the token name
func ValidateCloudflareToken(token string) (string, error) {
	if token == "" || token == PlaceholderToken {
		return "", fmt.Errorf("no token to validate")
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// Verify token
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/user/tokens/verify", nil)
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

	// Try to get token details to retrieve the name
	// This is optional - token may not have permission to read its own details
	// Requires "User: API Tokens: Read" permission which is not needed for deployment
	tokenID := verifyResp.Result.ID
	tokenReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/user/tokens/%s", tokenID), nil)
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

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s", accountID)
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
