---
page_title: "oneuptime_google_sec_ops_connection Data Source - oneuptime"
subcategory: "Other"
description: |-
  Connections to Google SecOps (Chronicle) tenants. Detection alerts are polled on an interval and ingested as security events.
---

# oneuptime_google_sec_ops_connection (Data Source)

Connections to Google SecOps (Chronicle) tenants. Detection alerts are polled on an interval and ingested as security events. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_google_sec_ops_connection" "by_name" {
  name = "example-google_sec_ops_connection"
}

data "oneuptime_google_sec_ops_connection" "by_id" {
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
- `region` (String) Google SecOps regional endpoint prefix, e.g. 'us' or 'europe'. Used to build the API base URL... Computed.
- `instance_resource_name` (String) The Chronicle instance resource name: projects/{project}/locations/{location}/instances/{instance}... Computed.
- `is_enabled` (Bool) Whether this connection is polled... Computed.
- `poll_interval_in_minutes` (Number) How often detection alerts are polled, in minutes... Computed.
- `last_polled_at` (String) A date time object.. Computed.
- `cursor` (String) Poll cursor: the newest detection timestamp already ingested, as an ISO string... Computed.
- `last_error` (String) The most recent poll error, if any. Cleared on the next successful poll... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
