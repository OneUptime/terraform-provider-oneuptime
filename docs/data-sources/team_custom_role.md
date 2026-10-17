---
page_title: "oneuptime_team_custom_role Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  The custom roles a team is given. Its members can do everything the role allows.
---

# oneuptime_team_custom_role (Data Source)

The custom roles a team is given. Its members can do everything the role allows.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team custom role may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team_custom_role" "example" {
  team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_team_custom_role" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `custom_role_id` (String) ID of the custom role given to the team. The ID of a `oneuptime_custom_role`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `team_id` (String) ID of the team given the role. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
