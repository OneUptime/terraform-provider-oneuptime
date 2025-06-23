# Example usage of oneuptime_user resource
resource "oneuptime_user" "example" {
  name        = "example-user"
  description = "Example user created by Terraform"
}

# Output the resource ID
output "user_id" {
  description = "ID of the created user"
  value       = oneuptime_user.example.id
}
