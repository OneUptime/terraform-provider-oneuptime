---
page_title: "oneuptime_monitor_status Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Manage monitor status in your project. Monitor Status are Operational, Degraded and Offline for example. Add custom status like Monitoring or more.
---

# oneuptime_monitor_status (Data Source)

Manage monitor status in your project. Monitor Status are Operational, Degraded and Offline for example. Add custom status like Monitoring or more.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor status may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_status" "example" {
  name = "Example monitor status"
}

# Or by id:
data "oneuptime_monitor_status" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_offline_state` (Boolean) Is this monitor in offline state?
- `is_operational_state` (Boolean) Is this monitor in operational state?
- `name` (String) Any friendly name of this object.
- `priority` (Number) Where this status sits in the project's list of monitor statuses, from the healthiest (the lowest number) down to the worst: where monitors are shown together, as on a status page or in a monitor group, the status furthest down the list wins. A new status without a number goes just above the offline status. Setting a number moves the status to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
