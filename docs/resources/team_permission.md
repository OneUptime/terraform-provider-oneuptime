---
page_title: "oneuptime_team_permission Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Permissions for your OneUptime team
---

# oneuptime_team_permission (Resource)

Permissions for your OneUptime team

## Example Usage

```terraform
resource "oneuptime_team_permission" "example" {
  team_id = oneuptime_team.example.id
}
```

## Schema

### Required

- `team_id` (String) ID of Team this permission belongs in. The ID of a `oneuptime_team`.

### Optional

- `is_block_permission` (Boolean) When on, the permission is taken away from the team's members instead of given. A block permission always wins over one that gives it. Defaults to `false`.
- `labels` (Set of String) Labels that limit the permission to the resources carrying one of them. Empty means every resource. IDs of `oneuptime_label` resources.
- `permission` (String) Permission. You can find list of permissions on the Permissions page. A JSON value: write it with `jsonencode()`.
- `scope` (String) Scope of this permission row. One of: All, Owned, Labels. Defaults to All so new permissions apply to every resource in the project unless explicitly narrowed. Allowed values: `All`, `Owned`, `Labels`. Defaults to `All`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `custom_role_id` (String) ID of the custom role this permission comes from. Empty when the permission was given to the team directly. The ID of a `oneuptime_custom_role`.
- `custom_role_permission_id` (String) ID of the custom role's own permission this one copies. Empty when the permission was given to the team directly. The ID of a `oneuptime_custom_role_permission`.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `team_custom_role_id` (String) ID of the custom role given to the team that brought this permission. Empty when the permission was given to the team directly. The ID of a `oneuptime_team_custom_role`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing team permission by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_team_permission.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_team_permission.example <id>
```
