---
page_title: "oneuptime_security_event_connection Data Source - oneuptime"
subcategory: "Other"
description: |-
  Managed connections to SIEM, EDR, cloud security and identity products. Records are polled on an interval and ingested as security events.
---

# oneuptime_security_event_connection (Data Source)

Managed connections to SIEM, EDR, cloud security and identity products. Records are polled on an interval and ingested as security events. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_security_event_connection" "by_name" {
  name = "example-security_event_connection"
}

data "oneuptime_security_event_connection" "by_id" {
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
- `description` (String) What this connection imports and why... Computed.
- `provider` (String) Which security product this connection polls, e.g. 'microsoft-sentinel' or 'crowdstrike-falcon'. Fixed once created... Computed.
- `config` (String) Provider-specific, non-secret settings such as tenant, workspace, region or base URL. Keys are defined by the provider catalog... Computed.
- `is_enabled` (Bool) Whether this connection is polled on its schedule... Computed.
- `poll_interval_in_minutes` (Number) How often new records are polled, in minutes... Computed.
- `alerting_only` (Bool) For providers that distinguish alerting from non-alerting records: import only the alerting ones... Computed.
- `last_successful_poll_at` (String) A date time object.. Computed.
- `last_event_ingested_at` (String) A date time object.. Computed.
- `last_poll_result` (String) The latest scheduled or on-demand poll result with counts and warnings... Computed.
- `last_polled_at` (String) A date time object.. Computed.
- `cursor` (String) Poll cursor: the end of the last completely processed creation-time window, as an ISO string... Computed.
- `last_error` (String) The most recent poll error with credentials redacted, if any. Cleared on the next successful poll... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
