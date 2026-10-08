---
page_title: "oneuptime_exception_instance Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  API endpoints for Exception Instance
---

# oneuptime_exception_instance (Data Source)

API endpoints for Exception Instance

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one exception instance may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_exception_instance" "example" {
  primary_entity_id = "example-primary-entity-id"
}

# Or by id:
data "oneuptime_exception_instance" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `attributes` (String) Attributes.
- `container_entity_key` (String) Container Entity Key.
- `environment` (String) Environment.
- `escaped` (Boolean) Exception Escaped.
- `exception_type` (String) Exception Type.
- `fingerprint` (String) Fingerprint.
- `host_entity_key` (String) Host Entity Key.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key.
- `message` (String) Exception Message.
- `parsed_frames` (String) Parsed Stack Frames.
- `primary_entity_id` (String) Service ID.
- `primary_entity_type` (String) Service Type.
- `release` (String) Release.
- `service_entity_key` (String) Service Entity Key.
- `session_id` (String) Session ID.
- `span_id` (String) Span ID.
- `span_name` (String) Span Name.
- `span_status_code` (Number) Span Status Code.
- `stack_trace` (String) Stack Trace.
- `time` (String) Time.
- `time_unix_nano` (String) Time (in Unix Nano).
- `trace_id` (String) Trace ID.

### Read-Only

- `attribute_keys` (Set of String) Attribute Keys.
- `entity_keys` (Set of String) Entity Keys.
- `project_id` (String) Project ID.
