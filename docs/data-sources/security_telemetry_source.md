---
page_title: "oneuptime_security_telemetry_source Data Source - oneuptime"
subcategory: "Security"
description: |-
  Whether OneUptime looks for security events in one kind of telemetry the project already sends it - traces, network flows, the audit log, session replays or the audit logs of its connected cloud accounts - and what the last look found.
---

# oneuptime_security_telemetry_source (Data Source)

Whether OneUptime looks for security events in one kind of telemetry the project already sends it - traces, network flows, the audit log, session replays or the audit logs of its connected cloud accounts - and what the last look found.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security telemetry source may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_telemetry_source" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_security_telemetry_source" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who turned this source on. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) The time up to which this telemetry has been read, as an ISO string: the end of the last window the scan finished.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether OneUptime looks for security events in this kind of the project's telemetry.
- `kind` (String) Which telemetry this source reads: 'traces', 'network-flows', 'audit-log', 'rum-sessions' or 'cloud-audit' (the audit logs of the cloud accounts connected under Cloud). Fixed once created.
- `last_error` (String) Why the last scan failed, if it did. Cleared by the next scan that succeeds.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_event_found_at` (String) When a security event was last found in this telemetry.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_scan_result` (String) What the last scan read and found: the records read, the security events found, and how many of each signal. A JSON value: read it with `jsondecode()`.
- `last_scanned_at` (String) When this telemetry was last read.
- `project_id` (String) ID of the project whose telemetry is read. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
