---
page_title: "oneuptime_on_call_duty_execution_log_timeline Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Timeline events for on-call duty policy execution log.
---

# oneuptime_on_call_duty_execution_log_timeline (Data Source)

Timeline events for on-call duty policy execution log.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call duty execution log timeline may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_duty_execution_log_timeline" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
}

# Or by id:
data "oneuptime_on_call_duty_execution_log_timeline" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_sent_to_user_id` (String) ID of the user who we sent alert to. The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_acknowledged` (Boolean) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, On-Call Admin, On-Call Member, On-Call Viewer, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]
- `on_call_duty_policy_escalation_rule_id` (String) ID of your On-Call Policy Escalation Rule where this timeline event belongs. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_execution_log_id` (String) ID of your On-Call Policy Execution Log where this timeline event belongs. The ID of a `oneuptime_on_call_duty_execution_log`.
- `on_call_duty_policy_id` (String) ID of your OneUptime on-call duty policy in which this object belongs. The ID of a `oneuptime_on_call_policy`.
- `on_call_duty_schedule_id` (String) Which schedule ID did the user belong to when the alert was sent? The ID of a `oneuptime_on_call_policy_schedule`.
- `overrided_by_user_id` (String) ID of the user this alert would have paged, when a user override sent it to the user in Alert Sent To User ID instead because they were away. Empty when no override applied. The ID of a `oneuptime_user` (see the data source).
- `status` (String) Status of this execution timeline event.
- `status_message` (String) Status message of this execution timeline event.
- `triggered_by_alert_episode_id` (String) ID of your OneUptime Alert Episode in which this object belongs. The ID of a `oneuptime_alert_episode`.
- `triggered_by_alert_id` (String) ID of your OneUptime Alert in which this object belongs. The ID of a `oneuptime_alert`.
- `triggered_by_incident_episode_id` (String) ID of your OneUptime Incident Episode in which this object belongs. The ID of a `oneuptime_incident_episode`.
- `triggered_by_incident_id` (String) ID of your OneUptime Incident in which this object belongs. The ID of a `oneuptime_incident`.
- `user_belongs_to_team_id` (String) Which team ID did the user belong to when the alert was sent? The ID of a `oneuptime_team`.
- `user_notification_event_type` (String) Type of event that triggered this on-call duty policy.

### Read-Only

- `acknowledged_at` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, On-Call Admin, On-Call Member, On-Call Viewer, Read On-Call Duty Policy Execution Log Timeline], Update: [No access - you don't have permission for this operation]
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
