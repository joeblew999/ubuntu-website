package main

import (
	"os"
	"os/exec"

	googleauth "github.com/joeblew999/ubuntu-website/internal/google/auth"
)

// terraformDir is the Terraform that enables the Google Cloud project's APIs (terraform.tfvars has its id).
const terraformDir = "terraform/google-mcp"

// runTerraform runs terraform <command> <args> in terraformDir, after making sure the application-default credentials
// it signs in with work: checked with Google, and signed in again in the browser if they've expired, so a stale
// sign-in is fixed here rather than failing Terraform with invalid_grant. Its input and output are the terminal's
// (apply asks to confirm; -auto-approve skips that).
func runTerraform(command string, args []string) {
	if command != "init" {
		if err := googleauth.EnsureGcloudCredentials(os.Stdout); err != nil {
			fail(err.Error())
		}
	}
	cmd := exec.Command("terraform", append([]string{"-chdir=" + terraformDir, command}, args...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		}
		fail(err.Error())
	}
}
