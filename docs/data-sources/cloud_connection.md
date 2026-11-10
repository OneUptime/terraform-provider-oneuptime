---
page_title: "oneuptime_cloud_connection Data Source - oneuptime"
subcategory: "Cloud"
description: |-
  Connections to AWS, Azure and Google Cloud accounts: OneUptime lists their resources and reads their key metrics on a schedule, with read-only access and nothing to install.
---

# oneuptime_cloud_connection (Data Source)

Connections to AWS, Azure and Google Cloud accounts: OneUptime lists their resources and reads their key metrics on a schedule, with read-only access and nothing to install.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one cloud connection may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_cloud_connection" "example" {
  name = "Example cloud connection"
}

# Or by id:
data "oneuptime_cloud_connection" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this connection is read on its schedule.
- `last_error` (String) Why the last read failed, with credentials removed. Cleared by the next read that works.
- `name` (String) Friendly name for this connection.
- `provider_value` (String) Which cloud this connection reads: 'aws', 'azure' or 'gcp'. Fixed once created.
- `service_id` (String) ID of the telemetry service the metrics this connection reads belong to. The ID of a `oneuptime_service`.

### Read-Only

- `config` (String) The cloud's settings that are not secret, such as a role ARN, a tenant or the regions to read. Keys are defined by the cloud. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_pull_result` (String) What the last read found: the accounts, the resources by type and region, and the metrics it wrote. A JSON value: read it with `jsondecode()`.
- `last_pulled_at` (String) When the connection was last read. Empty means never.
- `last_successful_pull_at` (String) When a read last finished without an error.
- `next_pull_at` (String) When the connection is read next. Empty means as soon as possible.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `stream_key_created_at` (String) When the connection's stream key was made. Empty means it streams nothing.
- `stream_state` (String) Azure: the subscription holding the connection's event hubs, and where OneUptime has read each to. A JSON value: read it with `jsondecode()`.
- `stream_summary` (String) What the connection's streams delivered: the accounts, regions and namespaces or resource types of metrics, the log groups, categories or log names of logs, and the last delivery refused. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
