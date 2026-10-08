---
page_title: "oneuptime_project_user_profile Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Stores user profile data including custom fields for each user in a project.
---

# oneuptime_project_user_profile (Resource)

Stores user profile data including custom fields for each user in a project.

## Example Usage

```terraform
resource "oneuptime_project_user_profile" "example" {
  user_id = data.oneuptime_user.example.id
}
```

## Schema

### Required

- `user_id` (String) ID of User this profile belongs to. The ID of a `oneuptime_user` (see the data source).

### Optional

- `custom_fields` (String) Custom Fields for this user in this project. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing project user profile by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_project_user_profile.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_project_user_profile.example <id>
```
