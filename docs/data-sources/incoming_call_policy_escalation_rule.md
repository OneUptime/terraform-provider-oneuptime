---
page_title: "oneuptime_incoming_call_policy_escalation_rule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage escalation rules for incoming call policies that define who to call and in what order
---

# oneuptime_incoming_call_policy_escalation_rule (Data Source)

Manage escalation rules for incoming call policies that define who to call and in what order

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming call policy escalation rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_call_policy_escalation_rule" "example" {
  last_modified_by_user_id = data.oneuptime_user.example.id
  user_email               = "alice@example.com"
}

# Or by id:
data "oneuptime_incoming_call_policy_escalation_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Optional description of this escalation rule.
- `escalate_after_seconds` (Number) How long, in seconds, the phone rings before the call moves on to the next rule. 20 when left out; a time below 5 or above 600 rings for 5 or 600, the limits Twilio takes.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incoming_call_policy_id` (String) ID of the parent Incoming Call Policy. The ID of a `oneuptime_incoming_call_policy`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Rule name (e.g., 'Primary On-Call', 'Backup Engineer').
- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule to route to (mutually exclusive with userId). The ID of a `oneuptime_on_call_policy_schedule`.
- `order` (Number) Where this rule sits in the escalation, lowest number first. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the user to route to directly (mutually exclusive with onCallDutyPolicyScheduleId). The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
