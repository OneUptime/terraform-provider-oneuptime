---
page_title: "oneuptime_log_recording_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Derived metrics computed every minute from logs: a count of matching logs, or the sum, average, min, max or a percentile of a numeric log attribute, optionally split by log attributes. Results are written into the metric store as a new series.
---

# oneuptime_log_recording_rule (Resource)

Derived metrics computed every minute from logs: a count of matching logs, or the sum, average, min, max or a percentile of a numeric log attribute, optionally split by log attributes. Results are written into the metric store as a new series.

## Example Usage

```terraform
resource "oneuptime_log_recording_rule" "example" {
  name        = "Example log recording rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this rule.

### Optional

- `definition` (String) Which logs count (filter: telemetryServiceIds, severityTexts, body, attributeFilters), how they are aggregated (aggregationType: Count, or Sum / Avg / Min / Max / P50 / P75 / P90 / P95 / P99 of the numeric valueAttribute), the optional groupByAttributes (up to 5 attribute keys) and the output metric's unit. A JSON value: write it with `jsonencode()`.
- `description` (String) What this recording rule computes and why.
- `is_enabled` (Boolean) Whether this rule is evaluated by the log recording rule cron, every minute. Defaults to `true`.
- `output_metric_name` (String) Name of the new metric this rule writes (e.g. sdwan.gateway.latency.ms). Leave it out and it is made from the rule's name - SD-WAN gateway latency becomes sd_wan_gateway_latency, with _2, _3 and so on added when another recording rule of the project already writes it. Keep it unique per project.

### Read-Only

- `computed_until` (String) The end of the last minute this rule has written into its metric. Null means it has not run yet.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this rule. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this recording rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing log recording rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_log_recording_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_log_recording_rule.example <id>
```
