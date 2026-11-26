---
page_title: "oneuptime_escalation_rule_team Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add teams to your escalation rules.
---

# oneuptime_escalation_rule_team (Data Source)

Add teams to your escalation rules.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one escalation rule team may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_team_on_call_duty_escalation_rule` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_escalation_rule_team" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_escalation_rule_team" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `on_call_duty_policy_escalation_rule_id` (String) ID of the escalation rule that pages the team. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of the on-call policy the escalation rule is in. The ID of a `oneuptime_on_call_policy`.
- `team_id` (String) ID of the team the escalation rule pages. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
