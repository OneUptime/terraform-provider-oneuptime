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
  name        = "Example metric recording rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this rule.

### Optional

- `definition` (String) Sources (aliased input metrics), arithmetic expression, and optional group-by attribute. A JSON value: write it with `jsonencode()`.
- `description` (String) What this recording rule computes and why.
- `is_enabled` (Boolean) Whether this rule is evaluated by the recording rule cron. Defaults to `true`.
- `output_metric_name` (String) Name of the new metric this rule writes (e.g. http.error_rate). Leave it out and it is made from the rule's name - HTTP error rate becomes http_error_rate, with _2, _3 and so on added when another recording rule of the project already writes it. Keep it unique per project.
- `sort_order` (Number) Not read when rules are evaluated: every enabled rule is evaluated each minute, on its own, whatever this holds. The dashboard lists rules by name. Defaults to `0`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this rule. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this recording rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing metric recording rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_metric_recording_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_metric_recording_rule.example <id>
```
