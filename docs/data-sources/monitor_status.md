---
page_title: "oneuptime_monitor_status Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Manage monitor status in your project. Monitor Status are Operational, Degraded and Offline for example. Add custom status like Monitoring or more.
---

# oneuptime_monitor_status (Data Source)

Manage monitor status in your project. Monitor Status are Operational, Degraded and Offline for example. Add custom status like Monitoring or more. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_monitor_status" "by_name" {
  name = "example-monitor_status"
}

data "oneuptime_monitor_status" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `color` (String) Color object. Computed.
- `is_operational_state` (Bool) Is this monitor in operational state?.. Computed.
- `is_offline_state` (Bool) Is this monitor in offline state?.. Computed.
- `priority` (Number) Order / Priority of this status. Behaves like an insertion slot: creating a status with priority P shifts every existing status with priority >= P up by one, and priority cannot be changed after creation (delete and recreate instead). When managing statuses declaratively, use high, gapped values (e.g. 101, 102, 103) created in ascending order so existing statuses are never shifted... Computed.
