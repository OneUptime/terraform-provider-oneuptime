---
page_title: "oneuptime_on_call_duty_execution_log Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Logs for on-call duty policy execution.
---

# oneuptime_on_call_duty_execution_log (Data Source)

Logs for on-call duty policy execution.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call duty execution log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_duty_execution_log" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
}

# Or by id:
data "oneuptime_on_call_duty_execution_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `acknowledged_by_team_id` (String) Team ID who acknowledged this object (if this object was acknowledged by a Team). The ID of a `oneuptime_team`.
- `acknowledged_by_user_id` (String) User ID who acknowledged this object (if this object was acknowledged by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_executed_escalation_rule_id` (String) ID of your On-Call Policy Last Executed Escalation Rule. The ID of a `oneuptime_escalation_rule`.
- `last_executed_escalation_rule_order` (Number) Which escalation rule was executed?
- `on_call_duty_policy_id` (String) ID of the On-Call Policy this execution log is a run of. Required when you create one. Empty for a direct page: people, a team or an on-call schedule paged without an on-call policy. The ID of a `oneuptime_on_call_policy`.
- `on_call_policy_execution_repeat_count` (Number) How many times did we execute this on-call policy?
- `page_again_number` (Number) For a run that pages again because its record stayed acknowledged without being resolved: which time, from 1. Empty for every other run.
- `paging_started_early_reason` (String) Why this run started paging before the working hours it waited for began: its record came to match an on-call rule that pages right away. Empty for every other run.
- `schedule_gap_retry_count` (Number) How many times the current escalation rule has been re-sampled because its target schedule(s) momentarily had no on-call user.
- `status` (String) Status of this execution.
- `status_message` (String) Status message of this execution.
- `triggered_by_alert_episode_id` (String) ID of the alert episode which triggered this on-call escalation policy. The ID of a `oneuptime_alert_episode`.
- `triggered_by_alert_id` (String) ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_alert`.
- `triggered_by_incident_episode_id` (String) ID of the incident episode which triggered this on-call escalation policy. The ID of a `oneuptime_incident_episode`.
- `triggered_by_incident_id` (String) ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_incident`.
- `triggered_by_user_id` (String) User ID who triggered this on-call policy. The ID of a `oneuptime_user` (see the data source).
- `triggered_by_workflow_id` (String) ID of the workflow whose step started this page, when a workflow did.
- `triggered_by_workflow_name` (String) Name of the workflow whose step started this page, as it was then.
- `user_notification_event_type` (String) Type of event that triggered this on-call duty policy.

### Read-Only

- `acknowledged_at` (String) When was this policy execution acknowledged?
- `acknowledgement_timeout_started_at` (String) When the acknowledgement timeout started counting for this run's record: its last acknowledgement. Empty when it is not counting.
- `created_at` (String) Date and Time when the object was created.
- `paging_starts_at` (String) When this run starts paging, when it waited to: an on-call rule that pages only during working hours held it until they began. Empty when it paged straight away.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
