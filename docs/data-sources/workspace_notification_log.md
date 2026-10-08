---
page_title: "oneuptime_workspace_notification_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all workspace activities including messages, channel creation, user invitations, and button interactions for Slack and Microsoft Teams.
---

# oneuptime_workspace_notification_log (Data Source)

Logs of all workspace activities including messages, channel creation, user invitations, and button interactions for Slack and Microsoft Teams.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one workspace notification log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_workspace_notification_log" "example" {
  workspace_type = "example-workspace-type"
}

# Or by id:
data "oneuptime_workspace_notification_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `action_type` (String) Type of workspace action performed.
- `alert_episode_id` (String) ID of Alert Episode associated with this message (if any). The ID of a `oneuptime_alert_episode`.
- `alert_id` (String) ID of Alert associated with this message (if any). The ID of a `oneuptime_alert`.
- `channel_id` (String) Channel ID where the message was sent.
- `channel_name` (String) Channel Name where the message was sent.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_episode_id` (String) ID of Incident Episode associated with this message (if any). The ID of a `oneuptime_incident_episode`.
- `incident_id` (String) ID of Incident associated with this message (if any). The ID of a `oneuptime_incident`.
- `message` (String) Content of the message.
- `monitor_id` (String) ID of Monitor associated with this message (if any). The ID of a `oneuptime_monitor`.
- `on_call_duty_policy_escalation_rule_id` (String) ID of On-Call Duty Policy Escalation Rule associated with this message (if any). The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of On-Call Duty Policy associated with this message (if any). The ID of a `oneuptime_on_call_policy`.
- `on_call_duty_policy_schedule_id` (String) ID of On-Call Duty Policy Schedule associated with this message (if any). The ID of a `oneuptime_on_call_policy_schedule`.
- `scheduled_maintenance_id` (String) ID of Scheduled Maintenance associated with this message (if any). The ID of a `oneuptime_scheduled_maintenance_event`.
- `status` (String) Status of the message.
- `status_message` (String) Status Message (if any).
- `status_page_announcement_id` (String) ID of Status Page Announcement associated with this message (if any). The ID of a `oneuptime_status_page_announcement`.
- `status_page_id` (String) ID of Status Page associated with this message (if any). The ID of a `oneuptime_status_page`.
- `team_id` (String) ID of Team associated with this message (if any). The ID of a `oneuptime_team`.
- `thread_id` (String) Thread ID of the message in the channel (if any).
- `user_id` (String) ID of User who initiated this workspace notification (if any). The ID of a `oneuptime_user` (see the data source).
- `workspace_type` (String) Type of Workspace - Slack, Microsoft Teams.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
