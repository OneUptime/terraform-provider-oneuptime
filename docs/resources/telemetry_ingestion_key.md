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
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `key_type` (String) Server keys are for backend services and OpenTelemetry collectors: full ingest, no origin checks. Browser keys are meant to be published in a web page, so they are write-only, restricted to trace / log / metric / session replay ingest, and are only accepted from the origins you list below. This cannot be changed after the key is created - create a new key instead...
- `allowed_origins` (String) Browser origins (scheme + host + port, for example https://app.example.com, or https://*.example.com for one level of subdomain) that may use this key. Required and strictly enforced on a Browser key: a request from an unlisted origin, or with no Origin header at all, is refused. Ignored entirely on a Server key...
- `pinned_service_name` (String) When set, every OpenTelemetry resource ingested with this key has its service.name REPLACED with this value. This is what stops data written with a scraped key from masquerading as another service: forged spans land in one service you can see and mute, instead of poisoning your backend services' dashboards and alerts...
- `is_enabled` (Bool) Turn this off to immediately stop accepting telemetry written with this key, without deleting it. Turn it back on to resume...
- `expires_at` (String) A date time object..
- `requests_per_minute_limit` (Number) Maximum ingest requests per minute accepted with this key. Leave empty to use the shipped default for a Browser key, and to leave a Server key unlimited. The limit is per key, across every client using it, so it has to clear your whole fleet - see DEFAULT_BROWSER_KEY_REQUESTS_PER_MINUTE for the default and the reasoning behind its size...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `secret_key` (String) A unique identifier for an object, represented as a UUID..
- `last_used_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_telemetry_ingestion_key.example <id>
```
