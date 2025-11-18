package env

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

// CommandOutput represents streaming command output
type CommandOutput struct {
	Output        string // Combined stdout/stderr output
	Error         error  // Command execution error
	LocalURL      string // Local preview URL (e.g., "http://localhost:1313")
	DeploymentURL string // Cloudflare Pages deployment URL (e.g., "https://project.pages.dev")
}

// Global Hugo server management
var (
	hugoServerCmd  *exec.Cmd
	hugoServerMux  sync.Mutex
	hugoServerPort = 1313 // Default Hugo server port
)

// extractDeploymentURL extracts the Cloudflare Pages deployment URL from wrangler output
func extractDeploymentURL(output string) string {
	// Wrangler output typically contains: "✨ Deployment complete! Take a peek over at https://abc123.project.pages.dev"
	// Pattern: https://[alphanumeric-].pages.dev
	re := regexp.MustCompile(`https://[a-z0-9-]+\.pages\.dev`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

// StartHugoServer starts a simple HTTP server for local preview of the built site
func StartHugoServer(mockMode bool) CommandOutput {
	hugoServerMux.Lock()
	defer hugoServerMux.Unlock()

	// Stop any existing server first
	if hugoServerCmd != nil {
		StopHugoServer()
	}

	if mockMode {
		localURL := fmt.Sprintf("http://localhost:%d", hugoServerPort)
		return CommandOutput{
			Output:   fmt.Sprintf("Starting preview server (mock mode)...\nServer running at %s", localURL),
			Error:    nil,
			LocalURL: localURL,
		}
	}

	// Use Hugo's built-in server to serve the built site
	// hugo server serves in production mode, no live reload
	localURL := fmt.Sprintf("http://localhost:%d", hugoServerPort)
	hugoServerCmd = exec.Command("hugo", "server", "-e", "production", "--disableLiveReload", "--port", fmt.Sprintf("%d", hugoServerPort))

	// Start the server in background
	if err := hugoServerCmd.Start(); err != nil {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("failed to start preview server: %w", err),
		}
	}

	return CommandOutput{
		Output:   fmt.Sprintf("Preview server started at %s", localURL),
		Error:    nil,
		LocalURL: localURL,
	}
}

// StopHugoServer stops the running Hugo server
func StopHugoServer() CommandOutput {
	hugoServerMux.Lock()
	defer hugoServerMux.Unlock()

	if hugoServerCmd == nil {
		return CommandOutput{
			Output: "No Hugo server is running",
			Error:  nil,
		}
	}

	// Kill the server process
	if err := hugoServerCmd.Process.Kill(); err != nil {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("failed to stop Hugo server: %w", err),
		}
	}

	hugoServerCmd = nil

	return CommandOutput{
		Output: "Hugo preview server stopped",
		Error:  nil,
	}
}

// BuildHugoSite runs `hugo --gc --minify` and returns streaming output
// Also starts a local preview server
func BuildHugoSite(mockMode bool) CommandOutput {
	if mockMode {
		localURL := fmt.Sprintf("http://localhost:%d", hugoServerPort)
		return CommandOutput{
			Output:   "Building Hugo site (mock mode)...\nBuild complete! (mock)\n\nStarting preview server...\nPreview server running at " + localURL,
			Error:    nil,
			LocalURL: localURL,
		}
	}

	// Build the site
	buildResult := runCommand("hugo", "--gc", "--minify")
	if buildResult.Error != nil {
		return buildResult
	}

	// Start preview server
	serverResult := StartHugoServer(mockMode)
	if serverResult.Error != nil {
		// Build succeeded but server failed - return build output with warning
		return CommandOutput{
			Output:   buildResult.Output + "\n\nWarning: Failed to start preview server: " + serverResult.Error.Error(),
			Error:    nil,
			LocalURL: "",
		}
	}

	// Combine build and server outputs
	return CommandOutput{
		Output:   buildResult.Output + "\n\n" + serverResult.Output,
		Error:    nil,
		LocalURL: serverResult.LocalURL,
	}
}

// DeployToPages runs `bunx wrangler pages deploy public --project-name={projectName}` and returns streaming output
func DeployToPages(projectName string, mockMode bool) CommandOutput {
	if mockMode {
		mockURL := fmt.Sprintf("https://%s.pages.dev", projectName)
		return CommandOutput{
			Output:        fmt.Sprintf("Deploying to Cloudflare Pages (mock mode)...\nProject: %s\nDeployment complete! (mock)\nURL: %s", projectName, mockURL),
			Error:         nil,
			DeploymentURL: mockURL,
		}
	}

	if projectName == "" {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("project name is required"),
		}
	}

	result := runCommand("bunx", "wrangler", "pages", "deploy", "public", "--project-name="+projectName)

	// Extract deployment URL from output
	if result.Error == nil {
		deploymentURL := extractDeploymentURL(result.Output)
		result.DeploymentURL = deploymentURL
	}

	return result
}

// CreatePagesProject runs `bunx wrangler pages project create {projectName} --production-branch=main`
// Returns success if project already exists (idempotent)
func CreatePagesProject(projectName string, mockMode bool) CommandOutput {
	if mockMode {
		return CommandOutput{
			Output: fmt.Sprintf("Creating Cloudflare Pages project (mock mode)...\nProject '%s' created successfully (mock)", projectName),
			Error:  nil,
		}
	}

	if projectName == "" {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("project name is required"),
		}
	}

	result := runCommand("bunx", "wrangler", "pages", "project", "create", projectName, "--production-branch=main")

	// Wrangler returns error if project exists - make it idempotent
	if result.Error != nil && strings.Contains(result.Output, "already exists") {
		return CommandOutput{
			Output: result.Output + "\n✓ Project already exists (idempotent success)",
			Error:  nil,
		}
	}

	return result
}

// BuildAndDeploy runs Hugo build followed by Wrangler deploy
func BuildAndDeploy(projectName string, mockMode bool) CommandOutput {
	// Step 1: Build Hugo site
	buildResult := BuildHugoSite(mockMode)
	if buildResult.Error != nil {
		return CommandOutput{
			Output: buildResult.Output,
			Error:  fmt.Errorf("build failed: %w", buildResult.Error),
		}
	}

	// Step 2: Deploy to Pages
	deployResult := DeployToPages(projectName, mockMode)
	if deployResult.Error != nil {
		return CommandOutput{
			Output:   buildResult.Output + "\n\n" + deployResult.Output,
			Error:    fmt.Errorf("deployment failed: %w", deployResult.Error),
			LocalURL: buildResult.LocalURL, // Preserve local URL even if deploy fails
		}
	}

	// Success - combine outputs and URLs
	return CommandOutput{
		Output:        buildResult.Output + "\n\n" + deployResult.Output,
		Error:         nil,
		LocalURL:      buildResult.LocalURL,
		DeploymentURL: deployResult.DeploymentURL,
	}
}

// runCommand executes a command and captures streaming output
func runCommand(name string, args ...string) CommandOutput {
	cmd := exec.Command(name, args...)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("failed to create stdout pipe: %w", err),
		}
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("failed to create stderr pipe: %w", err),
		}
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("failed to start command: %w", err),
		}
	}

	// Read output from both pipes
	var output strings.Builder
	done := make(chan error)

	// Read stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			output.WriteString(line + "\n")
		}
	}()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			output.WriteString(line + "\n")
		}
	}()

	// Wait for command to finish
	go func() {
		done <- cmd.Wait()
	}()

	// Wait for completion
	err = <-done

	// Ensure all output is read
	io.Copy(&output, stdoutPipe)
	io.Copy(&output, stderrPipe)

	if err != nil {
		return CommandOutput{
			Output: output.String(),
			Error:  fmt.Errorf("command failed: %w", err),
		}
	}

	return CommandOutput{
		Output: output.String(),
		Error:  nil,
	}
}
