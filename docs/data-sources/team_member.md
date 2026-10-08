---
page_title: "oneuptime_team_member Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  This model connects users and teams
---

# oneuptime_team_member (Data Source)

This model connects users and teams

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team member may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team_member" "example" {
  team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_team_member" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `has_accepted_invitation` (Boolean) Has this team member accepted invitation.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `team_id` (String) ID of Team this user belongs to. The ID of a `oneuptime_team`.
- `user_id` (String) ID of User who belongs to this team. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `invitation_accepted_at` (String) When did this team member accept invitation.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
