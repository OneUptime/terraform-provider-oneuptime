---
page_title: "oneuptime_serverless_function Resource - oneuptime"
subcategory: "Other"
description: |-
  Serverless / FaaS functions auto-discovered from OpenTelemetry (faas.name / cloud.platform). Examples: AWS Lambda, Google Cloud Functions, Azure Functions.
---

# oneuptime_serverless_function (Resource)

Serverless / FaaS functions auto-discovered from OpenTelemetry (faas.name / cloud.platform). Examples: AWS Lambda, Google Cloud Functions, Azure Functions.

## Example Usage

```terraform
resource "oneuptime_serverless_function" "example" {
  name                = "Example serverless function"
  function_identifier = "Example short text"
  description         = "Managed by Terraform"
}
```

## Schema

### Required

- `function_identifier` (String) Stable identifier from the faas.name OpenTelemetry resource attribute. Identity key for this function.
- `name` (String) Friendly name for this serverless function.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_archived` (Boolean) Is this serverless function archived? Archived serverless functions are hidden from lists but keep collecting telemetry. Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this function. Leave blank to use the project-wide default.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this function. Unset fields fall back to the function default, then the project's retention settings. A JSON value: write it with `jsonencode()`.

### Read-Only

- `agent_version` (String) Version of the OneUptime agent reporting this function.
- `archived_at` (String) When was this serverless function archived?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `cloud_account_id` (String) Last-seen cloud.account.id OpenTelemetry resource attribute.
- `cloud_platform` (String) Last-seen cloud.platform OpenTelemetry resource attribute, e.g. aws_lambda, gcp_cloud_functions, azure_functions.
- `cloud_provider` (String) Last-seen cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure.
- `cloud_region` (String) Last-seen cloud.region OpenTelemetry resource attribute, e.g. us-east-1.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `function_version` (String) Last-seen faas.version OpenTelemetry resource attribute.
- `id` (String) Unique identifier for the resource.
- `last_seen_at` (String) When telemetry was last received for this function.
- `otel_collector_status` (String) Whether telemetry is currently being received (connected) or has gone stale (disconnected).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `runtime_name` (String) Last-seen process.runtime.name OpenTelemetry resource attribute.
- `runtime_version` (String) Last-seen process.runtime.version OpenTelemetry resource attribute.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing serverless function by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_serverless_function.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_serverless_function.example <id>
```
