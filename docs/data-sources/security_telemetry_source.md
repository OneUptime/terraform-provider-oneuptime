---
page_title: "oneuptime_security_telemetry_source Data Source - oneuptime"
subcategory: "Other"
description: |-
  Whether OneUptime looks for security events in one kind of telemetry the project already sends it - traces, network flows, the audit log or session replays - and what the last look found.
---

# oneuptime_security_telemetry_source (Data Source)

Whether OneUptime looks for security events in one kind of telemetry the project already sends it - traces, network flows, the audit log or session replays - and what the last look found.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security telemetry source may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_telemetry_source" "example" {
  kind = "example-kind"
}

# Or by id:
data "oneuptime_security_telemetry_source" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who turned this source on. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) The time up to which this telemetry has been read, as an ISO string: the end of the last window the scan finished.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether OneUptime looks for security events in this kind of the project's telemetry.
- `kind` (String) Which telemetry this source reads: 'traces', 'network-flows', 'audit-log' or 'rum-sessions'. Fixed once created.
- `last_error` (String) Why the last scan failed, if it did. Cleared by the next scan that succeeds.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_event_found_at` (String) When a security event was last found in this telemetry.
- `last_scan_result` (String) What the last scan read and found: the records read, the security events found, and how many of each signal. A JSON value: write it with `jsonencode()`.
- `last_scanned_at` (String) When this telemetry was last read.
- `project_id` (String) ID of the project whose telemetry is read. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
