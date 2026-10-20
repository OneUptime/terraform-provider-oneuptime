---
page_title: "oneuptime_security_log_source Data Source - oneuptime"
subcategory: "Other"
description: |-
  Whether OneUptime looks for security events in the logs the project already sends it, and what the last look found.
---

# oneuptime_security_log_source (Data Source)

Whether OneUptime looks for security events in the logs the project already sends it, and what the last look found.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security log source may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_log_source" "example" {
  cursor = "example-cursor"
}

# Or by id:
data "oneuptime_security_log_source" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who turned this source on. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) The time up to which logs have been read, as an ISO string: the end of the last window the scan finished.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether OneUptime looks for security events in the project's logs.
- `last_error` (String) Why the last scan failed, if it did. Cleared by the next scan that succeeds.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_event_found_at` (String) When a security event was last found in the logs.
- `last_scan_result` (String) What the last scan read and found: the log lines read, the security events found, and how many of each kind. A JSON value: write it with `jsonencode()`.
- `last_scanned_at` (String) When the logs were last read.
- `project_id` (String) ID of the project whose logs are read. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
