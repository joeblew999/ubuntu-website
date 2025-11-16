package env

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const envFile = ".env"

// Environment variable keys - kept for backward compatibility with existing code
// New code should use getEnvKey() to get keys from struct tags
const (
	EnvCloudflareToken   = "CLOUDFLARE_API_TOKEN"
	EnvCloudflareAccount = "CLOUDFLARE_ACCOUNT_ID"
	EnvCloudflareProject = "CLOUDFLARE_PROJECT_NAME"
	EnvClaudeAPIKey      = "CLAUDE_API_KEY"
	EnvClaudeWorkspace   = "CLAUDE_WORKSPACE"
)

// Placeholder values - kept for backward compatibility
// New code should use isPlaceholder() helper function
const (
	PlaceholderToken = "your-token-here"
	PlaceholderKey   = "your-api-key-here"
)

// EnvConfig holds environment configuration
// All env vars are defined here with struct tags for metadata
// This is the single source of truth for all environment variables
type EnvConfig struct {
	CloudflareToken   string `env:"CLOUDFLARE_API_TOKEN" default:"your-token-here" comment:"Cloudflare credentials (for deployment)" validate:"cloudflare_token" required:"true"`
	CloudflareAccount string `env:"CLOUDFLARE_ACCOUNT_ID" default:"your-account-id" validate:"cloudflare_account"`
	CloudflareProject string `env:"CLOUDFLARE_PROJECT_NAME" default:"your-project-name"`
	ClaudeAPIKey      string `env:"CLAUDE_API_KEY" default:"your-api-key-here" comment:"Claude API key (for translation)" validate:"claude_api_key"`
	ClaudeWorkspace   string `env:"CLAUDE_WORKSPACE" default:"your-workspace-name" comment:"Claude Workspace (recommended for project isolation)"`
}

// getEnvKey returns the env key name for a struct field using reflection
func getEnvKey(field reflect.StructField) string {
	return field.Tag.Get("env")
}

// getDefaultValue returns the default value for a field
func getDefaultValue(field reflect.StructField) string {
	return field.Tag.Get("default")
}

// getComment returns the comment for a field (for grouping in .env file)
func getComment(field reflect.StructField) string {
	return field.Tag.Get("comment")
}

// getValidateName returns the validation function name for a field
func getValidateName(field reflect.StructField) string {
	return field.Tag.Get("validate")
}

// isRequired returns whether a field is required
func isRequired(field reflect.StructField) bool {
	return field.Tag.Get("required") == "true"
}

// setFieldByEnvKey sets a struct field value by env key name using reflection
func setFieldByEnvKey(cfg *EnvConfig, envKey string, value string) bool {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if getEnvKey(field) == envKey {
			v.Field(i).SetString(value)
			return true
		}
	}
	return false
}

// getFieldByEnvKey gets a struct field value by env key name using reflection
func getFieldByEnvKey(cfg *EnvConfig, envKey string) (string, bool) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if getEnvKey(field) == envKey {
			return v.Field(i).String(), true
		}
	}
	return "", false
}

// LoadEnv reads the .env file and returns the configuration
func LoadEnv() (*EnvConfig, error) {
	cfg := &EnvConfig{}

	file, err := os.Open(envFile)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // Return empty config if file doesn't exist
		}
		return nil, fmt.Errorf("failed to open .env: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Use reflection to set field by env key
		setFieldByEnvKey(cfg, key, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .env: %w", err)
	}

	return cfg, nil
}

// CreateEnv creates a new .env file with default values
func CreateEnv() error {
	cfg := &EnvConfig{}
	return WriteEnv(cfg)
}

// UpdateEnv updates a specific key in the .env file
func UpdateEnv(key, value string) error {
	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

	// Use reflection to update field by env key
	if !setFieldByEnvKey(cfg, key, value) {
		return fmt.Errorf("unknown environment key: %s", key)
	}

	// Write back the entire file
	return WriteEnv(cfg)
}

// WriteEnv writes the complete configuration to .env
func WriteEnv(cfg *EnvConfig) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	var content strings.Builder
	var lastComment string

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envKey := getEnvKey(field)
		if envKey == "" {
			continue // Skip fields without env tag
		}

		// Get field value or default
		value := v.Field(i).String()
		if value == "" {
			value = getDefaultValue(field)
		}

		// Add comment header if this field has one and it's different from last
		comment := getComment(field)
		if comment != "" && comment != lastComment {
			if i > 0 {
				content.WriteString("\n")
			}
			content.WriteString("# ")
			content.WriteString(comment)
			content.WriteString("\n")
			lastComment = comment
		}

		// Write the key=value line
		content.WriteString(envKey)
		content.WriteString("=")
		content.WriteString(value)
		content.WriteString("\n")
	}

	if err := os.WriteFile(envFile, []byte(content.String()), 0600); err != nil {
		return fmt.Errorf("failed to write .env: %w", err)
	}

	return nil
}

// EnvExists checks if .env file exists
func EnvExists() bool {
	_, err := os.Stat(envFile)
	return err == nil
}

// GetEnvPath returns the absolute path to the .env file
func GetEnvPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", wd, envFile), nil
}

// isPlaceholder checks if a value is a placeholder/default value
func isPlaceholder(value string) bool {
	return value == "" || strings.HasPrefix(value, "your-") || strings.HasPrefix(value, "your_")
}
