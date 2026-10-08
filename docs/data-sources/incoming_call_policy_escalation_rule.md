---
page_title: "oneuptime_incoming_call_policy_escalation_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Manage escalation rules for incoming call policies that define who to call and in what order
---

# oneuptime_incoming_call_policy_escalation_rule (Data Source)

Manage escalation rules for incoming call policies that define who to call and in what order

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incoming call policy escalation rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incoming_call_policy_escalation_rule" "example" {
  name = "Example incoming call policy escalation rule"
}

# Or by id:
data "oneuptime_incoming_call_policy_escalation_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Optional description of this escalation rule.
- `escalate_after_seconds` (Number) How long, in seconds, the phone rings before the call moves on to the next rule. 20 when left out; a time below 5 or above 600 rings for 5 or 600, the limits Twilio takes.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incoming_call_policy_id` (String) ID of the parent Incoming Call Policy. The ID of a `oneuptime_incoming_call_policy`.
- `name` (String) Rule name (e.g., 'Primary On-Call', 'Backup Engineer').
- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule to route to (mutually exclusive with userId). The ID of a `oneuptime_on_call_policy_schedule`.
- `order` (Number) Where this rule sits in the escalation, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `user_id` (String) ID of the user to route to directly (mutually exclusive with onCallDutyPolicyScheduleId). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
