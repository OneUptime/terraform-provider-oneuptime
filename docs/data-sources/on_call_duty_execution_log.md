---
page_title: "oneuptime_on_call_duty_execution_log Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Logs for on-call duty policy execution.
---

# oneuptime_on_call_duty_execution_log (Data Source)

Logs for on-call duty policy execution. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_on_call_duty_execution_log" "by_name" {
  name = "example-on_call_duty_execution_log"
}

data "oneuptime_on_call_duty_execution_log" "by_id" {
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
- `status` (String) Status of this execution.. Computed.
- `status_message` (String) Status message of this execution.. Computed.
- `user_notification_event_type` (String) Type of event that triggered this on-call duty policy... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `acknowledged_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `acknowledged_at` (String) A date time object.. Computed.
- `acknowledged_by_team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `last_executed_escalation_rule_order` (Number) Which escalation rule was executed?.. Computed.
- `last_executed_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_policy_execution_repeat_count` (Number) How many times did we execute this on-call policy?.. Computed.
- `schedule_gap_retry_count` (Number) How many times the current escalation rule has been re-sampled because its target schedule(s) momentarily had no on-call user... Computed.
- `triggered_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
