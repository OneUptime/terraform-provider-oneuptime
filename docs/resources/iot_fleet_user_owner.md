---
page_title: "oneuptime_iot_fleet_user_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add users as owners to your IoT fleets.
---

# oneuptime_iot_fleet_user_owner (Resource)

Add users as owners to your IoT fleets.

~> **Renamed:** this resource was called `oneuptime_io_t_fleet_user_owner` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing iot fleet user owner:

```terraform
moved {
  from = oneuptime_io_t_fleet_user_owner.example
  to   = oneuptime_iot_fleet_user_owner.example
}
```

## Example Usage

```terraform
resource "oneuptime_iot_fleet_user_owner" "example" {
  user_id      = data.oneuptime_user.example.id
  iot_fleet_id = oneuptime_iot_fleet.example.id
}
```

## Schema

### Required

- `iot_fleet_id` (String) ID of your OneUptime IoT Fleet in which this object belongs. The ID of a `oneuptime_iot_fleet`.
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing iot fleet user owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_iot_fleet_user_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_iot_fleet_user_owner.example <id>
```
