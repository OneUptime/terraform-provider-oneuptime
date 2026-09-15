---
page_title: "oneuptime_security_event_connection Resource - oneuptime"
subcategory: "Other"
description: |-
  Managed connections to SIEM, EDR, cloud security and identity products. Records are polled on an interval and ingested as security events.
---

# oneuptime_security_event_connection (Resource)

Managed connections to SIEM, EDR, cloud security and identity products. Records are polled on an interval and ingested as security events.

## Example Usage

```terraform
resource "oneuptime_security_event_connection" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  provider = "Example short text"
  secrets = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `provider` (String) Which security product this connection polls, e.g. 'microsoft-sentinel' or 'crowdstrike-falcon'. Fixed once created...
- `secrets` (String) Provider-specific secrets (client secret, API token, secret access key) as a JSON object. Encrypted at rest and never returned by the API...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) What this connection imports and why...
- `config` (String) Provider-specific, non-secret settings such as tenant, workspace, region or base URL. Keys are defined by the provider catalog...
- `is_enabled` (Bool) Whether this connection is polled on its schedule...
- `poll_interval_in_minutes` (Number) How often new records are polled, in minutes...
- `alerting_only` (Bool) For providers that distinguish alerting from non-alerting records: import only the alerting ones...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `last_successful_poll_at` (String) A date time object..
- `last_event_ingested_at` (String) A date time object..
- `last_poll_result` (String) The latest scheduled or on-demand poll result with counts and warnings...
- `last_polled_at` (String) A date time object..
- `cursor` (String) Poll cursor: the end of the last completely processed creation-time window, as an ISO string...
- `last_error` (String) The most recent poll error with credentials redacted, if any. Cleared on the next successful poll...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_security_event_connection.example <id>
```
