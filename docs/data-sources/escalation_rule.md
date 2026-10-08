---
page_title: "oneuptime_escalation_rule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call duty escalation rule for the on-call policy.
---

# oneuptime_escalation_rule (Data Source)

Manage on-call duty escalation rule for the on-call policy.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one escalation rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_escalation_rule" "example" {
  name = "Example escalation rule"
}

# Or by id:
data "oneuptime_escalation_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `escalate_after_in_minutes` (Number) How long should we wait before we execute the next escalation rule?
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name of this object.
- `on_call_duty_policy_id` (String) ID of your On-Call Policy where this escalation rule belongs. The ID of a `oneuptime_on_call_policy`.
- `order` (Number) Order of this rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
