---
page_title: "oneuptime_incoming_call_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Parent log for each incoming call instance. Groups all escalation attempts together.
---

# oneuptime_incoming_call_log (Data Source)

Parent log for each incoming call instance. Groups all escalation attempts together. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incoming_call_log" "by_name" {
  name = "example-incoming_call_log"
}

data "oneuptime_incoming_call_log" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incoming_call_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `caller_phone_number` (String) Phone object. Computed.
- `routing_phone_number` (String) Phone object. Computed.
- `call_provider_call_id` (String) Call provider's call identifier.. Computed.
- `status` (String) Current status of the incoming call.. Computed.
- `status_message` (String) Additional status information.. Computed.
- `call_duration_in_seconds` (Number) Total call duration in seconds.. Computed.
- `call_cost_in_usd_cents` (Number) Total cost for this call in USD cents.. Computed.
- `incoming_call_cost_in_usd_cents` (Number) Cost for incoming leg in USD cents.. Computed.
- `outgoing_call_cost_in_usd_cents` (Number) Cost for all forwarding attempts in USD cents.. Computed.
- `started_at` (String) A date time object.. Computed.
- `ended_at` (String) A date time object.. Computed.
- `answered_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `current_escalation_rule_order` (Number) The current escalation rule order being processed.. Computed.
- `repeat_count` (Number) Number of times the policy has been repeated.. Computed.
