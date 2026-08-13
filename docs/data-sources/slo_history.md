---
page_title: "oneuptime_slo_history Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for SLO History
---

# oneuptime_slo_history (Data Source)

API endpoints for SLO History Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_slo_history" "by_name" {
  name = "example-slo_history"
}

data "oneuptime_slo_history" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `slo_id` (String) SLO ID. Computed.
- `metric_name` (String) Metric Name. Computed.
- `bucket_start` (String) Bucket Start. Computed.
- `value` (Number) Value. Computed.
- `version` (String) Version. Computed.
