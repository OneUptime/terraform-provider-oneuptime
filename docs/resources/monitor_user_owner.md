---
page_title: "oneuptime_monitor_user_owner Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Add users as owners to your monitors.
---

# oneuptime_monitor_user_owner (Resource)

Add users as owners to your monitors.

## Example Usage

```terraform
resource "oneuptime_monitor_user_owner" "example" {
  user_id    = data.oneuptime_user.example.id
  monitor_id = oneuptime_monitor.example.id
}
```

## Schema

### Required

- `monitor_id` (String) ID of your OneUptime Monitor in which this object belongs. The ID of a `oneuptime_monitor`.
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).

### Optional

- `is_owner_notified` (Boolean) Are owners notified of this resource ownership? Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing monitor user owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_monitor_user_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_monitor_user_owner.example <id>
```
