---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_monitor Resource - oneuptime"
subcategory: ""
description: |-
  Monitor resource
---

# oneuptime_monitor (Resource)

Monitor resource

## Example Usage

```terraform
resource "oneuptime_monitor" "example" {
  project_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "example-resource"
  monitor_type = "example-monitor_type"
  current_monitor_status_id = "123e4567-e89b-12d3-a456-426614174000"
  is_owner_notified_of_resource_creation = "example-is_owner_notified_of_resource_creation"
  disable_active_monitoring = "example-disable_active_monitoring"
  description = "Example resource"
}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Required.
- `name` (String) Name. Required.
- `description` (String) Description. Optional.
- `labels` (List) Labels. Optional.
- `monitor_type` (String) Monitor Type. Required.
- `current_monitor_status_id` (String) A unique identifier for an object, represented as a UUID.. Required.
- `monitor_steps` (Map) MonitorSteps object. Optional.
- `monitoring_interval` (String) Monitoring Interval. Optional.
- `custom_fields` (Map) Custom Fields. Optional.
- `is_owner_notified_of_resource_creation` (Bool) Are Owners Notified Of Resource Creation?. Required.
- `disable_active_monitoring` (Bool) Disable Monitoring. Required.
- `incoming_request_monitor_heartbeat_checked_at` (Map) A date time object.. Optional.
- `telemetry_monitor_next_monitor_at` (Map) A date time object.. Optional.
- `telemetry_monitor_last_monitor_at` (Map) A date time object.. Optional.
- `server_monitor_request_received_at` (Map) A date time object.. Optional.
- `server_monitor_secret_key` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `incoming_request_secret_key` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `incoming_monitor_request` (Map) Incoming Monitor Request. Optional.
- `server_monitor_response` (Map) Server Monitor Response. Optional.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `slug` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `disable_active_monitoring_because_of_scheduled_maintenance_event` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]. Computed.
- `disable_active_monitoring_because_of_manual_incident` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]. Computed.
- `is_all_probes_disconnected_from_this_monitor` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]. Computed.
- `is_no_probe_enabled_on_this_monitor` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor], Update: [No access - you don't have permission for this operation]. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_monitor.example <id>
```
