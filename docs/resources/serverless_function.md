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
  name = "Example short text"
  function_identifier = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Friendly name for this serverless function..
- `function_identifier` (String) Stable identifier from the faas.name OpenTelemetry resource attribute. Identity key for this function...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this function. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this function. Unset fields fall back to the function default, then the project's retention settings...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this serverless function archived? Archived serverless functions are hidden from lists but keep collecting telemetry...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `cloud_platform` (String) Last-seen cloud.platform OpenTelemetry resource attribute, e.g. aws_lambda, gcp_cloud_functions, azure_functions...
- `cloud_provider` (String) Last-seen cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure...
- `cloud_region` (String) Last-seen cloud.region OpenTelemetry resource attribute, e.g. us-east-1...
- `cloud_account_id` (String) Last-seen cloud.account.id OpenTelemetry resource attribute...
- `function_version` (String) Last-seen faas.version OpenTelemetry resource attribute...
- `runtime_name` (String) Last-seen process.runtime.name OpenTelemetry resource attribute...
- `runtime_version` (String) Last-seen process.runtime.version OpenTelemetry resource attribute...
- `otel_collector_status` (String) Whether telemetry is currently being received (connected) or has gone stale (disconnected)...
- `agent_version` (String) Version of the OneUptime agent reporting this function...
- `last_seen_at` (String) A date time object..
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_serverless_function.example <id>
```
