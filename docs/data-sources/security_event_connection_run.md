---
page_title: "oneuptime_security_event_connection_run Data Source - oneuptime"
subcategory: "Security"
description: |-
  History of connection tests, previews, scheduled polls and historical imports for security event connections. Credentials are never included.
---

# oneuptime_security_event_connection_run (Data Source)

History of connection tests, previews, scheduled polls and historical imports for security event connections. Credentials are never included.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security event connection run may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_event_connection_run" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_security_event_connection_run" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `error` (String) The run failure with credentials redacted.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `requested_by_user_id` (String) ID of the user who requested this run. The ID of a `oneuptime_user` (see the data source).
- `security_event_connection_id` (String) ID of the connection for this run. The ID of a `oneuptime_security_event_connection`.
- `status` (String) Status of this connection run.
- `type` (String) Operation of this connection run.

### Read-Only

- `completed_at` (String) When this run completed.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project for this run. The ID of a `oneuptime_project`.
- `request` (String) Validated operation and selected time range. Contains no credentials. A JSON value: read it with `jsondecode()`.
- `result` (String) Counts, requested time range, checks and a bounded preview of records. A JSON value: read it with `jsondecode()`.
- `started_at` (String) When this run started.
- `updated_at` (String) Date and Time when the object was updated.
