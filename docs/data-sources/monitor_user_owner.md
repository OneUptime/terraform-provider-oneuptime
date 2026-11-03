---
page_title: "oneuptime_monitor_user_owner Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Add users as owners to your monitors.
---

# oneuptime_monitor_user_owner (Data Source)

Add users as owners to your monitors.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor user owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_user_owner" "example" {
  monitor_id = oneuptime_monitor.example.id
  user_email = "alice@example.com"
}

# Or by id:
data "oneuptime_monitor_user_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `monitor_id` (String) ID of the monitor the user owns. The ID of a `oneuptime_monitor`.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user who owns the monitor. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
