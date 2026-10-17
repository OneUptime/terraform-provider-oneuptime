---
page_title: "oneuptime_custom_role Resource - oneuptime"
subcategory: "Other"
description: |-
  A set of permissions you name once, such as Incident Responder, and give to teams. Changing a role changes every team that has it.
---

# oneuptime_custom_role (Resource)

A set of permissions you name once, such as Incident Responder, and give to teams. Changing a role changes every team that has it.

## Example Usage

```terraform
resource "oneuptime_custom_role" "example" {
  name        = "Example custom role"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) What the role is called, such as Incident Responder. Unique in the project, and never the name of a built-in role.

### Optional

- `description` (String) What the role is for, in a sentence or two.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `permission_count` (Number) How many permissions the role holds. OneUptime keeps it as permissions are added to the role and taken out of it.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing custom role by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_custom_role.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_custom_role.example <id>
```
