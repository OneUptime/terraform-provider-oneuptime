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
  on_call_duty_policy_id = "123e4567-e89b-12d3-a456-426614174000"
  status = "Example short text"
  status_message = "This is an example of longer text content that might be stored in this field."
  user_notification_event_type = "Example short text"
}
```

## Schema

### Required

- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID..
- `status` (String) Status of this execution..
- `status_message` (String) Status message of this execution..
- `user_notification_event_type` (String) Type of event that triggered this on-call duty policy...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `triggered_by_incident_id` (String) A unique identifier for an object, represented as a UUID..
- `triggered_by_alert_id` (String) A unique identifier for an object, represented as a UUID..
- `triggered_by_alert_episode_id` (String) A unique identifier for an object, represented as a UUID..
- `triggered_by_incident_episode_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `acknowledged_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `acknowledged_at` (String) A date time object..
- `acknowledged_by_team_id` (String) A unique identifier for an object, represented as a UUID..
- `last_executed_escalation_rule_order` (Number) Which escalation rule was executed?..
- `last_escalation_rule_executed_at` (String) A date time object..
- `last_executed_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID..
- `execute_next_escalation_rule_in_minutes` (Number) How many minutes should we wait before executing the next escalation rule?..
- `on_call_policy_execution_repeat_count` (Number) How many times did we execute this on-call policy?..
- `triggered_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `schedule_gap_retry_count` (Number) How many times the current escalation rule has been re-sampled because its target schedule(s) momentarily had no on-call user...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_on_call_duty_execution_log.example <id>
```
