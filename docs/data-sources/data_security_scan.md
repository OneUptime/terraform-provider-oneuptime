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
  data_store_id = oneuptime_data_store.example.id
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
- `data_store_id` (String) The store scanned, for a content or telemetry scan. The ID of a `oneuptime_data_store`.
- `duration_in_ms` (Number) How long the scan took.
- `error_message` (String) Why the scan failed, if it did.
- `findings_found` (Number) How many findings the scan made or updated.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `kind` (String) What the scan read: Configuration, Content or Telemetry.
- `objects_read` (Number) How many objects or files the scan sampled.
- `probe_id` (String) The probe that ran the scan, when one did. The ID of a `oneuptime_probe`.
- `rows_read` (Number) How many rows or records the scan sampled.
- `status` (String) Running, Succeeded, Partly succeeded, Failed or Skipped.
- `stores_read` (Number) How many stores a configuration read covered.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `finished_at` (String) When the scan finished.
- `project_id` (String) ID of the project this scan belongs to. The ID of a `oneuptime_project`.
- `read_log` (String) The objects, tables or services read, by name. Never their contents. A JSON value: read it with `jsondecode()`.
- `started_at` (String) When the scan started.
- `updated_at` (String) Date and Time when the object was updated.
- `warnings` (String) What the scan could not read, and why. A JSON value: read it with `jsondecode()`.
