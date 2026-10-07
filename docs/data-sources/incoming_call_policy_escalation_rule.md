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
- `order` (Number) Where this rule sits in the escalation, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them... Computed.
- `escalate_after_seconds` (Number) How long, in seconds, the phone rings before the call moves on to the next rule. 20 when left out; a time below 5 or above 600 rings for 5 or 600, the limits Twilio takes... Computed.
- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
