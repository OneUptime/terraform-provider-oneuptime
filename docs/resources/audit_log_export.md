---
page_title: "oneuptime_audit_log_export Resource - oneuptime"
subcategory: "Other"
description: |-
  Where the audit log is exported: S3 or S3-compatible storage, Google Cloud Storage, Azure Blob Storage, a webhook or a syslog receiver.
---

# oneuptime_audit_log_export (Resource)

Where the audit log is exported: S3 or S3-compatible storage, Google Cloud Storage, Azure Blob Storage, a webhook or a syslog receiver.

## Example Usage

```terraform
resource "oneuptime_audit_log_export" "example" {
  name             = "Example audit log export"
  destination_type = "Example short text"
}
```

## Schema

### Required

- `destination_type` (String) Where entries are sent: s3, s3-compatible, google-cloud-storage, azure-blob-storage, webhook or syslog. Fixed once created.
- `name` (String) A name for this export, such as Compliance archive.

### Optional

- `config` (String) The destination's settings that are not secret, such as the bucket and the region. Keys are defined by the destination. A JSON value: write it with `jsonencode()`.
- `is_enabled` (Boolean) Whether entries are exported. A paused export keeps its place. Defaults to `true`.
- `secrets` (String, Sensitive) The destination's secrets, such as an access key's secret, as a JSON object. Encrypted at rest and never returned by the API.
- `starts_at` (String) Entries recorded from this time on are exported. Empty on create means from now; an earlier date also sends the entries kept since then.

### Read-Only

- `consecutive_failure_count` (Number) How many deliveries in a row have failed. It sets how long the next retry waits.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this export. The ID of a `oneuptime_user` (see the data source).
- `cursor_created_at` (String) When the last entry delivered was recorded: where the next batch starts.
- `cursor_entry_id` (String) The ID of the last entry delivered, which orders the entries recorded in the same second.
- `exported_entry_count` (Number) How many entries have been delivered so far.
- `failing_since` (String) When deliveries started failing. Cleared by the next delivery that succeeds.
- `id` (String) Unique identifier for the resource.
- `last_attempt_at` (String) When OneUptime last tried to deliver entries.
- `last_error` (String) Why the last delivery failed, with credentials redacted. Cleared by the next delivery that succeeds.
- `last_exported_at` (String) When entries were last delivered.
- `next_attempt_at` (String) While deliveries fail: when OneUptime tries again. Saving the export's settings tries again at once.
- `project_id` (String) ID of the project whose audit log is exported. Empty for an export of the whole instance. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing audit log export by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_audit_log_export.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_audit_log_export.example <id>
```
