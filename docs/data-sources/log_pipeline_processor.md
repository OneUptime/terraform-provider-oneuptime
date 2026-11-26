---
page_title: "oneuptime_log_pipeline_processor Data Source - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Individual processors within a log pipeline that transform log data during ingestion.
---

# oneuptime_log_pipeline_processor (Data Source)

Individual processors within a log pipeline that transform log data during ingestion.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one log pipeline processor may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_log_pipeline_processor" "example" {
  name = "Example log pipeline processor"
}

# Or by id:
data "oneuptime_log_pipeline_processor" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this log pipeline processor. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this processor is active.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `log_pipeline_id` (String) ID of the log pipeline this processor belongs to. The ID of a `oneuptime_log_pipeline`.
- `name` (String) Friendly name for this processor.
- `processor_type` (String) The type of processor: GrokParser, KeyValueParser, AttributeRemapper, SeverityRemapper, or CategoryProcessor.
- `sort_order` (Number) Where this processor runs within its pipeline, lowest number first. A new processor is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them.

### Read-Only

- `configuration` (String) Processor-specific configuration as JSON (e.g., grok pattern, source/target fields, mapping rules). A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this log pipeline processor belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
