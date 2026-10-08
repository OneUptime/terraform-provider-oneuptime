---
page_title: "oneuptime_serverless_function Data Source - oneuptime"
subcategory: "Other"
description: |-
  Serverless / FaaS functions auto-discovered from OpenTelemetry (faas.name / cloud.platform). Examples: AWS Lambda, Google Cloud Functions, Azure Functions.
---

# oneuptime_serverless_function (Data Source)

Serverless / FaaS functions auto-discovered from OpenTelemetry (faas.name / cloud.platform). Examples: AWS Lambda, Google Cloud Functions, Azure Functions.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one serverless function may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_serverless_function" "example" {
  name = "Example serverless function"
}

# Or by id:
data "oneuptime_serverless_function" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `agent_version` (String) Version of the OneUptime agent reporting this function.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `cloud_account_id` (String) Last-seen cloud.account.id OpenTelemetry resource attribute.
- `cloud_platform` (String) Last-seen cloud.platform OpenTelemetry resource attribute, e.g. aws_lambda, gcp_cloud_functions, azure_functions.
- `cloud_provider` (String) Last-seen cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure.
- `cloud_region` (String) Last-seen cloud.region OpenTelemetry resource attribute, e.g. us-east-1.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `function_identifier` (String) Stable identifier from the faas.name OpenTelemetry resource attribute. Identity key for this function.
- `function_version` (String) Last-seen faas.version OpenTelemetry resource attribute.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_archived` (Boolean) Is this serverless function archived? Archived serverless functions are hidden from lists but keep collecting telemetry.
- `name` (String) Friendly name for this serverless function.
- `otel_collector_status` (String) Whether telemetry is currently being received (connected) or has gone stale (disconnected).
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this function. Leave blank to use the project-wide default.
- `runtime_name` (String) Last-seen process.runtime.name OpenTelemetry resource attribute.
- `runtime_version` (String) Last-seen process.runtime.version OpenTelemetry resource attribute.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `archived_at` (String) When was this serverless function archived?
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When telemetry was last received for this function.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this function. Unset fields fall back to the function default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
