---
page_title: "oneuptime_profile_sample Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for ProfileSample
---

# oneuptime_profile_sample (Data Source)

API endpoints for ProfileSample

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one profile sample may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_profile_sample" "example" {
  primary_entity_id = "example-primary-entity-id"
}

# Or by id:
data "oneuptime_profile_sample" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `container_entity_key` (String) Container Entity Key.
- `host_entity_key` (String) Host Entity Key.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key.
- `labels` (String) Labels.
- `primary_entity_id` (String) Service ID.
- `primary_entity_type` (String) Service Type.
- `profile_id` (String) Profile ID.
- `profile_type` (String) Profile Type.
- `service_entity_key` (String) Service Entity Key.
- `span_id` (String) Span ID.
- `stacktrace_hash` (String) Stacktrace Hash.
- `time` (String) Time.
- `time_unix_nano` (String) Time (in Unix Nano).
- `trace_id` (String) Trace ID.
- `value` (String) Value.

### Read-Only

- `entity_keys` (Set of String) Entity Keys.
- `frame_types` (Set of String) Frame Types.
- `project_id` (String) Project ID.
- `stacktrace` (Set of String) Stacktrace.
