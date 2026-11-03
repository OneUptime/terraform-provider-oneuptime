---
page_title: "oneuptime_escalation_rule_schedule Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add on-call schedules to your escalation rules.
---

# oneuptime_escalation_rule_schedule (Resource)

Add on-call schedules to your escalation rules.

~> **Renamed:** this resource was called `oneuptime_schedules_on_call_duty_escalation_rule` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing escalation rule schedule:

```terraform
moved {
  from = oneuptime_schedules_on_call_duty_escalation_rule.example
  to   = oneuptime_escalation_rule_schedule.example
}
```

## Example Usage

```terraform
resource "oneuptime_escalation_rule_schedule" "example" {
  on_call_duty_policy_id                 = oneuptime_on_call_policy.example.id
  on_call_duty_policy_escalation_rule_id = oneuptime_escalation_rule.example.id
}
```

## Schema

### Required

- `on_call_duty_policy_escalation_rule_id` (String) ID of the escalation rule that pages the schedule. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of the on-call policy the escalation rule is in. The ID of a `oneuptime_on_call_policy`.

### Optional

- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule the escalation rule pages: whoever is on call in it. The ID of a `oneuptime_on_call_policy_schedule`.
- `page_next_in_rotation` (Boolean) When true, the escalation rule pages the next person in this schedule's rotation - the person after whoever is on call, passing over anyone who is away - instead of whoever is on call. Use it on the level after the one that pages the schedule, so a page nobody acknowledges goes to the backup from the same rotation. Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing escalation rule schedule by its escalation rule and its `on_call_duty_policy_schedule_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_escalation_rule_schedule.example
  id = "<on_call_duty_policy_escalation_rule_id>/<on_call_duty_policy_schedule_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_escalation_rule_schedule.example <on_call_duty_policy_escalation_rule_id>/<on_call_duty_policy_schedule_id>
```

Its own ID works too.
