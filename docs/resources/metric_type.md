---
page_title: "oneuptime_metric_type Resource - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  List of all the metrics ingested with OpenTelemetry
---

# oneuptime_metric_type (Resource)

List of all the metrics ingested with OpenTelemetry

## Example Usage

```terraform
resource "oneuptime_metric_type" "example" {
  name = "Example short text"
  description = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `services` (Set) List of services this metric is related to..
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Metric description..
- `unit` (String) Metric description..
- `is_monotonic` (Bool) Whether this metric is a monotonic counter (only ever increases), as reported by OpenTelemetry at ingest. Null when the instrument type does not carry monotonicity (e.g. gauges)...
- `aggregation_temporality` (String) OpenTelemetry aggregation temporality of this metric (Delta or Cumulative), as reported at ingest. Null when unknown...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_metric_type.example <id>
```
