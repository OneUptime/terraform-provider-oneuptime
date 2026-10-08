# Look up an existing monitor status by name.
data "oneuptime_monitor_status" "example" {
  name = "Offline"
}

output "monitor_status_id" {
  description = "ID of the monitor status"
  value       = data.oneuptime_monitor_status.example.id
}
