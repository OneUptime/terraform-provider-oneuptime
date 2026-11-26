---
page_title: "oneuptime_security_event_connection Resource - oneuptime"
subcategory: "Security"
description: |-
  Managed connections to SIEM, EDR, cloud security and identity products. Records are polled on an interval and ingested as security events.
---

# oneuptime_security_event_connection (Resource)

Managed connections to SIEM, EDR, cloud security and identity products. Records are polled on an interval and ingested as security events.

## Example Usage

```terraform
resource "oneuptime_security_event_connection" "example" {
  name           = "Example security event connection"
  provider_value = "microsoft-sentinel"
  secrets        = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  description    = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this connection.
- `provider_value` (String) Which security product this connection polls, e.g. 'microsoft-sentinel' or 'crowdstrike-falcon'. Fixed once created. Allowed values: `microsoft-sentinel`, `microsoft-defender-xdr`, `crowdstrike-falcon`, `splunk`, `elastic-security`, `aws-security-hub`, `okta`, `google-secops`, `microsoft-entra-id`, `google-workspace`, `cloudflare`.
- `secrets` (String, Sensitive) Provider-specific secrets (client secret, API token, secret access key) as a JSON object. Encrypted at rest and never returned by the API.

### Optional

- `alerting_only` (Boolean) For providers that distinguish alerting from non-alerting records: import only the alerting ones. Defaults to `true`.
- `case_minimum_severity` (String) Imported alerts at or above this severity (Low, Medium, High or Critical) open security cases. Empty imports them as events only. Defaults to `High`.
- `config` (String) Provider-specific, non-secret settings such as tenant, workspace, region or base URL. Keys are defined by the provider catalog. A JSON value: write it with `jsonencode()`.
- `description` (String) What this connection imports and why.
- `is_enabled` (Boolean) Whether this connection is polled on its schedule. Defaults to `true`.
- `poll_interval_in_minutes` (Number) How often new records are polled, in minutes. Defaults to `5`.

### Read-Only

- `blocking_settings` (String) Whether OneUptime enforces your blocklist through this connection, and in which zones. Changed through the connection's blocking routes. A JSON value: read it with `jsondecode()`.
- `blocking_status` (String) What OneUptime last found of its lists and rules at the provider, written by OneUptime. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `cursor` (String) Poll cursor: the end of the last completely processed creation-time window, as an ISO string.
- `id` (String) Unique identifier for the resource.
- `last_error` (String) The most recent poll error with credentials redacted, if any. Cleared on the next successful poll.
- `last_event_ingested_at` (String) When a new security event was last imported from this connection.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `last_poll_result` (String) The latest scheduled or on-demand poll result with counts and warnings. A JSON value: read it with `jsondecode()`.
- `last_polled_at` (String) When this connection was last polled. Null means it has never run.
- `last_successful_poll_at` (String) When a complete poll last finished successfully, including an empty result.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security event connection by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_event_connection.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_event_connection.example <id>
```
