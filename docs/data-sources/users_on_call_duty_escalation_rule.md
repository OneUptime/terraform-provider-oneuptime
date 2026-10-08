---
page_title: "oneuptime_users_on_call_duty_escalation_rule Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Manage on-call duty escalation rule for the on-call policy.
---

# oneuptime_users_on_call_duty_escalation_rule (Data Source)

Manage on-call duty escalation rule for the on-call policy.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one users on call duty escalation rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_users_on_call_duty_escalation_rule" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
}

# Or by id:
data "oneuptime_users_on_call_duty_escalation_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `on_call_duty_policy_escalation_rule_id` (String) ID of your On-Call Policy Escalation Rule where this user belongs. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of your On-Call Policy where this escalation rule belongs. The ID of a `oneuptime_on_call_policy`.
- `user_id` (String) ID of the user who is in this escalation rule. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
