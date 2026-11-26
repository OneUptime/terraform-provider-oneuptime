---
page_title: "oneuptime_data_asset_usage Data Source - oneuptime"
subcategory: "Other"
description: |-
  A service seen reading or writing a data asset in its database spans, with the operations it ran and when it last ran them.
---

# oneuptime_data_asset_usage (Data Source)

A service seen reading or writing a data asset in its database spans, with the operations it ran and when it last ran them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data asset usage may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_asset_usage" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_data_asset_usage" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `access_kind` (String) Whether the service writes the table's rows or reads them.
- `call_count` (Number) How many calls the service made to the table in the last hour.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_asset_id` (String) ID of the table or view the service read or wrote. The ID of a `oneuptime_data_asset`.
- `data_store_id` (String) ID of the data store the table is in.
- `error_count` (Number) How many of those calls failed.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ai_app` (Boolean) Whether the service also calls a model: reading this table, it is an AI app whose answers depend on it.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `service_id` (String) ID of the service whose spans read or wrote the table. The ID of a `oneuptime_service`.
- `service_name` (String) The service's name when it was last seen using the table.
- `usage_key` (String) The service and the way it uses the table, one row per pair: what the next discovery updates.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `endpoints` (String) The service's busiest endpoints that made these calls, from the server spans of the same traces. A JSON value: read it with `jsondecode()`.
- `first_seen_at` (String) When the service was first seen using the table this way.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_seen_at` (String) When the service was last seen using the table this way.
- `operations` (String) The statements the service ran on the table, such as INSERT or SELECT. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
