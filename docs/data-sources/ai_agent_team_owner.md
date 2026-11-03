---
page_title: "oneuptime_ai_agent_team_owner Data Source - oneuptime"
subcategory: "AI"
description: |-
  Add teams as owners to your AI agents.
---

# oneuptime_ai_agent_team_owner (Data Source)

Add teams as owners to your AI agents.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai agent team owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_ai_agent_owner_team` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_ai_agent_team_owner" "example" {
  team_id = oneuptime_team.example.id
}

# Or by id:
data "oneuptime_ai_agent_team_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_agent_id` (String) ID of the AI agent the team owns. The ID of a `oneuptime_ai_agent`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `team_id` (String) ID of the team that owns the AI agent. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
