---
page_title: "oneuptime_span Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  API endpoints for Span
---

# oneuptime_span (Data Source)

API endpoints for Span Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_span" "by_name" {
  name = "example-span"
}

data "oneuptime_span" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `primary_entity_id` (String) Service ID. Computed.
- `primary_entity_type` (String) Service Type. Computed.
- `start_time` (String) Start Time. Computed.
- `end_time` (String) End Time. Computed.
- `start_time_unix_nano` (String) Start Time in Unix Nano. Computed.
- `duration_unix_nano` (Number) Duration in Unix Nano. Computed.
- `end_time_unix_nano` (String) End Time. Computed.
- `trace_id` (String) Trace ID. Computed.
- `span_id` (String) Span ID. Computed.
- `session_id` (String) Session ID. Computed.
- `parent_span_id` (String) Parent Span ID. Computed.
- `trace_state` (String) Trace State. Computed.
- `attributes` (String) Attributes. Computed.
- `attribute_keys` (Set) Attribute Keys. Computed.
- `entity_keys` (Set) Entity Keys. Computed.
- `service_entity_key` (String) Service Entity Key. Computed.
- `host_entity_key` (String) Host Entity Key. Computed.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key. Computed.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key. Computed.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key. Computed.
- `container_entity_key` (String) Container Entity Key. Computed.
- `events` (Set) Events. Computed.
- `links` (String) Links. Computed.
- `status_code` (Number) Status Code. Computed.
- `status_message` (String) Status Message. Computed.
- `kind` (String) Kind. Computed.
- `has_exception` (Bool) Has Exception. Computed.
- `is_root_span` (Bool) Is Root Span. Computed.
- `is_llm_span` (Bool) Is LLM Span. Computed.
- `llm_system` (String) LLM System. Computed.
- `llm_operation` (String) LLM Operation. Computed.
- `llm_request_model` (String) LLM Request Model. Computed.
- `llm_response_model` (String) LLM Response Model. Computed.
- `llm_agent_name` (String) LLM Agent Name. Computed.
- `llm_tool_name` (String) LLM Tool Name. Computed.
- `llm_input_tokens` (Number) LLM Input Tokens. Computed.
- `llm_output_tokens` (Number) LLM Output Tokens. Computed.
- `llm_total_tokens` (Number) LLM Total Tokens. Computed.
- `llm_cost` (Number) LLM Cost (USD). Computed.
