---
page_title: "oneuptime_incoming_call_log_item Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Child log for each escalation attempt / user ring within a call.
---

# oneuptime_incoming_call_log_item (Data Source)

Child log for each escalation attempt / user ring within a call.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming call log item may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_call_log_item" "example" {
  last_modified_by_user_id = data.oneuptime_user.example.id
  user_email               = "alice@example.com"
}

# Or by id:
data "oneuptime_incoming_call_log_item" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `call_cost_in_usd_cents` (Number) Cost for this dial attempt in USD cents.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `dial_duration_in_seconds` (Number) How long this dial lasted in seconds.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incoming_call_log_id` (String) ID of the parent Incoming Call Log. The ID of a `oneuptime_incoming_call_log` (see the data source).
- `incoming_call_policy_escalation_rule_id` (String) ID of the escalation rule used. The ID of a `oneuptime_incoming_call_policy_escalation_rule`.
- `incoming_call_policy_id` (String) ID of the incoming call policy whose number took the call this attempt is part of. The ID of a `oneuptime_incoming_call_policy`.
- `is_answered` (Boolean) Whether this user answered the call.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `status` (String) Status of this dial attempt.
- `status_message` (String) Additional status information.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) User ID who was called. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ended_at` (String) When dial ended.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `started_at` (String) When dial started.
- `updated_at` (String) Date and Time when the object was updated.
- `user_phone_number` (String) Phone number that was dialed.
