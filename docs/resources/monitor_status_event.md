---
page_title: "oneuptime_monitor_status_event Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Change state of the monitor (Operational to Offline for example)
---

# oneuptime_monitor_status_event (Resource)

Change state of the monitor (Operational to Offline for example)

## Example Usage

```terraform
resource "oneuptime_monitor_status_event" "example" {
  monitor_id = "123e4567-e89b-12d3-a456-426614174000"
  monitor_status_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `monitor_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_status_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `root_cause` (String) What is the root cause of this status change?..
- `ends_at` (String) A date time object..
- `starts_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `is_owner_notified` (Bool) Are owners notified of status change?..
- `status_change_log` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Monitor Admin, Monitor Member, Monitor Viewer, Read Monitor Status Timeline], Update: [No access - you don't have permission for this operation].

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_monitor_status_event.example <id>
```
