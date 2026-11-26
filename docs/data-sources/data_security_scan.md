---
page_title: "oneuptime_data_security_scan Data Source - oneuptime"
subcategory: "Data Security"
description: |-
  Every scan Data Security ran: what it read, how much, what it found, how long it took and what went wrong.
---

# oneuptime_data_security_scan (Data Source)

Every scan Data Security ran: what it read, how much, what it found, how long it took and what went wrong.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data security scan may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_security_scan" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_data_security_scan" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `bytes_read` (Number) How many bytes the scan read.
- `cloud_connection_id` (String) The connection read, for a configuration read. The ID of a `oneuptime_cloud_connection`.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `data_store_id` (String) The store scanned, for a content or telemetry scan. The ID of a `oneuptime_data_store`.
- `duration_in_ms` (Number) How long the scan took.
- `error_message` (String) Why the scan failed, if it did.
- `findings_found` (Number) How many findings the scan made or updated.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `kind` (String) What the scan read: Configuration, Content or Telemetry.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `objects_read` (Number) How many objects or files the scan sampled.
- `probe_id` (String) The probe that ran the scan, when one did. The ID of a `oneuptime_probe`.
- `rows_read` (Number) How many rows or records the scan sampled.
- `status` (String) Running, Succeeded, Partly succeeded, Failed or Skipped.
- `stores_read` (Number) How many stores a configuration read covered.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `finished_at` (String) When the scan finished.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this scan belongs to. The ID of a `oneuptime_project`.
- `read_log` (String) The objects, tables or services read, by name. Never their contents. A JSON value: read it with `jsondecode()`.
- `started_at` (String) When the scan started.
- `updated_at` (String) Date and Time when the object was updated.
- `warnings` (String) What the scan could not read, and why. A JSON value: read it with `jsondecode()`.
