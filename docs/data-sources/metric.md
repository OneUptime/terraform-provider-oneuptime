---
page_title: "oneuptime_metric Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  API endpoints for Metric
---

# oneuptime_metric (Data Source)

API endpoints for Metric

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one metric may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_metric" "example" {
  name = "Example metric"
}

# Or by id:
data "oneuptime_metric" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `aggregation_temporality` (String) Aggregation Temporality.
- `attributes` (String) Attributes.
- `bucket_counts` (String) Bucket Counts.
- `container_entity_key` (String) Container Entity Key.
- `count_value` (String) Count.
- `explicit_bounds` (String) Explicit Bounds.
- `host_entity_key` (String) Host Entity Key.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_monotonic` (Boolean) Is Monotonic.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key.
- `max` (Number) Max.
- `metric_point_type` (String) Metric Point Type.
- `min` (Number) Min.
- `name` (String) Name.
- `negative_bucket_counts` (String) Negative Bucket Counts.
- `negative_offset` (Number) Negative Bucket Offset.
- `positive_bucket_counts` (String) Positive Bucket Counts.
- `positive_offset` (Number) Positive Bucket Offset.
- `primary_entity_id` (String) Service ID.
- `primary_entity_type` (String) Service Type.
- `scale` (Number) Scale.
- `service_entity_key` (String) Service Entity Key.
- `span_id` (String) Span ID.
- `start_time` (String) Start Time.
- `start_time_unix_nano` (String) Start Time (in Unix Nano).
- `sum` (Number) Sum.
- `summary_quantiles` (String) Summary Quantiles.
- `summary_values` (String) Summary Values.
- `time` (String) Time.
- `time_unix_nano` (String) Time (in Unix Nano).
- `trace_id` (String) Trace ID.
- `value` (Number) Value.
- `zero_count` (String) Zero Count.

### Read-Only

- `attribute_keys` (Set of String) Attribute Keys.
- `entity_keys` (Set of String) Entity Keys.
- `project_id` (String) Project ID.
