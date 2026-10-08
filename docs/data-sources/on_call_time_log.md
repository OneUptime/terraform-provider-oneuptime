---
page_title: "oneuptime_on_call_time_log Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call duty user overrides, for example if the user is on leave you can override the on-call duty policy for that user so all the alerts will be routed to the other user.
---

# oneuptime_on_call_time_log (Data Source)

Manage on-call duty user overrides, for example if the user is on leave you can override the on-call duty policy for that user so all the alerts will be routed to the other user.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call time log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_time_log" "example" {
  on_call_duty_policy_id = "example-on-call-duty-policy-id"
}

# Or by id:
data "oneuptime_on_call_time_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_info` (String) More information about this log record.
- `on_call_duty_policy_escalation_rule_id` (String) ID of your On-Call Policy Escalation Rule ID where this escalation rule belongs.
- `on_call_duty_policy_id` (String) ID of your On-Call Policy where this escalation rule belongs.
- `on_call_duty_policy_schedule_id` (String) ID of your On-Call Policy Schedule where this escalation rule belongs.
- `team_id` (String) ID of your On-Call Policy Team ID where this escalation rule belongs.
- `user_id` (String) User ID for which this log belongs. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ends_at` (String) When does this end?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `starts_at` (String) When does this start?
- `updated_at` (String) Date and Time when the object was updated.
