---
page_title: "oneuptime_service_level_objective Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Define Service Level Objectives (SLOs) with targets, compliance windows and error budgets, and track how much error budget remains.
---

# oneuptime_service_level_objective (Data Source)

Define Service Level Objectives (SLOs) with targets, compliance windows and error budgets, and track how much error budget remains. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_service_level_objective" "by_name" {
  name = "example-service_level_objective"
}

data "oneuptime_service_level_objective" "by_id" {
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
- `description` (String) Description of this Service Level Objective.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `is_enabled` (Bool) Whether this Service Level Objective is enabled. Disabled SLOs are not evaluated... Computed.
- `sli_type` (String) Type of Service Level Indicator this objective measures (Monitor Uptime or Metric).. Computed.
- `multi_monitor_mode` (String) How downtime is counted when multiple monitors are attached. 'Any Monitor Down' counts time when any monitor is down. 'Monitor Seconds Average' averages downtime across monitors... Computed.
- `monitors` (Set) Monitors whose uptime is measured by this Service Level Objective (for Monitor Uptime SLIs)... Computed.
- `monitor_labels` (Set) Monitor labels that automatically attach monitors to this SLO. Any monitor in the project carrying at least one of these labels is added to the Monitors list, and is removed again when it stops carrying any of them... Computed.
- `auto_added_monitors` (Set) Monitors that were attached to this SLO by its label rule rather than by hand. Maintained by the server... Computed.
- `downtime_monitor_statuses` (Set) List of monitor statuses that are considered as "down" for this Service Level Objective... Computed.
- `metric_query_config` (String) Query configuration for Metric SLIs: metric name, good-event predicate and optional attribute filters... Computed.
- `target_percentage` (Number) Target of this Service Level Objective as a percentage (e.g. 99.9). Must be less than 100... Computed.
- `window_type` (String) Type of compliance window for this objective (Rolling or Calendar Month).. Computed.
- `window_days` (Number) Length of the rolling compliance window in days (e.g. 7, 28, 30 or 90). Ignored for Calendar Month windows... Computed.
- `timezone` (String) IANA timezone (e.g. America/New_York) used for Calendar Month window boundaries. Defaults to UTC when not set... Computed.
- `at_risk_threshold_percentage` (Number) Percentage of remaining error budget at which the SLO status changes to At Risk. For example, 20 means the status becomes At Risk when less than 20% of the error budget remains... Computed.
- `current_sli_percentage` (Number) Current Service Level Indicator over the compliance window, as a percentage. Computed by the worker... Computed.
- `error_budget_remaining_percentage` (Number) Percentage of the error budget that remains. Can be negative when the budget is exhausted. Computed by the worker... Computed.
- `error_budget_remaining_seconds` (Number) Seconds of error budget that remain. Can be negative when the budget is exhausted. Computed by the worker... Computed.
- `error_budget_total_seconds` (Number) Total seconds of error budget for the compliance window. Computed by the worker... Computed.
- `current_burn_rate` (Number) Rate at which the error budget is currently being consumed. A burn rate of 1 exhausts the budget exactly at the end of the window. Computed by the worker... Computed.
- `slo_status` (String) Current status of this Service Level Objective (Healthy, At Risk, Budget Exhausted, Misconfigured, Paused). Computed by the worker... Computed.
- `status_change_notification_sent_at` (String) A date time object.. Computed.
- `last_evaluated_at` (String) A date time object.. Computed.
- `next_evaluation_at` (String) A date time object.. Computed.
- `last_accumulated_bucket_end_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
