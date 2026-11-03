---
page_title: "oneuptime_escalation_rule_schedule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add on-call schedules to your escalation rules.
---

# oneuptime_escalation_rule_schedule (Data Source)

Add on-call schedules to your escalation rules.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one escalation rule schedule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_schedules_on_call_duty_escalation_rule` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_escalation_rule_schedule" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
}

# Or by id:
data "oneuptime_escalation_rule_schedule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `on_call_duty_policy_escalation_rule_id` (String) ID of the escalation rule that pages the schedule. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of the on-call policy the escalation rule is in. The ID of a `oneuptime_on_call_policy`.
- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule the escalation rule pages: whoever is on call in it. The ID of a `oneuptime_on_call_policy_schedule`.
- `page_next_in_rotation` (Boolean) When true, the escalation rule pages the next person in this schedule's rotation - the person after whoever is on call, passing over anyone who is away - instead of whoever is on call. Use it on the level after the one that pages the schedule, so a page nobody acknowledges goes to the backup from the same rotation.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
