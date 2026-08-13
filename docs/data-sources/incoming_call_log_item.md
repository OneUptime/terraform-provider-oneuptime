---
page_title: "oneuptime_incoming_call_log_item Data Source - oneuptime"
subcategory: "Other"
description: |-
  Child log for each escalation attempt / user ring within a call.
---

# oneuptime_incoming_call_log_item (Data Source)

Child log for each escalation attempt / user ring within a call. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incoming_call_log_item" "by_name" {
  name = "example-incoming_call_log_item"
}

data "oneuptime_incoming_call_log_item" "by_id" {
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
- `incoming_call_log_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incoming_call_policy_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_phone_number` (String) Phone object. Computed.
- `status` (String) Status of this dial attempt.. Computed.
- `status_message` (String) Additional status information.. Computed.
- `dial_duration_in_seconds` (Number) How long this dial lasted in seconds.. Computed.
- `call_cost_in_usd_cents` (Number) Cost for this dial attempt in USD cents.. Computed.
- `started_at` (String) A date time object.. Computed.
- `ended_at` (String) A date time object.. Computed.
- `is_answered` (Bool) Whether this user answered the call.. Computed.
