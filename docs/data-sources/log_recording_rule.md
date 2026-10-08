---
page_title: "oneuptime_log_recording_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Derived metrics computed every minute from logs: a count of matching logs, or the sum, average, min, max or a percentile of a numeric log attribute, optionally split by log attributes. Results are written into the metric store as a new series.
---

# oneuptime_log_recording_rule (Data Source)

Derived metrics computed every minute from logs: a count of matching logs, or the sum, average, min, max or a percentile of a numeric log attribute, optionally split by log attributes. Results are written into the metric store as a new series.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one log recording rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_log_recording_rule" "example" {
  name = "Example log recording rule"
}

# Or by id:
data "oneuptime_log_recording_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this rule. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What this recording rule computes and why.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is evaluated by the log recording rule cron, every minute.
- `name` (String) Friendly name for this rule.
- `output_metric_name` (String) Name of the new metric this rule writes (e.g. sdwan.gateway.latency.ms). Leave it out and it is made from the rule's name - SD-WAN gateway latency becomes sd_wan_gateway_latency, with _2, _3 and so on added when another recording rule of the project already writes it. Keep it unique per project.

### Read-Only

- `computed_until` (String) The end of the last minute this rule has written into its metric. Null means it has not run yet.
- `created_at` (String) Date and Time when the object was created.
- `definition` (String) Which logs count (filter: telemetryServiceIds, severityTexts, body, attributeFilters), how they are aggregated (aggregationType: Count, or Sum / Avg / Min / Max / P50 / P75 / P90 / P95 / P99 of the numeric valueAttribute), the optional groupByAttributes (up to 5 attribute keys) and the output metric's unit. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of the project this recording rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
