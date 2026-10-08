# Example usage of oneuptime_monitor
# The project's own statuses and severities, looked up by name.
data "oneuptime_monitor_status" "operational" {
  name = "Operational"
}

data "oneuptime_monitor_status" "offline" {
  name = "Offline"
}

data "oneuptime_incident_severity" "critical" {
  name = "Critical Incident"
}

resource "oneuptime_monitor" "example" {
  name                = "Example website"
  description         = "Checks https://example.com every minute"
  monitor_type        = "Website"
  monitoring_interval = "* * * * *"

  monitor_steps = [{
    monitor_destination      = "https://example.com"
    monitor_destination_type = "URL"
    request_type             = "GET"

    # Evaluated top to bottom; the first that matches wins.
    criteria = [
      {
        name                  = "Offline"
        description           = "The website does not answer"
        filter_condition      = "Any"
        change_monitor_status = true
        monitor_status_id     = data.oneuptime_monitor_status.offline.id
        create_incidents      = true

        filters = [
          { check_on = "Is Online", filter_type = "False" },
        ]

        incidents = [{
          title                 = "Example website is down"
          description           = "The website did not respond to the probe."
          incident_severity_id  = data.oneuptime_incident_severity.critical.id
          auto_resolve_incident = true
        }]
      },
      {
        name                  = "Online"
        description           = "The website answers"
        filter_condition      = "All"
        change_monitor_status = true
        monitor_status_id     = data.oneuptime_monitor_status.operational.id

        filters = [
          { check_on = "Is Online", filter_type = "True" },
        ]
      },
    ]
  }]
}

output "monitor_id" {
  description = "ID of the created monitor"
  value       = oneuptime_monitor.example.id
}
