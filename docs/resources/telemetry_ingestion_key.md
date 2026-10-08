---
page_title: "oneuptime_telemetry_ingestion_key Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Manage Telemetry Ingestion Keys for your project
---

# oneuptime_telemetry_ingestion_key (Resource)

Manage Telemetry Ingestion Keys for your project

## Example Usage

```terraform
resource "oneuptime_telemetry_ingestion_key" "example" {
  name        = "Example telemetry ingestion key"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `allowed_origins` (String) Web origins (for example https://app.example.com or https://*.example.com) and exact React Native identities (for example app://com.example.mobile) that may use this key. Required on a Browser key. Web requests need a listed Origin; mobile replay requests without Origin need a listed app identity. app:// entries cannot use wildcards and are self-asserted identifiers, not platform attestation. Ignored on a Server key. A JSON value: write it with `jsonencode()`.
- `description` (String) Friendly description that will help you remember.
- `expires_at` (String) Date and time after which this key stops being accepted. Empty means it never expires, which is how every key behaved before this column existed. Setting one on a Browser key bounds how long a scraped copy stays useful.
- `is_enabled` (Boolean) Turn this off to immediately stop accepting telemetry written with this key, without deleting it. Turn it back on to resume. Defaults to `true`.
- `key_type` (String) Server keys are for backend services and OpenTelemetry collectors: full ingest, no origin checks. Browser keys are write-only client keys: listed web origins may send traces, logs, metrics and session replay, while listed app:// identities currently authorize React Native session replay only. This cannot be changed after the key is created - create a new key instead. Defaults to `Server`.
- `pinned_service_name` (String) When set, every OpenTelemetry resource ingested with this key has its service.name REPLACED with this value. This is what stops data written with a scraped key from masquerading as another service: forged spans land in one service you can see and mute, instead of poisoning your backend services' dashboards and alerts.
- `requests_per_minute_limit` (Number) Maximum ingest requests per minute accepted with this key. Leave empty to use the shipped default for a Browser key, and to leave a Server key unlimited. The limit is per key, across every client using it, so it has to clear your whole fleet - see DEFAULT_BROWSER_KEY_REQUESTS_PER_MINUTE for the default and the reasoning behind its size.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_used_at` (String) The last time telemetry was accepted with this key. Empty means it has never been used since this was recorded. Use it to find keys that are safe to rotate or delete.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `secret_key` (String) Secret Telemetry Ingestion Key.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing telemetry ingestion key by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_telemetry_ingestion_key.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_telemetry_ingestion_key.example <id>
```
