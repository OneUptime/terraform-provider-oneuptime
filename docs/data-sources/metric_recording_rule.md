---
page_title: "oneuptime_metric_recording_rule Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Derived metrics computed on a schedule from an expression over other metrics. Results are written back into the metric store as a new series.
---

# oneuptime_metric_recording_rule (Data Source)

Derived metrics computed on a schedule from an expression over other metrics. Results are written back into the metric store as a new series.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one metric recording rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_metric_recording_rule" "example" {
  name = "Example metric recording rule"
}

# Or by id:
data "oneuptime_metric_recording_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this rule. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What this recording rule computes and why.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is evaluated by the recording rule cron.
- `name` (String) Friendly name for this rule.
- `output_metric_name` (String) Name of the new metric this rule writes (e.g. http.error_rate). Leave it out and it is made from the rule's name - HTTP error rate becomes http_error_rate, with _2, _3 and so on added when another recording rule of the project already writes it. Keep it unique per project.
- `sort_order` (Number) Not read when rules are evaluated: every enabled rule is evaluated each minute, on its own, whatever this holds. The dashboard lists rules by name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `definition` (String) Sources (aliased input metrics), arithmetic expression, and optional group-by attribute. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of the project this recording rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
