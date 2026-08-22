---
page_title: "oneuptime_monitor Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything.
---

# oneuptime_monitor (Resource)

Monitor is anything that monitors your API, Websites, IP, Network or more. You can also create static monitor that does not monitor anything.

## Example Usage

```terraform
resource "oneuptime_monitor" "example" {
  name = "Example short text"
  monitor_type = "Manual"
  description = "This is an example of longer text content that might be stored in this field."
  monitor_steps = [
    {
      monitor_destination      = "https://your-service.example.com"
      monitor_destination_type = "URL"
      request_type             = "GET"
      criteria = [
        {
          name             = "Check if online"
          filter_condition = "All"
          filters = [
            {
              check_on = "Is Online"
            }
          ]
        }
      ]
    }
  ]
}
```

## Schema

### Required

- `name` (String) Any friendly name for this monitor..
- `monitor_type` (String) What is the type of this monitor? Website? API? etc... Allowed values: `Manual`, `Website`, `API`, `Ping`, `Kubernetes`, `Docker`, `Host`, `Podman`, `Docker Swarm`, `Proxmox`, `Ceph`, `IoT Device`, `IP`, `Incoming Request`, `Incoming Email`, `Port`, `Server`, `SSL Certificate`, `SQL Query`, `Synthetic Monitor`, `Custom JavaScript Code`, `Logs`, `Metrics`, `Traces`, `Exceptions`, `Profiles`, `Security Events`, `Network Device`, `DNS`, `DNSSEC`, `Domain`, `External Status Page`.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `depends_on_monitors` (Set) Parent monitors this monitor depends on. When a parent is offline (or in one of the configured suppression statuses), alerts and incidents from this monitor are suppressed at creation time — the monitor keeps evaluating and its status timeline still updates...
- `suppress_alerts_when_parent_monitor_statuses` (Set) Parent monitor statuses that suppress this monitor's alerts and incidents. When empty, statuses flagged as offline suppress (the default). Only used when Depends On Monitors is set...
- `monitor_template_id` (String) A unique identifier for an object, represented as a UUID..
- `current_monitor_status_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_steps` (Block List) MonitorSteps object.
- `monitoring_interval` (String) How often would you like OneUptime to monitor this resource?..
- `custom_fields` (String) Custom Fields on this resource...
- `disable_active_monitoring` (Bool) Disable active monitoring for this resource?..
- `incoming_request_monitor_heartbeat_checked_at` (String) A date time object..
- `telemetry_monitor_next_monitor_at` (String) A date time object..
- `telemetry_monitor_last_monitor_at` (String) A date time object..
- `server_monitor_request_received_at` (String) A date time object..
- `incoming_monitor_request` (String) Incoming Monitor Request for Incoming Request Monitor..
- `server_monitor_response` (String) Server Monitor Response for Server Monitor..
- `minimum_probe_agreement` (Number) Minimum number of probes that must agree on a status before the monitor status changes. If null, all enabled and connected probes must agree...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `is_owner_notified_of_resource_creation` (Bool) Are owners notified of when this resource is created?..
- `disable_active_monitoring_because_of_scheduled_maintenance_event` (Bool) Disable Monitoring because of Ongoing Scheduled Maintenance Event..
- `disable_active_monitoring_because_of_manual_incident` (Bool) Disable Monitoring because of Incident which is creeated manually by user...
- `server_monitor_secret_key` (String) A unique identifier for an object, represented as a UUID..
- `incoming_request_secret_key` (String) A unique identifier for an object, represented as a UUID..
- `incoming_email_secret_key` (String) A unique identifier for an object, represented as a UUID..
- `incoming_email_monitor_last_email_received_at` (String) A date time object..
- `incoming_email_monitor_request` (String) This field is for Incoming Email Monitor only. Last email data received...
- `incoming_email_monitor_heartbeat_checked_at` (String) A date time object..
- `is_all_probes_disconnected_from_this_monitor` (Bool) All Probes Disconnected From This Monitor. Is this monitor not being monitored?..
- `is_no_probe_enabled_on_this_monitor` (Bool) No Probe Enabled On This Monitor. Is this monitor not being monitored?..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_monitor.example <id>
```
