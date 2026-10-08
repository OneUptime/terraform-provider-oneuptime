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

}
```

## Schema

### Optional

- `is_block_permission` (Boolean) Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this permission is scoped at. IDs of `oneuptime_label` resources.
- `permission` (String) Permission. You can find list of permissions on the Permissions page. A JSON value: write it with `jsonencode()`.
- `scope` (String) Scope of this permission row. One of: All, Owned, Labels. Defaults to All so new permissions apply to every resource in the project unless explicitly narrowed. Defaults to `All`.
- `team_id` (String) ID of Team this permission belongs in. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
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
