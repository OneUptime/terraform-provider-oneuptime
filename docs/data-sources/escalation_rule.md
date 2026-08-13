---
page_title: "oneuptime_escalation_rule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call duty escalation rule for the on-call policy.
---

# oneuptime_escalation_rule (Data Source)

Manage on-call duty escalation rule for the on-call policy. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_escalation_rule" "by_name" {
  name = "example-escalation_rule"
}

data "oneuptime_escalation_rule" "by_id" {
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
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `escalate_after_in_minutes` (Number) How long should we wait before we execute the next escalation rule?.. Computed.
- `order` (Number) Order of this rule.. Computed.
