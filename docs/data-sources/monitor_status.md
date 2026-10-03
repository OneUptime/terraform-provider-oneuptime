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
- `priority` (Number) Where this status sits in the project's list of monitor statuses, from the healthiest (the lowest number) down to the worst: where monitors are shown together, as on a status page or in a monitor group, the status furthest down the list wins. A new status without a number goes just above the offline status. Setting a number moves the status to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them... Computed.
