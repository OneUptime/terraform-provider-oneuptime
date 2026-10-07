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
- `output_metric_name` (String) Name of the new metric this rule writes (e.g. sdwan.gateway.latency.ms). Leave it out and it is made from the rule's name - SD-WAN gateway latency becomes sd_wan_gateway_latency, with _2, _3 and so on added when another recording rule of the project already writes it. Keep it unique per project...
- `definition` (String) Which logs count (filter: telemetryServiceIds, severityTexts, body, attributeFilters), how they are aggregated (aggregationType: Count, or Sum / Avg / Min / Max / P50 / P75 / P90 / P95 / P99 of the numeric valueAttribute), the optional groupByAttributes (up to 5 attribute keys) and the output metric's unit...
- `is_enabled` (Bool) Whether this rule is evaluated by the log recording rule cron, every minute...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `computed_until` (String) A date time object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_log_recording_rule.example <id>
```
