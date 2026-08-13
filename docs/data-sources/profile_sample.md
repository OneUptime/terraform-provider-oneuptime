---
page_title: "oneuptime_profile_sample Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for ProfileSample
---

# oneuptime_profile_sample (Data Source)

API endpoints for ProfileSample Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_profile_sample" "by_name" {
  name = "example-profile_sample"
}

data "oneuptime_profile_sample" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `primary_entity_id` (String) Service ID. Computed.
- `primary_entity_type` (String) Service Type. Computed.
- `profile_id` (String) Profile ID. Computed.
- `trace_id` (String) Trace ID. Computed.
- `span_id` (String) Span ID. Computed.
- `time` (String) Time. Computed.
- `time_unix_nano` (String) Time (in Unix Nano). Computed.
- `stacktrace` (Set) Stacktrace. Computed.
- `stacktrace_hash` (String) Stacktrace Hash. Computed.
- `frame_types` (Set) Frame Types. Computed.
- `value` (String) Value. Computed.
- `profile_type` (String) Profile Type. Computed.
- `labels` (String) Labels. Computed.
- `entity_keys` (Set) Entity Keys. Computed.
- `service_entity_key` (String) Service Entity Key. Computed.
- `host_entity_key` (String) Host Entity Key. Computed.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key. Computed.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key. Computed.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key. Computed.
- `container_entity_key` (String) Container Entity Key. Computed.
