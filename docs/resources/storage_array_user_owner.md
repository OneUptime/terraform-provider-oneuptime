---
page_title: "oneuptime_storage_array_user_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add users as owners to your storage arrays.
---

# oneuptime_storage_array_user_owner (Resource)

Add users as owners to your storage arrays.

## Example Usage

```terraform
resource "oneuptime_storage_array_user_owner" "example" {
  user_id          = data.oneuptime_user.example.id
  storage_array_id = oneuptime_storage_array.example.id
}
```

## Schema

### Required

- `storage_array_id` (String) ID of your OneUptime Storage Array in which this object belongs. The ID of a `oneuptime_storage_array`.
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing storage array user owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_storage_array_user_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_storage_array_user_owner.example <id>
```
