---
page_title: "oneuptime_audit_log_export Data Source - oneuptime"
subcategory: "Other"
description: |-
  Where the audit log is exported: S3 or S3-compatible storage, Google Cloud Storage, Azure Blob Storage, a webhook or a syslog receiver.
---

# oneuptime_audit_log_export (Data Source)

Where the audit log is exported: S3 or S3-compatible storage, Google Cloud Storage, Azure Blob Storage, a webhook or a syslog receiver.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one audit log export may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_audit_log_export" "example" {
  name = "Example audit log export"
}

# Or by id:
data "oneuptime_audit_log_export" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `consecutive_failure_count` (Number) How many deliveries in a row have failed. It sets how long the next retry waits.
- `created_by_user_id` (String) ID of the user who created this export. The ID of a `oneuptime_user` (see the data source).
- `cursor_entry_id` (String) The ID of the last entry delivered, which orders the entries recorded in the same second.
- `destination_type` (String) Where entries are sent: s3, s3-compatible, google-cloud-storage, azure-blob-storage, webhook or syslog. Fixed once created.
- `exported_entry_count` (Number) How many entries have been delivered so far.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether entries are exported. A paused export keeps its place.
- `last_error` (String) Why the last delivery failed, with credentials redacted. Cleared by the next delivery that succeeds.
- `name` (String) A name for this export, such as Compliance archive.

### Read-Only

- `config` (String) The destination's settings that are not secret, such as the bucket and the region. Keys are defined by the destination. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `cursor_created_at` (String) When the last entry delivered was recorded: where the next batch starts.
- `failing_since` (String) When deliveries started failing. Cleared by the next delivery that succeeds.
- `last_attempt_at` (String) When OneUptime last tried to deliver entries.
- `last_exported_at` (String) When entries were last delivered.
- `next_attempt_at` (String) While deliveries fail: when OneUptime tries again. Saving the export's settings tries again at once.
- `project_id` (String) ID of the project whose audit log is exported. Empty for an export of the whole instance. The ID of a `oneuptime_project`.
- `starts_at` (String) Entries recorded from this time on are exported. Empty on create means from now; an earlier date also sends the entries kept since then.
- `updated_at` (String) Date and Time when the object was updated.
