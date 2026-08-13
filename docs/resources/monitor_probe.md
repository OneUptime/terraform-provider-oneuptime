---
page_title: "oneuptime_monitor_probe Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Add probes to monitor your resource from multiple locations around the world.
---

# oneuptime_monitor_probe (Resource)

Add probes to monitor your resource from multiple locations around the world.

## Example Usage

```terraform
resource "oneuptime_monitor_probe" "example" {
  probe_id = "123e4567-e89b-12d3-a456-426614174000"
  monitor_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `probe_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `last_ping_at` (String) A date time object..
- `next_ping_at` (String) A date time object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_enabled` (Bool) Monitor probe is_enabled.

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `last_monitoring_log` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Probe], Update: [No access - you don't have permission for this operation].

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_monitor_probe.example <id>
```
