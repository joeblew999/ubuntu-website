package main

// tfDestroy is google-mcp tf-destroy, the task google-mcp:tf-destroy: the Terraform-managed API enablement removed (terraform/google-mcp), its arguments passed to
// terraform.
func tfDestroy(args []string) { runTerraform("destroy", args) }
