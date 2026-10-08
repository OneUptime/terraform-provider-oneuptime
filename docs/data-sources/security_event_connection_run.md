---
page_title: "oneuptime_security_event_connection_run Data Source - oneuptime"
subcategory: "Other"
description: |-
  History of connection tests, previews, scheduled polls and historical imports for security event connections. Credentials are never included.
---

# oneuptime_security_event_connection_run (Data Source)

History of connection tests, previews, scheduled polls and historical imports for security event connections. Credentials are never included.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security event connection run may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_event_connection_run" "example" {
  security_event_connection_id = oneuptime_security_event_connection.example.id
}

# Or by id:
data "oneuptime_security_event_connection_run" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `error` (String) The run failure with credentials redacted.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `requested_by_user_id` (String) ID of the user who requested this run. The ID of a `oneuptime_user` (see the data source).
- `security_event_connection_id` (String) ID of the connection for this run. The ID of a `oneuptime_security_event_connection`.
- `status` (String) Status of this connection run.
- `type` (String) Operation of this connection run.

### Read-Only

- `completed_at` (String) When this run completed.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project for this run. The ID of a `oneuptime_project`.
- `request` (String) Validated operation and selected time range. Contains no credentials. A JSON value: write it with `jsonencode()`.
- `result` (String) Counts, requested time range, checks and a bounded preview of records. A JSON value: write it with `jsonencode()`.
- `started_at` (String) When this run started.
- `updated_at` (String) Date and Time when the object was updated.
