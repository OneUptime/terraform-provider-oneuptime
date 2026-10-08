---
page_title: "oneuptime_team_member Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  This model connects users and teams
---

# oneuptime_team_member (Resource)

This model connects users and teams

## Example Usage

```terraform
resource "oneuptime_team_member" "example" {
  user_id = data.oneuptime_user.example.id
}
```

## Schema

### Required

- `user_id` (String) ID of User who belongs to this team. The ID of a `oneuptime_user` (see the data source).

### Optional

- `has_accepted_invitation` (Boolean) Has this team member accepted invitation. Defaults to `false`.
- `invitation_accepted_at` (String) When did this team member accept invitation.
- `team_id` (String) ID of Team this user belongs to. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing team member by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_team_member.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_team_member.example <id>
```
