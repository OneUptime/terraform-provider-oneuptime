---
page_title: "oneuptime_exception_instance Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  API endpoints for Exception Instance
---

# oneuptime_exception_instance (Data Source)

API endpoints for Exception Instance Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_exception_instance" "by_name" {
  name = "example-exception_instance"
}

data "oneuptime_exception_instance" "by_id" {
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
- `exception_type` (String) Exception Type. Computed.
- `stack_trace` (String) Stack Trace. Computed.
- `message` (String) Exception Message. Computed.
- `span_status_code` (Number) Span Status Code. Computed.
- `escaped` (Bool) Exception Escaped. Computed.
- `trace_id` (String) Trace ID. Computed.
- `span_id` (String) Span ID. Computed.
- `session_id` (String) Session ID. Computed.
- `fingerprint` (String) Fingerprint. Computed.
- `span_name` (String) Span Name. Computed.
- `release` (String) Release. Computed.
- `environment` (String) Environment. Computed.
- `parsed_frames` (String) Parsed Stack Frames. Computed.
- `attributes` (String) Attributes. Computed.
- `attribute_keys` (Set) Attribute Keys. Computed.
- `entity_keys` (Set) Entity Keys. Computed.
- `service_entity_key` (String) Service Entity Key. Computed.
- `host_entity_key` (String) Host Entity Key. Computed.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key. Computed.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key. Computed.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key. Computed.
- `container_entity_key` (String) Container Entity Key. Computed.
