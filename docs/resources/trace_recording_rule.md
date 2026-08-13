---
page_title: "oneuptime_trace_recording_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Derived metrics computed on a schedule from aggregations over spans. Results are written into the metric store as a new series.
---

# oneuptime_trace_recording_rule (Resource)

Derived metrics computed on a schedule from aggregations over spans. Results are written into the metric store as a new series.

## Example Usage

```terraform
resource "oneuptime_trace_recording_rule" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  output_metric_name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `output_metric_name` (String) Name of the new metric this rule writes (e.g. http.error_rate). Must be unique per project...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) What this recording rule computes and why...
- `definition` (String) Sources (aliased span aggregations), arithmetic expression, and optional group-by attribute...
- `is_enabled` (Bool) Whether this rule is evaluated by the recording rule cron...
- `sort_order` (Number) Evaluation order when multiple rules exist...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_trace_recording_rule.example <id>
```
