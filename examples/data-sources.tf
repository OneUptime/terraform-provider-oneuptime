# Example usage of oneuptime_user data source
data "oneuptime_user" "example" {
  name = "example-user"
}

# Output the data source result
output "user_result" {
  description = "Result of the user data source"
  value       = data.oneuptime_user.example
}
