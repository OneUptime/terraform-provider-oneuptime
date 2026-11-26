---
page_title: "oneuptime_trace_recording_rule Data Source - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Derived metrics computed on a schedule from aggregations over spans. Results are written into the metric store as a new series.
---

# oneuptime_trace_recording_rule (Data Source)

Derived metrics computed on a schedule from aggregations over spans. Results are written into the metric store as a new series.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one trace recording rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_trace_recording_rule" "example" {
  name = "Example trace recording rule"
}

# Or by id:
data "oneuptime_trace_recording_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this rule. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What this recording rule computes and why.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is evaluated by the recording rule cron.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Friendly name for this rule.
- `output_metric_name` (String) Name of the new metric this rule writes (e.g. http.error_rate). Leave it out and it is made from the rule's name - HTTP error rate becomes http_error_rate, with _2, _3 and so on added when another recording rule of the project already writes it. Keep it unique per project.
- `sort_order` (Number) Not read when rules are evaluated: every enabled rule is evaluated each minute, on its own, whatever this holds. The dashboard lists rules by name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `definition` (String) Sources (aliased span aggregations), arithmetic expression, and optional group-by attribute. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this recording rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
