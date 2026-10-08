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
  name        = "Example metric type"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `aggregation_temporality` (String) OpenTelemetry aggregation temporality of this metric (Delta or Cumulative), as reported at ingest. Null when unknown.
- `description` (String) Metric description.
- `is_monotonic` (Boolean) Whether this metric is a monotonic counter (only ever increases), as reported by OpenTelemetry at ingest. Null when the instrument type does not carry monotonicity (e.g. gauges).
- `services` (Set of String) List of services this metric is related to. IDs of `oneuptime_service` resources.
- `unit` (String) Metric description.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing metric type by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_metric_type.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_metric_type.example <id>
```
