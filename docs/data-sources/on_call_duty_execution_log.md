---
page_title: "oneuptime_on_call_duty_execution_log Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Each time an on-call policy runs: what started it, how far it got, and who it paged.
---

# oneuptime_on_call_duty_execution_log (Data Source)

Each time an on-call policy runs: what started it, how far it got, and who it paged.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call duty execution log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_duty_execution_log" "example" {
  created_by_actor_type = "User"
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
- `ai_hold_outcome` (String) For a run OneUptime AI held while it investigated and fixed: Holding while it waits, then how the hold ended - Resolved by AI, Recovered on its own, Handled by a person, Handed off, Hold time ran out or Page now pressed. Empty for every other run.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_observer_page` (Boolean) Whether this run pages its people as observers: they are kept informed, and acknowledging their page only says they have seen it - the incident or alert is not acknowledged, and nobody else's page stops.
- `last_executed_escalation_rule_id` (String) ID of your On-Call Policy Last Executed Escalation Rule. The ID of a `oneuptime_escalation_rule`.
- `last_executed_escalation_rule_order` (Number) Which escalation rule was executed?
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `on_call_duty_policy_id` (String) ID of the On-Call Policy this execution log is a run of. Required when you create one. Empty for a direct page: people, a team or an on-call schedule paged without an on-call policy. The ID of a `oneuptime_on_call_policy`.
- `on_call_policy_execution_repeat_count` (Number) How many times did we execute this on-call policy?
- `page_again_number` (Number) For a run that pages again because its record stayed acknowledged without being resolved: which time, from 1. Empty for every other run.
- `page_message` (String) What the person who paged from Slack or Microsoft Teams wrote to the people this page reaches.
- `paging_started_early_reason` (String) Why this run started paging before the working hours it waited for began: its record came to match an on-call rule that pages right away. Empty for every other run.
- `schedule_gap_retry_count` (Number) How many times the current escalation rule has been re-sampled because its target schedule(s) momentarily had no on-call user.
- `status` (String) Status of this execution.
- `status_message` (String) Status message of this execution.
- `test_page_fire_drill_id` (String) For a test started by a monitor's fire drill: the drill, the same on the test of every on-call policy it started. Empty for every other run.
- `test_page_fire_drill_monitor_id` (String) For a test started by a monitor's fire drill: ID of the monitor whose failure the drill rehearsed. Empty for every other run, and once the monitor is deleted. The ID of a `oneuptime_monitor`.
- `test_page_record_kind` (String) For a test page: the kind of record it stands in for - Incident or Alert - whose On-Call Rules it runs. Empty for every other run.
- `test_page_severity_id` (String) For a test page: the incident or alert severity whose On-Call Rules it runs. Empty for every other run.
- `test_page_walks_escalation` (Boolean) For a test of an on-call policy: whether it walks every escalation level, each after a short wait, rather than paging the first level only.
- `triggered_by_alert_episode_id` (String) ID of the alert episode which triggered this on-call escalation policy. The ID of a `oneuptime_alert_episode`.
- `triggered_by_alert_id` (String) ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_alert`.
- `triggered_by_incident_episode_id` (String) ID of the incident episode which triggered this on-call escalation policy. The ID of a `oneuptime_incident_episode`.
- `triggered_by_incident_id` (String) ID of the incident which triggered this on-call escalation policy. The ID of a `oneuptime_incident`.
- `triggered_by_user_id` (String) User ID who triggered this on-call policy. The ID of a `oneuptime_user` (see the data source).
- `triggered_by_workflow_id` (String) ID of the workflow whose step started this page, when a workflow did.
- `triggered_by_workflow_name` (String) Name of the workflow whose step started this page, as it was then.
- `user_notification_event_type` (String) Type of event that triggered this on-call policy.

### Read-Only

- `acknowledged_at` (String) When was this policy execution acknowledged?
- `acknowledgement_timeout_started_at` (String) When the acknowledgement timeout started counting for this run's record: its last acknowledgement. Empty when it is not counting.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `paging_starts_at` (String) When this run starts paging, when it waited to: an on-call rule that pages only during working hours held it until they began. Empty when it paged straight away.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `test_page_fire_drill_failure` (String) For a test started by a monitor's fire drill: the failure the drill rehearsed, as its dialog named it - the criteria's name, and how long an alert stays open before an alert incident rule declares the incident tested. Empty for every other run. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
