---
page_title: "oneuptime_monitor_status_event Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Change state of the monitor (Operational to Offline for example)
---

# oneuptime_monitor_status_event (Data Source)

Change state of the monitor (Operational to Offline for example)

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor status event may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_status_event" "example" {
  monitor_id = oneuptime_monitor.example.id
}

# Or by id:
data "oneuptime_monitor_status_event" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of status change?
- `monitor_id` (String) Relation to Monitor ID Resource in which this object belongs. The ID of a `oneuptime_monitor`.
- `monitor_status_id` (String) Relation to Monitor Status ID Resource in which this object belongs. The ID of a `oneuptime_monitor_status`.
- `root_cause` (String) What is the root cause of this status change?

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ends_at` (String) When did this status change end?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `starts_at` (String) When did this status change?
- `status_change_log` (String) A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
