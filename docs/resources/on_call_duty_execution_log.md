---
page_title: "oneuptime_on_call_duty_execution_log Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Logs for on-call duty policy execution.
---

# oneuptime_on_call_duty_execution_log (Resource)

Logs for on-call duty policy execution.

## Example Usage

```terraform
resource "oneuptime_on_call_duty_execution_log" "example" {
  on_call_duty_policy_id       = oneuptime_on_call_policy.example.id
  status                       = "Example short text"
  status_message               = "This is an example of longer text content that might be stored in this field."
  user_notification_event_type = "Example short text"
}
```

## Schema

### Required

- `on_call_duty_policy_id` (String) ID of your On-Call Policy which belongs to this execution log event. The ID of a `oneuptime_on_call_policy`.
- `status` (String) Status of this execution.
- `status_message` (String) Status message of this execution.
- `user_notification_event_type` (String) Type of event that triggered this on-call duty policy.

### Optional

- `acknowledged_at` (String) When was this policy execution acknowledged?
- `acknowledged_by_team_id` (String) Team ID who acknowledged this object (if this object was acknowledged by a Team). The ID of a `oneuptime_team`.
- `execute_next_escalation_rule_in_minutes` (Number) How many minutes should we wait before executing the next escalation rule?
- `last_escalation_rule_executed_at` (String) When was the escalation rule executed?
- `last_executed_escalation_rule_id` (String) ID of your On-Call Policy Last Executed Escalation Rule. The ID of a `oneuptime_escalation_rule`.
- `last_executed_escalation_rule_order` (Number) Which escalation rule was executed?
- `on_call_policy_execution_repeat_count` (Number) How many times did we execute this on-call policy? Defaults to `1`.
- `triggered_by_alert_episode_id` (String) ID of the alert episode which triggered this on-call escalation policy. The ID of a `oneuptime_alert_episode`.
- `triggered_by_alert_id` (String) ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_alert`.
- `triggered_by_incident_episode_id` (String) ID of the incident episode which triggered this on-call escalation policy. The ID of a `oneuptime_incident_episode`.
- `triggered_by_incident_id` (String) ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_incident`.

### Read-Only

- `acknowledged_by_user_id` (String) User ID who acknowledged this object (if this object was acknowledged by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `schedule_gap_retry_count` (Number) How many times the current escalation rule has been re-sampled because its target schedule(s) momentarily had no on-call user.
- `triggered_by_user_id` (String) User ID who triggered this on-call policy. The ID of a `oneuptime_user` (see the data source).
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call duty execution log by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_duty_execution_log.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_duty_execution_log.example <id>
```
