---
page_title: "oneuptime_user_notification_log Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Log events for user notifications
---

# oneuptime_user_notification_log (Data Source)

Log events for user notifications

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one user notification log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_user_notification_log" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_user_notification_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `acknowledged_by_user_id` (String) User ID who acknowledged this object (if this object was acknowledged by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `on_call_duty_policy_escalation_rule_id` (String) ID of your On-Call Policy Escalation Rule which belongs to this log event. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_execution_log_id` (String) ID of your On-Call Policy execution log which belongs to this log event. The ID of a `oneuptime_on_call_duty_execution_log`.
- `on_call_duty_policy_execution_log_timeline_id` (String) ID of your On-Call Policy Execution Log where this timeline event belongs. The ID of a `oneuptime_on_call_duty_execution_log_timeline` (see the data source).
- `on_call_duty_policy_id` (String) ID of your On-Call Policy which belongs to this execution log event. The ID of a `oneuptime_on_call_policy`.
- `on_call_duty_schedule_id` (String) Which schedule ID did the user belong to when the alert was sent? The ID of a `oneuptime_on_call_policy_schedule`.
- `overrided_by_user_id` (String) ID of the user this alert would have paged, when a user override sent it to this log's user instead because they were away. Empty when no override applied. The ID of a `oneuptime_user` (see the data source).
- `status` (String) Status of this execution.
- `status_message` (String) Status message of this execution.
- `triggered_by_alert_episode_id` (String) ID of the Alert Episode which triggered this on-call escalation policy. The ID of a `oneuptime_alert_episode`.
- `triggered_by_alert_id` (String) ID of the Alert which triggered this on-call escalation policy. The ID of a `oneuptime_alert`.
- `triggered_by_incident_episode_id` (String) ID of the Incident Episode which triggered this on-call escalation policy. The ID of a `oneuptime_incident_episode`.
- `triggered_by_incident_id` (String) ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_incident`.
- `user_belongs_to_team_id` (String) Which team did the user belong to when the alert was sent? The ID of a `oneuptime_team`.
- `user_id` (String) User ID who this log belongs to. The ID of a `oneuptime_user` (see the data source).
- `user_notification_event_type` (String) Notification Event Type of this execution.

### Read-Only

- `acknowledged_at` (String)
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
