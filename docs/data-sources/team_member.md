---
page_title: "oneuptime_team_member Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Add users to your teams.
---

# oneuptime_team_member (Data Source)

Add users to your teams.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team member may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team_member" "example" {
  team_id    = oneuptime_team.example.id
  user_email = "alice@example.com"
}

# Or by id:
data "oneuptime_team_member" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `has_accepted_invitation` (Boolean) Whether the person has joined: false while their invitation waits, true once they accept it. Someone in the project already joins at once.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `team_id` (String) ID of the team the person is in. The ID of a `oneuptime_team`.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the person in the team. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `invitation_accepted_at` (String) When the person accepted their invitation.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
