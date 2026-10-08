---
page_title: "oneuptime_storage_array_team_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your storage arrays.
---

# oneuptime_storage_array_team_owner (Resource)

Add teams as owners to your storage arrays.

## Example Usage

```terraform
resource "oneuptime_storage_array_team_owner" "example" {
  team_id          = oneuptime_team.example.id
  storage_array_id = oneuptime_storage_array.example.id
}
```

## Schema

### Required

- `storage_array_id` (String) ID of your OneUptime Storage Array in which this object belongs. The ID of a `oneuptime_storage_array`.
- `team_id` (String) ID of your OneUptime Team in which this object belongs. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing storage array team owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_storage_array_team_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_storage_array_team_owner.example <id>
```
