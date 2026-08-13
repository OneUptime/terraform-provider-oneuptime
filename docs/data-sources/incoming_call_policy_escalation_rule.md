---
page_title: "oneuptime_incoming_call_policy_escalation_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Manage escalation rules for incoming call policies that define who to call and in what order
---

# oneuptime_incoming_call_policy_escalation_rule (Data Source)

Manage escalation rules for incoming call policies that define who to call and in what order Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incoming_call_policy_escalation_rule" "by_name" {
  name = "example-incoming_call_policy_escalation_rule"
}

data "oneuptime_incoming_call_policy_escalation_rule" "by_id" {
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
- `description` (String) Optional description of this escalation rule.. Computed.
- `order` (Number) Execution order (1, 2, 3...).. Computed.
- `escalate_after_seconds` (Number) Seconds before escalating to next rule.. Computed.
- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
