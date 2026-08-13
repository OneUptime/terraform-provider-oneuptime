---
page_title: "oneuptime_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  API endpoints for Log
---

# oneuptime_log (Data Source)

API endpoints for Log Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_log" "by_name" {
  name = "example-log"
}

data "oneuptime_log" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `primary_entity_id` (String) Service ID. Computed.
- `primary_entity_type` (String) Service Type. Computed.
- `time` (String) Time. Computed.
- `time_unix_nano` (String) Time (in Unix Nano). Computed.
- `severity_text` (String) Severity Text. Computed.
- `severity_number` (Number) Severity Number. Computed.
- `attributes` (String) Attributes. Computed.
- `attribute_keys` (Set) Attribute Keys. Computed.
- `entity_keys` (Set) Entity Keys. Computed.
- `service_entity_key` (String) Service Entity Key. Computed.
- `host_entity_key` (String) Host Entity Key. Computed.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key. Computed.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key. Computed.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key. Computed.
- `container_entity_key` (String) Container Entity Key. Computed.
- `trace_id` (String) Trace ID. Computed.
- `span_id` (String) Span ID. Computed.
- `session_id` (String) Session ID. Computed.
- `body` (String) Log Body. Computed.
- `observed_time_unix_nano` (String) Observed Time (in Unix Nano). Computed.
- `dropped_attributes_count` (Number) Dropped Attributes Count. Computed.
- `flags` (Number) Flags. Computed.
