terraform {
  required_providers {
    oneuptime = {
      source  = "oneuptime/oneuptime"
      version = "~> 14.0"
    }
  }
}

provider "oneuptime" {
  oneuptime_url = "https://oneuptime.com" # Optional: defaults to oneuptime.com; or ONEUPTIME_URL
  api_key       = var.oneuptime_api_key  # or ONEUPTIME_API_KEY
}

# Configure variables
variable "oneuptime_api_key" {
  description = "Project API key for oneuptime"
  type        = string
  sensitive   = true
}
