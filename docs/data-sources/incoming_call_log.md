---
page_title: "oneuptime_incoming_call_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Parent log for each incoming call instance. Groups all escalation attempts together.
---

# oneuptime_incoming_call_log (Data Source)

Parent log for each incoming call instance. Groups all escalation attempts together.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming call log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_call_log" "example" {
  incoming_call_policy_id = oneuptime_incoming_call_policy.example.id
}

# Or by id:
data "oneuptime_incoming_call_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `answered_by_user_id` (String) User ID who answered the call. The ID of a `oneuptime_user` (see the data source).
- `call_cost_in_usd_cents` (Number) Total cost for this call in USD cents.
- `call_duration_in_seconds` (Number) Total call duration in seconds.
- `call_provider_call_id` (String) Call provider's call identifier.
- `current_escalation_rule_order` (Number) The current escalation rule order being processed.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incoming_call_cost_in_usd_cents` (Number) Cost for incoming leg in USD cents.
- `incoming_call_policy_id` (String) ID of the Incoming Call Policy. The ID of a `oneuptime_incoming_call_policy`.
- `outgoing_call_cost_in_usd_cents` (Number) Cost for all forwarding attempts in USD cents.
- `repeat_count` (Number) Number of times the policy has been repeated.
- `status` (String) Current status of the incoming call.
- `status_message` (String) Additional status information.

### Read-Only

- `caller_phone_number` (String) Incoming caller's phone number.
- `created_at` (String) Date and Time when the object was created.
- `ended_at` (String) When the call ended.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `routing_phone_number` (String) The routing number that was called.
- `started_at` (String) When the call started.
- `updated_at` (String) Date and Time when the object was updated.
