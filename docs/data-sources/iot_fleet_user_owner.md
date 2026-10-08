---
page_title: "oneuptime_iot_fleet_user_owner Data Source - oneuptime"
subcategory: "Other"
description: |-
  Add users as owners to your IoT fleets.
---

# oneuptime_iot_fleet_user_owner (Data Source)

Add users as owners to your IoT fleets.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one iot fleet user owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_io_t_fleet_user_owner` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_iot_fleet_user_owner" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_iot_fleet_user_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `iot_fleet_id` (String) ID of your OneUptime IoT Fleet in which this object belongs. The ID of a `oneuptime_iot_fleet`.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
