---
page_title: "oneuptime_on_call_duty_execution_log_timeline Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Timeline events for on-call duty policy execution log.
---

# oneuptime_on_call_duty_execution_log_timeline (Data Source)

Timeline events for on-call duty policy execution log. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_on_call_duty_execution_log_timeline" "by_name" {
  name = "example-on_call_duty_execution_log_timeline"
}

data "oneuptime_on_call_duty_execution_log_timeline" "by_id" {
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
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `triggered_by_incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `triggered_by_alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `triggered_by_alert_episode_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `triggered_by_incident_episode_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_execution_log_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_notification_event_type` (String) Type of event that triggered this on-call duty policy... Computed.
- `alert_sent_to_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_belongs_to_team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_schedule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_message` (String) Status message of this execution timeline event.. Computed.
- `status` (String) Status of this execution timeline event.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_acknowledged` (Bool) On call duty execution log timeline is_acknowledged. Computed.
- `acknowledged_at` (String) A date time object.. Computed.
- `overrided_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
