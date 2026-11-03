---
page_title: "oneuptime_escalation_rule_user Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add users to your escalation rules.
---

# oneuptime_escalation_rule_user (Data Source)

Add users to your escalation rules.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one escalation rule user may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_users_on_call_duty_escalation_rule` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_escalation_rule_user" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
  user_email             = "alice@example.com"
}

# Or by id:
data "oneuptime_escalation_rule_user" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `on_call_duty_policy_escalation_rule_id` (String) ID of the escalation rule that pages the person. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of the on-call policy the escalation rule is in. The ID of a `oneuptime_on_call_policy`.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the person this escalation rule pages. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
