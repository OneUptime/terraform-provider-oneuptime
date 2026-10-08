---
page_title: "oneuptime_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  API endpoints for Log
---

# oneuptime_log (Data Source)

API endpoints for Log

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_log" "example" {
  primary_entity_id = "example-primary-entity-id"
}

# Or by id:
data "oneuptime_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `attributes` (String) Attributes.
- `body` (String) Log Body.
- `container_entity_key` (String) Container Entity Key.
- `dropped_attributes_count` (Number) Dropped Attributes Count.
- `flags` (Number) Flags.
- `host_entity_key` (String) Host Entity Key.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key.
- `observed_time_unix_nano` (String) Observed Time (in Unix Nano).
- `primary_entity_id` (String) Service ID.
- `primary_entity_type` (String) Service Type.
- `service_entity_key` (String) Service Entity Key.
- `session_id` (String) Session ID.
- `severity_number` (Number) Severity Number.
- `severity_text` (String) Severity Text.
- `span_id` (String) Span ID.
- `time` (String) Time.
- `time_unix_nano` (String) Time (in Unix Nano).
- `trace_id` (String) Trace ID.

### Read-Only

- `attribute_keys` (Set of String) Attribute Keys.
- `entity_keys` (Set of String) Entity Keys.
- `project_id` (String) Project ID.
