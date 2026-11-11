---
page_title: "oneuptime_team_team_owner Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Add teams as owners to your teams. The members of an owning team look after who is in the owned team and what it may do.
---

# oneuptime_team_team_owner (Data Source)

Add teams as owners to your teams. The members of an owning team look after who is in the owned team and what it may do.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team team owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team_team_owner" "example" {
  owned_team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_team_team_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `owned_team_id` (String) ID of the team the owning team owns. The ID of a `oneuptime_team`.
- `team_id` (String) ID of the team that owns the team. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
