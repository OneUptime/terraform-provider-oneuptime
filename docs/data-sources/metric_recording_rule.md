---
page_title: "oneuptime_metric_recording_rule Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Derived metrics computed on a schedule from an expression over other metrics. Results are written back into the metric store as a new series.
---

# oneuptime_metric_recording_rule (Data Source)

Derived metrics computed on a schedule from an expression over other metrics. Results are written back into the metric store as a new series. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_metric_recording_rule" "by_name" {
  name = "example-metric_recording_rule"
}

data "oneuptime_metric_recording_rule" "by_id" {
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
- `description` (String) What this recording rule computes and why... Computed.
- `output_metric_name` (String) Name of the new metric this rule writes (e.g. http.error_rate). Must be unique per project... Computed.
- `definition` (String) Sources (aliased input metrics), arithmetic expression, and optional group-by attribute... Computed.
- `is_enabled` (Bool) Whether this rule is evaluated by the recording rule cron... Computed.
- `sort_order` (Number) Evaluation order when multiple rules exist... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
