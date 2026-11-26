---
page_title: "oneuptime_service_level_objective Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Define Service Level Objectives (SLOs) with targets, compliance windows and error budgets, and track how much error budget remains.
---

# oneuptime_service_level_objective (Data Source)

Define Service Level Objectives (SLOs) with targets, compliance windows and error budgets, and track how much error budget remains.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one service level objective may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_service_level_objective" "example" {
  name = "Example service level objective"
}

# Or by id:
data "oneuptime_service_level_objective" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `at_risk_threshold_percentage` (Number) Percentage of remaining error budget at which the SLO status changes to At Risk. For example, 20 means the status becomes At Risk when less than 20% of the error budget remains.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_burn_rate` (Number) Rate at which the error budget is currently being consumed. A burn rate of 1 exhausts the budget exactly at the end of the window. Computed by the worker.
- `current_error_budget_policy_tier_id` (String) ID of the error budget policy tier this Service Level Objective is in. Empty while its error budget is above every tier, or when it has no error budget policy. Computed by the worker.
- `current_error_budget_policy_tier_name` (String) Name of the error budget policy tier this Service Level Objective is in. Empty while its error budget is above every tier, or when it has no error budget policy. Computed by the worker.
- `current_sli_percentage` (Number) Current Service Level Indicator over the compliance window, as a percentage. Computed by the worker.
- `deployments_allowed` (Boolean) False while this Service Level Objective is in a tier of its error budget policy that freezes deployments, and true otherwise, so a CI/CD pipeline can check it before it releases. Kept as it is while the SLO is not measured. Computed by the worker.
- `description` (String) Description of this Service Level Objective.
- `error_budget_remaining_events` (Number) For request SLIs: how many more bad requests the target allows in the compliance window. Negative when the budget is exhausted. Computed by the worker.
- `error_budget_remaining_percentage` (Number) Percentage of the error budget that remains. Can be negative when the budget is exhausted. Computed by the worker.
- `error_budget_remaining_seconds` (Number) Seconds of error budget that remain. Can be negative when the budget is exhausted. Computed by the worker.
- `error_budget_total_events` (Number) For request SLIs: how many bad requests the target allows in the compliance window. Computed by the worker.
- `error_budget_total_seconds` (Number) Total seconds of error budget for the compliance window. Computed by the worker.
- `good_event_count` (Number) For request SLIs: the requests in the compliance window that counted as good (did not fail, or were fast enough). Computed by the worker.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_archived` (Boolean) Archived SLOs are hidden from lists and are not evaluated.
- `is_enabled` (Boolean) Whether this Service Level Objective is enabled. Disabled SLOs are not evaluated.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `latency_threshold_in_milliseconds` (Number) For the Request Latency SLI type: a request served in this many milliseconds or less counts as good. A whole number from 1 to 3600000.
- `multi_monitor_mode` (String) How downtime is counted when multiple monitors are attached. 'Any Monitor Down' counts time when any monitor is down. 'Monitor Seconds Average' averages downtime across monitors.
- `name` (String) Name of this Service Level Objective.
- `request_filter` (String) For request SLI types: count only the requests that match this traces search of span attributes and the span's name, such as @http.route:/checkout* -@http.route:/health. Empty counts every request of the SLO's services.
- `sli_type` (String) What this objective measures: Monitor Uptime (the share of time its monitors are up), Request Success (the share of its services' requests that did not fail) or Request Latency (the share of its services' requests served within the latency threshold). Metric is reserved and not measured yet.
- `slo_status` (String) Current status of this Service Level Objective (Healthy, At Risk, Budget Exhausted, Misconfigured, Paused). Computed by the worker.
- `slug` (String) Friendly globally unique name for your object.
- `target_percentage` (Number) Target of this Service Level Objective as a percentage (e.g. 99.9). Must be less than 100.
- `timezone` (String) IANA timezone (e.g. America/New_York) used for Calendar Month window boundaries. Defaults to UTC when not set.
- `total_event_count` (Number) For request SLIs: every request in the compliance window that this Service Level Objective measures. Computed by the worker.
- `window_days` (Number) Length of the rolling compliance window in days (e.g. 7, 28, 30 or 90). Ignored for Calendar Month windows.
- `window_type` (String) Type of compliance window for this objective (Rolling or Calendar Month).

### Read-Only

- `archived_at` (String) When this Service Level Objective was archived.
- `auto_added_monitors` (Set of String) Monitors that were attached to this SLO by its monitor rules rather than by hand. Maintained by the server. IDs of `oneuptime_monitor` resources.
- `created_at` (String) Date and Time when the object was created.
- `current_error_budget_policy_tier_entered_at` (String) When this Service Level Objective entered the error budget policy tier it is in. Empty while it is in no tier. Computed by the worker.
- `deployments_frozen_at` (String) When deployments of this Service Level Objective were frozen: when it entered a tier of its error budget policy that freezes deployments, or the tier it is in started freezing them. Empty while deployments are allowed. Computed by the worker.
- `downtime_monitor_statuses` (Set of String) List of monitor statuses that are considered as "down" for this Service Level Objective. IDs of `oneuptime_monitor_status` resources.
- `first_accumulated_bucket_start_at` (String) For request SLIs: the start of the oldest 5-minute step whose good and total requests were counted, so the time this objective's numbers cover. Computed by the worker.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.
- `last_accumulated_bucket_end_at` (String) For request SLIs: the end of the newest 5-minute step whose good and total requests were counted. Computed by the worker.
- `last_evaluated_at` (String) The last time this Service Level Objective was evaluated. Computed by the worker.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `metric_query_config` (String) Query configuration for Metric SLIs: metric name, good-event predicate and optional attribute filters. A JSON value: read it with `jsondecode()`.
- `monitor_labels` (Set of String) Deprecated: superseded by SLO Monitor Rules and no longer read by the SLO engine. Existing labels were migrated into a monitor rule named "Auto-add monitors with labels". Kept only for compatibility during upgrades: labels written here to an SLO with no monitor rules are turned into that rule, and are ignored once the SLO has monitor rules. Use SLO Monitor Rules instead. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Monitors whose uptime is measured by this Service Level Objective (for Monitor Uptime SLIs). IDs of `oneuptime_monitor` resources.
- `next_evaluation_at` (String) When this Service Level Objective is next due for evaluation. Computed by the worker.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `services` (Set of String) Services whose requests this Service Level Objective measures, for the Request Success and Request Latency SLI types. Their requests are counted from their traces. IDs of `oneuptime_service` resources.
- `status_change_notification_sent_at` (String) The last time a status-change notification was sent to owners. Computed by the worker.
- `updated_at` (String) Date and Time when the object was updated.
