---
page_title: "oneuptime_team_team_owner Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Add teams as owners to your teams. The members of an owning team look after who is in the owned team and what it may do.
---

# oneuptime_team_team_owner (Resource)

Add teams as owners to your teams. The members of an owning team look after who is in the owned team and what it may do.

## Example Usage

```terraform
resource "oneuptime_team_team_owner" "example" {
  owned_team_id = oneuptime_team.example.id
  team_id       = oneuptime_team.example.id
}
```

## Schema

### Required

- `owned_team_id` (String) ID of the team the owning team owns. The ID of a `oneuptime_team`.
- `team_id` (String) ID of the team that owns the team. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing team team owner by its team and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_team_team_owner.example
  id = "<owned_team_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_team_team_owner.example <owned_team_id>/<team_id>
```

Its own ID works too.
