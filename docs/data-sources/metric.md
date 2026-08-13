---
page_title: "oneuptime_metric Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  API endpoints for Metric
---

# oneuptime_metric (Data Source)

API endpoints for Metric Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_metric" "by_name" {
  name = "example-metric"
}

data "oneuptime_metric" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `primary_entity_id` (String) Service ID. Computed.
- `primary_entity_type` (String) Service Type. Computed.
- `aggregation_temporality` (String) Aggregation Temporality. Computed.
- `metric_point_type` (String) Metric Point Type. Computed.
- `time` (String) Time. Computed.
- `start_time` (String) Start Time. Computed.
- `time_unix_nano` (String) Time (in Unix Nano). Computed.
- `start_time_unix_nano` (String) Start Time (in Unix Nano). Computed.
- `attributes` (String) Attributes. Computed.
- `attribute_keys` (Set) Attribute Keys. Computed.
- `entity_keys` (Set) Entity Keys. Computed.
- `service_entity_key` (String) Service Entity Key. Computed.
- `host_entity_key` (String) Host Entity Key. Computed.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key. Computed.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key. Computed.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key. Computed.
- `container_entity_key` (String) Container Entity Key. Computed.
- `is_monotonic` (Bool) Is Monotonic. Computed.
- `count` (String) Count. Computed.
- `sum` (Number) Sum. Computed.
- `value` (Number) Value. Computed.
- `min` (Number) Min. Computed.
- `max` (Number) Max. Computed.
- `bucket_counts` (String) Bucket Counts. Computed.
- `explicit_bounds` (String) Explicit Bounds. Computed.
- `scale` (Number) Scale. Computed.
- `zero_count` (String) Zero Count. Computed.
- `positive_offset` (Number) Positive Bucket Offset. Computed.
- `positive_bucket_counts` (String) Positive Bucket Counts. Computed.
- `negative_offset` (Number) Negative Bucket Offset. Computed.
- `negative_bucket_counts` (String) Negative Bucket Counts. Computed.
- `summary_quantiles` (String) Summary Quantiles. Computed.
- `summary_values` (String) Summary Values. Computed.
- `trace_id` (String) Trace ID. Computed.
- `span_id` (String) Span ID. Computed.
