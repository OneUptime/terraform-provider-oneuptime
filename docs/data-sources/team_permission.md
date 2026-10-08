---
page_title: "oneuptime_team_permission Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Permissions for your OneUptime team
---

# oneuptime_team_permission (Data Source)

Permissions for your OneUptime team

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team permission may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team_permission" "example" {
  team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_team_permission" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_block_permission` (Boolean) Permissions - Create: [Project Owner, Project Admin, Edit Team Permissions], Read: [Project Owner, Project Admin, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Teams], Update: [Project Owner, Project Admin, Edit Team Permissions, Edit Team]
- `scope` (String) Scope of this permission row. One of: All, Owned, Labels. Defaults to All so new permissions apply to every resource in the project unless explicitly narrowed.
- `team_id` (String) ID of Team this permission belongs in. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this permission is scoped at. IDs of `oneuptime_label` resources.
- `permission` (String) Permission. You can find list of permissions on the Permissions page. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
