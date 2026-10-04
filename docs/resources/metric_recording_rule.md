---
page_title: "oneuptime_metric_recording_rule Resource - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Derived metrics computed on a schedule from an expression over other metrics. Results are written back into the metric store as a new series.
---

# oneuptime_metric_recording_rule (Resource)

Derived metrics computed on a schedule from an expression over other metrics. Results are written back into the metric store as a new series.

## Example Usage

```terraform
resource "oneuptime_metric_recording_rule" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) What this recording rule computes and why...
- `output_metric_name` (String) Name of the new metric this rule writes (e.g. http.error_rate). Leave it out and it is made from the rule's name - HTTP error rate becomes http_error_rate, with _2, _3 and so on added when another recording rule of the project already writes it. Keep it unique per project...
- `definition` (String) Sources (aliased input metrics), arithmetic expression, and optional group-by attribute...
- `is_enabled` (Bool) Whether this rule is evaluated by the recording rule cron...
- `sort_order` (Number) Not read when rules are evaluated: every enabled rule is evaluated each minute, on its own, whatever this holds. The dashboard lists rules by name...

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
terraform import oneuptime_metric_recording_rule.example <id>
```
