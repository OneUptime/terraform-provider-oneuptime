---
page_title: "oneuptime_on_call_time_log Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call duty user overrides, for example if the user is on leave you can override the on-call duty policy for that user so all the alerts will be routed to the other user.
---

# oneuptime_on_call_time_log (Resource)

Manage on-call duty user overrides, for example if the user is on leave you can override the on-call duty policy for that user so all the alerts will be routed to the other user.

## Example Usage

```terraform
resource "oneuptime_on_call_time_log" "example" {
  user_id   = data.oneuptime_user.example.id
  starts_at = "2030-01-01T00:00:00Z"
}
```

## Schema

### Required

- `starts_at` (String) When does this start?
- `user_id` (String) User ID for which this log belongs. The ID of a `oneuptime_user` (see the data source).

### Optional

- `ends_at` (String) When does this end?
- `more_info` (String) More information about this log record.
- `on_call_duty_policy_escalation_rule_id` (String) ID of your On-Call Policy Escalation Rule ID where this escalation rule belongs.
- `on_call_duty_policy_id` (String) ID of your On-Call Policy where this escalation rule belongs.
- `on_call_duty_policy_schedule_id` (String) ID of your On-Call Policy Schedule where this escalation rule belongs.
- `team_id` (String) ID of your On-Call Policy Team ID where this escalation rule belongs.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call time log by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_time_log.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_time_log.example <id>
```
