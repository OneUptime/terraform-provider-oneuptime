---
page_title: "oneuptime_telemetry_ingestion_key Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Manage Telemetry Ingestion Keys for your project
---

# oneuptime_telemetry_ingestion_key (Data Source)

Manage Telemetry Ingestion Keys for your project Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_telemetry_ingestion_key" "by_name" {
  name = "example-telemetry_ingestion_key"
}

data "oneuptime_telemetry_ingestion_key" "by_id" {
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
- `description` (String) Friendly description that will help you remember.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `secret_key` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `key_type` (String) Server keys are for backend services and OpenTelemetry collectors: full ingest, no origin checks. Browser keys are write-only client keys: listed web origins may send traces, logs, metrics and session replay, while listed app:// identities currently authorize React Native session replay only. This cannot be changed after the key is created - create a new key instead... Computed.
- `allowed_origins` (String) Web origins (for example https://app.example.com or https://*.example.com) and exact React Native identities (for example app://com.example.mobile) that may use this key. Required on a Browser key. Web requests need a listed Origin; mobile replay requests without Origin need a listed app identity. app:// entries cannot use wildcards and are self-asserted identifiers, not platform attestation. Ignored on a Server key... Computed.
- `pinned_service_name` (String) When set, every OpenTelemetry resource ingested with this key has its service.name REPLACED with this value. This is what stops data written with a scraped key from masquerading as another service: forged spans land in one service you can see and mute, instead of poisoning your backend services' dashboards and alerts... Computed.
- `is_enabled` (Bool) Turn this off to immediately stop accepting telemetry written with this key, without deleting it. Turn it back on to resume... Computed.
- `expires_at` (String) A date time object.. Computed.
- `last_used_at` (String) A date time object.. Computed.
- `requests_per_minute_limit` (Number) Maximum ingest requests per minute accepted with this key. Leave empty to use the shipped default for a Browser key, and to leave a Server key unlimited. The limit is per key, across every client using it, so it has to clear your whole fleet - see DEFAULT_BROWSER_KEY_REQUESTS_PER_MINUTE for the default and the reasoning behind its size... Computed.
