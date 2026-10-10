---
page_title: "oneuptime_span Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  API endpoints for Span
---

# oneuptime_span (Data Source)

API endpoints for Span

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one span may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_span" "example" {
  name = "Example span"
}

# Or by id:
data "oneuptime_span" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `attributes` (String) Attributes.
- `container_entity_key` (String) Container Entity Key.
- `duration_unix_nano` (Number) Duration in Unix Nano.
- `end_time` (String) End Time.
- `end_time_unix_nano` (String) End Time.
- `has_exception` (Boolean) Has Exception.
- `host_entity_key` (String) Host Entity Key.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_llm_span` (Boolean) Is LLM Span.
- `is_root_span` (Boolean) Is Root Span.
- `k8s_cluster_entity_key` (String) Kubernetes Cluster Entity Key.
- `k8s_node_entity_key` (String) Kubernetes Node Entity Key.
- `k8s_pod_entity_key` (String) Kubernetes Pod Entity Key.
- `kind` (String) Kind.
- `llm_agent_name` (String) LLM Agent Name.
- `llm_call_kind` (String) LLM Call Kind.
- `llm_conversation_id` (String) LLM Conversation ID.
- `llm_cost` (Number) LLM Cost (USD).
- `llm_input_tokens` (Number) LLM Input Tokens.
- `llm_operation` (String) LLM Operation.
- `llm_output_tokens` (Number) LLM Output Tokens.
- `llm_request_model` (String) LLM Request Model.
- `llm_response_model` (String) LLM Response Model.
- `llm_system` (String) LLM System.
- `llm_team` (String) LLM Team.
- `llm_tool_name` (String) LLM Tool Name.
- `llm_total_tokens` (Number) LLM Total Tokens.
- `llm_user_email` (String) LLM User Email.
- `llm_user_id` (String) LLM User ID.
- `llm_user_message_preview` (String) LLM User Message Preview.
- `name` (String) Name.
- `parent_span_id` (String) Parent Span ID.
- `primary_entity_id` (String) Service ID.
- `primary_entity_type` (String) Service Type.
- `service_entity_key` (String) Service Entity Key.
- `session_id` (String) Session ID.
- `span_id` (String) Span ID.
- `start_time` (String) Start Time.
- `start_time_unix_nano` (String) Start Time in Unix Nano.
- `status_code` (Number) Status Code.
- `status_message` (String) Status Message.
- `trace_id` (String) Trace ID.
- `trace_state` (String) Trace State.

### Read-Only

- `attribute_keys` (Set of String) Attribute Keys.
- `entity_keys` (Set of String) Entity Keys.
- `events` (Set of String) Events.
- `links` (String) Links. A JSON value: write it with `jsonencode()`.
- `llm_issues` (Set of String) LLM Answer Issues.
- `project_id` (String) Project ID.
