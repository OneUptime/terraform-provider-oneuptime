---
page_title: "oneuptime_team_custom_role Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  The custom roles a team is given. Its members can do everything the role allows.
---

# oneuptime_team_custom_role (Resource)

The custom roles a team is given. Its members can do everything the role allows.

## Example Usage

```terraform
resource "oneuptime_team_custom_role" "example" {
  team_id        = oneuptime_team.example.id
  custom_role_id = oneuptime_custom_role.example.id
}
```

## Schema

### Required

- `custom_role_id` (String) ID of the custom role given to the team. The ID of a `oneuptime_custom_role`.
- `team_id` (String) ID of the team given the role. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing team custom role by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_team_custom_role.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_team_custom_role.example <id>
```
