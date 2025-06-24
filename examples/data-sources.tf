# Example usage of oneuptime_user_data data source
data "oneuptime_user_data" "example" {
  name = "example-user_data"
}

# Output the data source result
output "user_data_result" {
  description = "Result of the user_data data source"
  value       = data.oneuptime_user_data.example
}
