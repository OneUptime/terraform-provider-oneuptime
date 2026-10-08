---
page_title: "oneuptime_slo_history Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for SLO History
---

# oneuptime_slo_history (Data Source)

API endpoints for SLO History

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one slo history may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_slo_history" "example" {
  slo_id = "example-slo-id"
}

# Or by id:
data "oneuptime_slo_history" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `metric_name` (String) Metric Name.
- `slo_id` (String) SLO ID.
- `value` (Number) Value.
- `version` (String) Version.

### Read-Only

- `bucket_start` (String) Bucket Start.
- `project_id` (String) Project ID.
