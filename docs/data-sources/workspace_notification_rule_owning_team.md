---
page_title: "oneuptime_workspace_notification_rule_owning_team Data Source - oneuptime"
subcategory: "Slack, Teams & Video Calls"
description: |-
  Teams that own a workspace notification rule: only their members, and project admins, can change or delete it.
---

# oneuptime_workspace_notification_rule_owning_team (Data Source)

Teams that own a workspace notification rule: only their members, and project admins, can change or delete it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one workspace notification rule owning team may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_workspace_notification_rule_owning_team" "example" {
  team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_workspace_notification_rule_owning_team" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `team_id` (String) ID of the team that owns the rule. The ID of a `oneuptime_team`.
- `workspace_notification_rule_id` (String) ID of the workspace notification rule the team owns. The ID of a `oneuptime_workspace_notification_rule`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
