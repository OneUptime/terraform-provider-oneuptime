---
page_title: "oneuptime_profile Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Profile
---

# oneuptime_profile (Data Source)

API endpoints for Profile Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_profile" "by_name" {
  name = "example-profile"
}

data "oneuptime_profile" "by_id" {
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
- `start_time` (String) Start Time. Computed.
- `end_time` (String) End Time. Computed.
- `start_time_unix_nano` (String) Start Time in Unix Nano. Computed.
- `end_time_unix_nano` (String) End Time in Unix Nano. Computed.
- `duration_nano` (String) Duration in Nanoseconds. Computed.
- `profile_type` (String) Profile Type. Computed.
- `unit` (String) Unit. Computed.
- `period_type` (String) Period Type. Computed.
- `period` (String) Period. Computed.
- `attributes` (String) Attributes. Computed.
- `attribute_keys` (Set) Attribute Keys. Computed.
- `entity_keys` (Set) Entity Keys. Computed.
- `service_entity_key` (String) Service Entity Key. Computed.
- `host_entity_key` (String) Host Entity Key. Computed.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key. Computed.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key. Computed.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key. Computed.
- `container_entity_key` (String) Container Entity Key. Computed.
- `sample_count` (Number) Sample Count. Computed.
- `original_payload_format` (String) Original Payload Format. Computed.
