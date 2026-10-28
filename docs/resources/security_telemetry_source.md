---
page_title: "oneuptime_security_telemetry_source Resource - oneuptime"
subcategory: "Other"
description: |-
  Whether OneUptime looks for security events in one kind of telemetry the project already sends it - traces, network flows, the audit log or session replays - and what the last look found.
---

# oneuptime_security_telemetry_source (Resource)

Whether OneUptime looks for security events in one kind of telemetry the project already sends it - traces, network flows, the audit log or session replays - and what the last look found.

## Example Usage

```terraform
resource "oneuptime_security_telemetry_source" "example" {
  kind = "Example short text"
}
```

## Schema

### Required

- `kind` (String) Which telemetry this source reads: 'traces', 'network-flows', 'audit-log' or 'rum-sessions'. Fixed once created.

### Optional

- `is_enabled` (Boolean) Whether OneUptime looks for security events in this kind of the project's telemetry. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who turned this source on. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) The time up to which this telemetry has been read, as an ISO string: the end of the last window the scan finished.
- `id` (String) Unique identifier for the resource.
- `last_error` (String) Why the last scan failed, if it did. Cleared by the next scan that succeeds.
- `last_event_found_at` (String) When a security event was last found in this telemetry.
- `last_scan_result` (String) What the last scan read and found: the records read, the security events found, and how many of each signal. A JSON value: write it with `jsonencode()`.
- `last_scanned_at` (String) When this telemetry was last read.
- `project_id` (String) ID of the project whose telemetry is read. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security telemetry source by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_telemetry_source.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_telemetry_source.example <id>
```
