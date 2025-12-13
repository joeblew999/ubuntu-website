# Google MCP Server - API Setup
# Enables all required Google APIs for the MCP server
#
# Prerequisites:
#   1. gcloud auth application-default login
#   2. Set TF_VAR_project_id or create terraform.tfvars
#
# Usage:
#   terraform init
#   terraform plan
#   terraform apply

terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 4.0"
    }
  }
}

provider "google" {
  project = var.project_id
}

variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

# APIs required for Google MCP Server
locals {
  required_apis = [
    "gmail.googleapis.com",
    "calendar-json.googleapis.com",
    "drive.googleapis.com",
    "sheets.googleapis.com",
    "docs.googleapis.com",
    "slides.googleapis.com",
  ]
}

# Enable each API
resource "google_project_service" "apis" {
  for_each = toset(local.required_apis)

  project = var.project_id
  service = each.value

  # Don't disable on destroy - safer for shared projects
  disable_on_destroy = false

  # Don't fail if already enabled
  disable_dependent_services = false
}

output "enabled_apis" {
  description = "List of enabled APIs"
  value       = [for api in google_project_service.apis : api.service]
}

output "project_id" {
  description = "Project ID where APIs were enabled"
  value       = var.project_id
}
