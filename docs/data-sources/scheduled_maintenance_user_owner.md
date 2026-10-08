---
page_title: "oneuptime_scheduled_maintenance_user_owner Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Add users as owners to your scheduled maintenance event.
---

# oneuptime_scheduled_maintenance_user_owner (Data Source)

Add users as owners to your scheduled maintenance event.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled maintenance user owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_user_owner" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_scheduled_maintenance_user_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `scheduled_maintenance_id` (String) ID of your OneUptime ScheduledMaintenance in which this object belongs. The ID of a `oneuptime_scheduled_maintenance_event`.
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
