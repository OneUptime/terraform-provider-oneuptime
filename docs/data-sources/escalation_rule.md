---
page_title: "oneuptime_escalation_rule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call escalation rule for the on-call policy.
---

# oneuptime_escalation_rule (Data Source)

Manage on-call escalation rule for the on-call policy.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one escalation rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_escalation_rule" "example" {
  name = "Example escalation rule"
}

# Or by id:
data "oneuptime_escalation_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `escalate_after_in_minutes` (Number) How long should we wait before we execute the next escalation rule?
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.
- `on_call_duty_policy_id` (String) ID of your On-Call Policy where this escalation rule belongs. The ID of a `oneuptime_on_call_policy`.
- `order` (Number) This rule's level in its on-call policy: the lowest is paged first, then the next, and so on. Rules that share a level are paged in the order of their IDs. Setting it places this rule and moves no other, so a configuration that numbers every rule of a policy ends up exactly as written. Left out when you create a rule, it becomes the policy's last level.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
