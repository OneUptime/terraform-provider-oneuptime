---
page_title: "oneuptime_on_call_duty_policy_user_owner Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add users as owners to your onCallDutyPolicys.
---

# oneuptime_on_call_duty_policy_user_owner (Data Source)

Add users as owners to your onCallDutyPolicys. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_on_call_duty_policy_user_owner" "by_name" {
  name = "example-on_call_duty_policy_user_owner"
}

data "oneuptime_on_call_duty_policy_user_owner" "by_id" {
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
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_owner_notified` (Bool) Are owners notified of this resource ownership?.. Computed.
