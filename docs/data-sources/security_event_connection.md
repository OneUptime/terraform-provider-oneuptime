---
page_title: "oneuptime_security_event_connection Data Source - oneuptime"
subcategory: "Other"
description: |-
  Managed connections to SIEM, EDR, cloud security and identity products. Records are polled on an interval and ingested as security events.
---

# oneuptime_security_event_connection (Data Source)

Managed connections to SIEM, EDR, cloud security and identity products. Records are polled on an interval and ingested as security events.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security event connection may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_event_connection" "example" {
  name = "Example security event connection"
}

# Or by id:
data "oneuptime_security_event_connection" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alerting_only` (Boolean) For providers that distinguish alerting from non-alerting records: import only the alerting ones.
- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) Poll cursor: the end of the last completely processed creation-time window, as an ISO string.
- `description` (String) What this connection imports and why.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this connection is polled on its schedule.
- `last_error` (String) The most recent poll error with credentials redacted, if any. Cleared on the next successful poll.
- `name` (String) Friendly name for this connection.
- `poll_interval_in_minutes` (Number) How often new records are polled, in minutes.
- `provider_value` (String) Which security product this connection polls, e.g. 'microsoft-sentinel' or 'crowdstrike-falcon'. Fixed once created.

### Read-Only

- `config` (String) Provider-specific, non-secret settings such as tenant, workspace, region or base URL. Keys are defined by the provider catalog. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_event_ingested_at` (String) When a new security event was last imported from this connection.
- `last_poll_result` (String) The latest scheduled or on-demand poll result with counts and warnings. A JSON value: write it with `jsonencode()`.
- `last_polled_at` (String) When this connection was last polled. Null means it has never run.
- `last_successful_poll_at` (String) When a complete poll last finished successfully, including an empty result.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
