# Example usage of oneuptime_probe resource
resource "oneuptime_probe" "example" {
  name        = "example-probe"
  description = "Example probe created by Terraform"
}

# Output the resource ID
output "probe_id" {
  description = "ID of the created probe"
  value       = oneuptime_probe.example.id
}
