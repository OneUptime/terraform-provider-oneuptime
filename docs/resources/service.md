---
page_title: "oneuptime_service Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Services is a collection of services that you have in your organization. It can be a collection of services that you are monitoring or services that you are providing to your customers. It can be anything that you want to keep track of.
---

# oneuptime_service (Resource)

Services is a collection of services that you have in your organization. It can be a collection of services that you are monitoring or services that you are providing to your customers. It can be anything that you want to keep track of.

## Example Usage

```terraform
resource "oneuptime_service" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this service archived? Archived services are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `service_color` (String) Color object.
- `service_language` (String) Language in which this service is written.
- `tech_stack` (String) Tech stack used in the service. This will help other developers understand the service better...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this service. Leave blank to use the project-wide default...
- `metric_cardinality_budget` (Number) Max number of distinct metric series this service may emit per metric. When exceeded, the highest-cardinality attribute is auto-bucketed. Null inherits the project default...
- `metric_downsampling_retention_days` (String) Per-tier retention override (raw, 1m, 5m, 1h, 1d) in days. Null fields inherit the project default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this service (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the service default, then the project's retention settings...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `last_seen_at` (String) A date time object..
- `service_version` (String) Last-seen value of the service.version OpenTelemetry resource attribute...
- `deployment_environment` (String) Last-seen value of the deployment.environment.name (or deployment.environment) OpenTelemetry resource attribute, e.g. production, staging...
- `service_namespace` (String) Last-seen value of the service.namespace OpenTelemetry resource attribute...
- `runtime_name` (String) Last-seen value of the process.runtime.name OpenTelemetry resource attribute, e.g. nodejs, go, OpenJDK Runtime Environment...
- `runtime_version` (String) Last-seen value of the process.runtime.version OpenTelemetry resource attribute...
- `telemetry_sdk_language` (String) Last-seen value of the telemetry.sdk.language OpenTelemetry resource attribute, e.g. java, dotnet, nodejs, python, go. Drives technology-specific golden metrics on the service overview...
- `cloud_provider` (String) Last-seen value of the cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure...
- `cloud_platform` (String) Last-seen value of the cloud.platform OpenTelemetry resource attribute, e.g. aws_ecs, gcp_cloud_run, aws_lambda...
- `cloud_region` (String) Last-seen value of the cloud.region OpenTelemetry resource attribute, e.g. us-east-1...
- `cloud_account_id` (String) Last-seen value of the cloud.account.id OpenTelemetry resource attribute...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_service.example <id>
```
