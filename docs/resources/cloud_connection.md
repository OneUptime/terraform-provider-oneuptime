---
page_title: "oneuptime_cloud_connection Resource - oneuptime"
subcategory: "Other"
description: |-
  Connections to AWS, Azure and Google Cloud accounts: OneUptime lists their resources and reads their key metrics on a schedule, with read-only access and nothing to install.
---

# oneuptime_cloud_connection (Resource)

Connections to AWS, Azure and Google Cloud accounts: OneUptime lists their resources and reads their key metrics on a schedule, with read-only access and nothing to install.

## Example Usage

```terraform
resource "oneuptime_cloud_connection" "example" {
  name           = "Example cloud connection"
  provider_value = "Example short text"
  secrets        = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `name` (String) Friendly name for this connection.
- `provider_value` (String) Which cloud this connection reads: 'aws', 'azure' or 'gcp'. Fixed once created.
- `secrets` (String, Sensitive) The cloud's secrets (a secret access key, a client secret, a service account key) as a JSON object. Encrypted at rest and never returned by the API.

### Optional

- `config` (String) The cloud's settings that are not secret, such as a role ARN, a tenant or the regions to read. Keys are defined by the cloud. A JSON value: write it with `jsonencode()`.
- `is_enabled` (Boolean) Whether this connection is read on its schedule. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_error` (String) Why the last read failed, with credentials removed. Cleared by the next read that works.
- `last_pull_result` (String) What the last read found: the accounts, the resources by type and region, and the metrics it wrote. A JSON value: write it with `jsonencode()`.
- `last_pulled_at` (String) When the connection was last read. Empty means never.
- `last_successful_pull_at` (String) When a read last finished without an error.
- `next_pull_at` (String) When the connection is read next. Empty means as soon as possible.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `service_id` (String) ID of the telemetry service the metrics this connection reads belong to. The ID of a `oneuptime_service`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing cloud connection by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_cloud_connection.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_cloud_connection.example <id>
```
