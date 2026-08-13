---
page_title: "oneuptime_workspace_notification_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all workspace activities including messages, channel creation, user invitations, and button interactions for Slack and Microsoft Teams.
---

# oneuptime_workspace_notification_log (Data Source)

Logs of all workspace activities including messages, channel creation, user invitations, and button interactions for Slack and Microsoft Teams. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_workspace_notification_log" "by_name" {
  name = "example-workspace_notification_log"
}

data "oneuptime_workspace_notification_log" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `workspace_type` (String) Type of Workspace - Slack, Microsoft Teams.. Computed.
- `channel_id` (String) Channel ID where the message was sent.. Computed.
- `channel_name` (String) Channel Name where the message was sent.. Computed.
- `thread_id` (String) Thread ID of the message in the channel (if any).. Computed.
- `message` (String) Content of the message.. Computed.
- `status_message` (String) Status Message (if any).. Computed.
- `status` (String) Status of the message.. Computed.
- `action_type` (String) Type of workspace action performed.. Computed.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_episode_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_episode_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `scheduled_maintenance_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_announcement_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
