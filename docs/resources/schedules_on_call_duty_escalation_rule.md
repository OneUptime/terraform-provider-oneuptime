---
page_title: "oneuptime_schedules_on_call_duty_escalation_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Manage schedules for on-call duty policy escalation rules.
---

# oneuptime_schedules_on_call_duty_escalation_rule (Resource)

Manage schedules for on-call duty policy escalation rules.

## Example Usage

```terraform
resource "oneuptime_schedules_on_call_duty_escalation_rule" "example" {
  on_call_duty_policy_id                 = oneuptime_on_call_policy.example.id
  on_call_duty_policy_escalation_rule_id = oneuptime_escalation_rule.example.id
}
```

## Schema

### Required

- `on_call_duty_policy_escalation_rule_id` (String) ID of your On-Call Policy Escalation Rule where this user belongs. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of your On-Call Policy where this escalation rule belongs. The ID of a `oneuptime_on_call_policy`.

### Optional

- `on_call_duty_policy_schedule_id` (String) ID of the on call schedule which is in this escalation rule. The ID of a `oneuptime_on_call_policy_schedule`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing schedules on call duty escalation rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_schedules_on_call_duty_escalation_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_schedules_on_call_duty_escalation_rule.example <id>
```
