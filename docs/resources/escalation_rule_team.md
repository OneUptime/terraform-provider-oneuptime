---
page_title: "oneuptime_escalation_rule_team Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add teams to your escalation rules.
---

# oneuptime_escalation_rule_team (Resource)

Add teams to your escalation rules.

~> **Renamed:** this resource was called `oneuptime_team_on_call_duty_escalation_rule` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing escalation rule team:

```terraform
moved {
  from = oneuptime_team_on_call_duty_escalation_rule.example
  to   = oneuptime_escalation_rule_team.example
}
```

## Example Usage

```terraform
resource "oneuptime_escalation_rule_team" "example" {
  on_call_duty_policy_id                 = oneuptime_on_call_policy.example.id
  on_call_duty_policy_escalation_rule_id = oneuptime_escalation_rule.example.id
}
```

## Schema

### Required

- `on_call_duty_policy_escalation_rule_id` (String) ID of the escalation rule that pages the team. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of the on-call policy the escalation rule is in. The ID of a `oneuptime_on_call_policy`.

### Optional

- `team_id` (String) ID of the team the escalation rule pages. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing escalation rule team by its escalation rule and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_escalation_rule_team.example
  id = "<on_call_duty_policy_escalation_rule_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_escalation_rule_team.example <on_call_duty_policy_escalation_rule_id>/<team_id>
```

Its own ID works too.
