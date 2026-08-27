terraform {
  required_providers {
    oneuptime = {
      source = "oneuptime/oneuptime"
      version = "12.0.24"
    }
  }
}

provider "oneuptime" {
  oneuptime_url = "oneuptime.com" # Optional, defaults to oneuptime.com (provider appends /api automatically)
  api_key       = var.oneuptime_api_key
}

# Configure variables
variable "oneuptime_api_key" {
  description = "API key for oneuptime"
  type        = string
  sensitive   = true
}
