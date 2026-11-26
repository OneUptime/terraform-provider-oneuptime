---
page_title: "oneuptime_exception_saved_view Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Save and reuse views of the exceptions list, including the current search, filters, time range, status, sort and page size.
---

# oneuptime_exception_saved_view (Data Source)

Save and reuse views of the exceptions list, including the current search, filters, time range, status, sort and page size.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one exception saved view may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_exception_saved_view" "example" {
  name = "Example exception saved view"
}

# Or by id:
data "oneuptime_exception_saved_view" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this saved exceptions view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_default` (Boolean) Whether this saved exceptions view should be applied by default.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Friendly name for this saved exceptions view.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this saved exceptions view belongs to. The ID of a `oneuptime_project`.
- `query` (String) Serialized exceptions list view state (search, filters, time range, page size, status, lens, assignee filter and sort) for this saved view. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
