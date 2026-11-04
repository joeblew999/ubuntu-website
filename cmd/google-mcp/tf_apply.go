package main

// tfApply is google-mcp tf-apply, the task google-mcp:tf-apply: the Google Cloud project's APIs enabled as Terraform says (terraform/google-mcp), its arguments passed to
// terraform.
func tfApply(args []string) { runTerraform("apply", args) }
