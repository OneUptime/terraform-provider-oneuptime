---
page_title: "oneuptime_google_sec_ops_connection_run Data Source - oneuptime"
subcategory: "Other"
description: |-
  History of connection tests, previews, scheduled polls and historical imports. Credentials are never included.
---

# oneuptime_google_sec_ops_connection_run (Data Source)

History of connection tests, previews, scheduled polls and historical imports. Credentials are never included. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_google_sec_ops_connection_run" "by_name" {
  name = "example-google_sec_ops_connection_run"
}

data "oneuptime_google_sec_ops_connection_run" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `google_sec_ops_connection_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `requested_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `type` (String) Operation of this connection run... Computed.
- `status` (String) Status of this connection run... Computed.
- `started_at` (String) A date time object.. Computed.
- `completed_at` (String) A date time object.. Computed.
- `request` (String) Validated operation and selected time range. Contains no credentials... Computed.
- `result` (String) Counts, requested time range, checks and a bounded preview of detections... Computed.
- `error` (String) The run failure with credentials redacted... Computed.
