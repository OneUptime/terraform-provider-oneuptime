---
page_title: "oneuptime_user_on_call_log_timeline Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Timeline events for user on-call log.
---

# oneuptime_user_on_call_log_timeline (Data Source)

Timeline events for user on-call log.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one user on call log timeline may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_user_on_call_log_timeline" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_user_on_call_log_timeline" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_acknowledged` (Boolean)
- `on_call_duty_policy_escalation_rule_id` (String) ID of your On-Call Policy Escalation Rule where this timeline event belongs. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_execution_log_id` (String) ID of your On-Call Policy Execution Log Timeline where this timeline event belongs. The ID of a `oneuptime_on_call_duty_execution_log`.
- `on_call_duty_policy_execution_log_timeline_id` (String) ID of your On-Call Policy Execution Log where this timeline event belongs. The ID of a `oneuptime_on_call_duty_execution_log_timeline` (see the data source).
- `on_call_duty_policy_id` (String) ID of your OneUptime on-call duty policy in which this object belongs. The ID of a `oneuptime_on_call_policy`.
- `status` (String) Status of this execution timeline event.
- `status_message` (String) Status message of this execution timeline event.
- `triggered_by_alert_episode_id` (String) ID of your OneUptime Alert Episode in which this object belongs. The ID of a `oneuptime_alert_episode`.
- `triggered_by_alert_id` (String) ID of your OneUptime Alert in which this object belongs. The ID of a `oneuptime_alert`.
- `triggered_by_incident_episode_id` (String) ID of your OneUptime Incident Episode in which this object belongs. The ID of a `oneuptime_incident_episode`.
- `triggered_by_incident_id` (String) ID of your OneUptime Incident in which this object belongs. The ID of a `oneuptime_incident`.
- `user_belongs_to_team_id` (String) Which team did the user belong to when the alert was sent? The ID of a `oneuptime_team`.
- `user_call_id` (String) ID of User Call in which this object belongs.
- `user_email_id` (String) ID of User Email in which this object belongs.
- `user_id` (String) User ID who this log belongs to. The ID of a `oneuptime_user` (see the data source).
- `user_microsoft_teams_id` (String) ID of User Microsoft Teams in which this object belongs.
- `user_notification_event_type` (String) Notification Event Type of this execution.
- `user_notification_log_id` (String) ID of your OneUptime User Notification Log in which this object belongs. The ID of a `oneuptime_user_notification_log` (see the data source).
- `user_notification_rule_id` (String) ID of your OneUptime User Notification Rule in which this object belongs.
- `user_push_id` (String) ID of User Push in which this object belongs.
- `user_slack_id` (String) ID of User Slack in which this object belongs.
- `user_sms_id` (String) ID of User SMS in which this object belongs.
- `user_telegram_id` (String) ID of User Telegram in which this object belongs.
- `user_webhook_id` (String) ID of User Webhook in which this object belongs.
- `user_whats_app_id` (String) ID of User WhatsApp in which this object belongs.

### Read-Only

- `acknowledged_at` (String)
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
