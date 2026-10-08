---
page_title: "oneuptime_service_user_owner Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Add users as owners to your Service.
---

# oneuptime_service_user_owner (Data Source)

Add users as owners to your Service.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one service user owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_service_user_owner" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_service_user_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `service_id` (String) ID of your OneUptime Service in which this object belongs. The ID of a `oneuptime_service`.
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
