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
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_burn_rate` (Number) Rate at which the error budget is currently being consumed. A burn rate of 1 exhausts the budget exactly at the end of the window. Computed by the worker.
- `current_sli_percentage` (Number) Current Service Level Indicator over the compliance window, as a percentage. Computed by the worker.
- `description` (String) Description of this Service Level Objective.
- `error_budget_remaining_percentage` (Number) Percentage of the error budget that remains. Can be negative when the budget is exhausted. Computed by the worker.
- `error_budget_remaining_seconds` (Number) Seconds of error budget that remain. Can be negative when the budget is exhausted. Computed by the worker.
- `error_budget_total_seconds` (Number) Total seconds of error budget for the compliance window. Computed by the worker.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_archived` (Boolean) Archived SLOs are hidden from lists and are not evaluated.
- `is_enabled` (Boolean) Whether this Service Level Objective is enabled. Disabled SLOs are not evaluated.
- `multi_monitor_mode` (String) How downtime is counted when multiple monitors are attached. 'Any Monitor Down' counts time when any monitor is down. 'Monitor Seconds Average' averages downtime across monitors.
- `name` (String) Name of this Service Level Objective.
- `sli_type` (String) Type of Service Level Indicator this objective measures (Monitor Uptime or Metric).
- `slo_status` (String) Current status of this Service Level Objective (Healthy, At Risk, Budget Exhausted, Misconfigured, Paused). Computed by the worker.
- `slug` (String) Friendly globally unique name for your object.
- `target_percentage` (Number) Target of this Service Level Objective as a percentage (e.g. 99.9). Must be less than 100.
- `timezone` (String) IANA timezone (e.g. America/New_York) used for Calendar Month window boundaries. Defaults to UTC when not set.
- `window_days` (Number) Length of the rolling compliance window in days (e.g. 7, 28, 30 or 90). Ignored for Calendar Month windows.
- `window_type` (String) Type of compliance window for this objective (Rolling or Calendar Month).

### Read-Only

- `archived_at` (String) When this Service Level Objective was archived.
- `auto_added_monitors` (Set of String) Monitors that were attached to this SLO by its monitor rules rather than by hand. Maintained by the server. IDs of `oneuptime_monitor` resources.
- `created_at` (String) Date and Time when the object was created.
- `downtime_monitor_statuses` (Set of String) List of monitor statuses that are considered as "down" for this Service Level Objective. IDs of `oneuptime_monitor_status` resources.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_accumulated_bucket_end_at` (String) Accumulation cursor for Metric SLIs: end of the last bucket whose good/total counts were persisted. Computed by the worker.
- `last_evaluated_at` (String) The last time this Service Level Objective was evaluated. Computed by the worker.
- `metric_query_config` (String) Query configuration for Metric SLIs: metric name, good-event predicate and optional attribute filters. A JSON value: write it with `jsonencode()`.
- `monitor_labels` (Set of String) Deprecated: superseded by SLO Monitor Rules and no longer read by the SLO engine. Existing labels were migrated into a monitor rule named "Auto-add monitors with labels". Kept only for compatibility during upgrades: labels written here to an SLO with no monitor rules are turned into that rule, and are ignored once the SLO has monitor rules. Use SLO Monitor Rules instead. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Monitors whose uptime is measured by this Service Level Objective (for Monitor Uptime SLIs). IDs of `oneuptime_monitor` resources.
- `next_evaluation_at` (String) When this Service Level Objective is next due for evaluation. Computed by the worker.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_change_notification_sent_at` (String) The last time a status-change notification was sent to owners. Computed by the worker.
- `updated_at` (String) Date and Time when the object was updated.
