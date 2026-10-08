---
page_title: "oneuptime_metric_type Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  List of all the metrics ingested with OpenTelemetry
---

# oneuptime_metric_type (Data Source)

List of all the metrics ingested with OpenTelemetry

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one metric type may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_metric_type" "example" {
  name = "Example metric type"
}

# Or by id:
data "oneuptime_metric_type" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `aggregation_temporality` (String) OpenTelemetry aggregation temporality of this metric (Delta or Cumulative), as reported at ingest. Null when unknown.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Metric description.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_monotonic` (Boolean) Whether this metric is a monotonic counter (only ever increases), as reported by OpenTelemetry at ingest. Null when the instrument type does not carry monotonicity (e.g. gauges).
- `name` (String) Any friendly name of this object.
- `unit` (String) Metric description.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `services` (Set of String) List of services this metric is related to. IDs of `oneuptime_service` resources.
- `updated_at` (String) Date and Time when the object was updated.
