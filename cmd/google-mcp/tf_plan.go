package main

// tfPlan is google-mcp tf-plan, the task google-mcp:tf-plan: what Terraform would change in the Google Cloud project's APIs (terraform/google-mcp), its arguments passed to
// terraform.
func tfPlan(args []string) { runTerraform("plan", args) }
