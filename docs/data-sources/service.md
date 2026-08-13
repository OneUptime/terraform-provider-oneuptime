---
page_title: "oneuptime_service Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Services is a collection of services that you have in your organization. It can be a collection of services that you are monitoring or services that you are providing to your customers. It can be anything that you want to keep track of.
---

# oneuptime_service (Data Source)

Services is a collection of services that you have in your organization. It can be a collection of services that you are monitoring or services that you are providing to your customers. It can be anything that you want to keep track of. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_service" "by_name" {
  name = "example-service"
}

data "oneuptime_service" "by_id" {
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
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this service archived? Archived services are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `service_color` (String) Color object. Computed.
- `service_language` (String) Language in which this service is written. Computed.
- `tech_stack` (String) Tech stack used in the service. This will help other developers understand the service better... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this service. Leave blank to use the project-wide default... Computed.
- `metric_cardinality_budget` (Number) Max number of distinct metric series this service may emit per metric. When exceeded, the highest-cardinality attribute is auto-bucketed. Null inherits the project default... Computed.
- `metric_downsampling_retention_days` (String) Per-tier retention override (raw, 1m, 5m, 1h, 1d) in days. Null fields inherit the project default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this service (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the service default, then the project's retention settings... Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `service_version` (String) Last-seen value of the service.version OpenTelemetry resource attribute... Computed.
- `deployment_environment` (String) Last-seen value of the deployment.environment.name (or deployment.environment) OpenTelemetry resource attribute, e.g. production, staging... Computed.
- `service_namespace` (String) Last-seen value of the service.namespace OpenTelemetry resource attribute... Computed.
- `runtime_name` (String) Last-seen value of the process.runtime.name OpenTelemetry resource attribute, e.g. nodejs, go, OpenJDK Runtime Environment... Computed.
- `runtime_version` (String) Last-seen value of the process.runtime.version OpenTelemetry resource attribute... Computed.
- `telemetry_sdk_language` (String) Last-seen value of the telemetry.sdk.language OpenTelemetry resource attribute, e.g. java, dotnet, nodejs, python, go. Drives technology-specific golden metrics on the service overview... Computed.
- `cloud_provider` (String) Last-seen value of the cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure... Computed.
- `cloud_platform` (String) Last-seen value of the cloud.platform OpenTelemetry resource attribute, e.g. aws_ecs, gcp_cloud_run, aws_lambda... Computed.
- `cloud_region` (String) Last-seen value of the cloud.region OpenTelemetry resource attribute, e.g. us-east-1... Computed.
- `cloud_account_id` (String) Last-seen value of the cloud.account.id OpenTelemetry resource attribute... Computed.
