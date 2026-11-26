---
page_title: "oneuptime_log_drop_filter Data Source - oneuptime"
subcategory: "Logs, Traces & Metrics"
description: |-
  Configure rules to drop or sample logs before storage to reduce volume and cost.
---

# oneuptime_log_drop_filter (Data Source)

Configure rules to drop or sample logs before storage to reduce volume and cost.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one log drop filter may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_log_drop_filter" "example" {
  name = "Example log drop filter"
}

# Or by id:
data "oneuptime_log_drop_filter" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `action` (String) What to do with matching logs: 'drop' to discard entirely, 'sample' to keep a percentage.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this log drop filter. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of what this drop filter does.
- `dropped_count` (Number) Total number of logs this filter has discarded since it was created.
- `filter_query` (String) Filter expression that identifies which logs to drop or sample.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this drop filter is active.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Friendly name for this drop filter.
- `sample_percentage` (Number) When action is 'sample', the percentage of matching logs to keep (1-99).
- `sort_order` (Number) Where this filter is evaluated among the project's log drop filters, lowest number first. A new filter is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_dropped_at` (String) When this filter most recently discarded a log. Null means it has never matched anything.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this log drop filter belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
