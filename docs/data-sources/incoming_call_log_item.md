---
page_title: "oneuptime_incoming_call_log_item Data Source - oneuptime"
subcategory: "Other"
description: |-
  Child log for each escalation attempt / user ring within a call.
---

# oneuptime_incoming_call_log_item (Data Source)

Child log for each escalation attempt / user ring within a call.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming call log item may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_call_log_item" "example" {
  incoming_call_log_id = data.oneuptime_incoming_call_log.example.id
}

# Or by id:
data "oneuptime_incoming_call_log_item" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `call_cost_in_usd_cents` (Number) Cost for this dial attempt in USD cents.
- `dial_duration_in_seconds` (Number) How long this dial lasted in seconds.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incoming_call_log_id` (String) ID of the parent Incoming Call Log. The ID of a `oneuptime_incoming_call_log` (see the data source).
- `incoming_call_policy_escalation_rule_id` (String) ID of the escalation rule used. The ID of a `oneuptime_incoming_call_policy_escalation_rule`.
- `is_answered` (Boolean) Whether this user answered the call.
- `status` (String) Status of this dial attempt.
- `status_message` (String) Additional status information.
- `user_id` (String) User ID who was called. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ended_at` (String) When dial ended.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `started_at` (String) When dial started.
- `updated_at` (String) Date and Time when the object was updated.
- `user_phone_number` (String) Phone number that was dialed.
