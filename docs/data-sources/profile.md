---
page_title: "oneuptime_profile Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Profile
---

# oneuptime_profile (Data Source)

API endpoints for Profile

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one profile may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_profile" "example" {
  primary_entity_id = "example-primary-entity-id"
}

# Or by id:
data "oneuptime_profile" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `attributes` (String) Attributes.
- `container_entity_key` (String) Container Entity Key.
- `duration_nano` (String) Duration in Nanoseconds.
- `end_time` (String) End Time.
- `end_time_unix_nano` (String) End Time in Unix Nano.
- `host_entity_key` (String) Host Entity Key.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key.
- `original_payload_format` (String) Original Payload Format.
- `period` (String) Period.
- `period_type` (String) Period Type.
- `primary_entity_id` (String) Service ID.
- `primary_entity_type` (String) Service Type.
- `profile_id` (String) Profile ID.
- `profile_type` (String) Profile Type.
- `sample_count` (Number) Sample Count.
- `service_entity_key` (String) Service Entity Key.
- `span_id` (String) Span ID.
- `start_time` (String) Start Time.
- `start_time_unix_nano` (String) Start Time in Unix Nano.
- `trace_id` (String) Trace ID.
- `unit` (String) Unit.

### Read-Only

- `attribute_keys` (Set of String) Attribute Keys.
- `entity_keys` (Set of String) Entity Keys.
- `project_id` (String) Project ID.
