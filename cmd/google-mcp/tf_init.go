package main

// tfInit is google-mcp tf-init, the task google-mcp:tf-init: Terraform initialized (its providers) for the Google Cloud project's APIs (terraform/google-mcp), its arguments passed to
// terraform.
func tfInit(args []string) { runTerraform("init", args) }
