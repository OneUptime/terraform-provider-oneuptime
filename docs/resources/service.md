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
  name        = "Example service"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_archived` (Boolean) Is this service archived? Archived services are hidden from lists but keep collecting telemetry. Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `metric_cardinality_budget` (Number) Max number of distinct metric series this service may emit per metric. When exceeded, the highest-cardinality attribute is auto-bucketed. Null inherits the project default.
- `metric_downsampling_retention_days` (String) Per-tier retention override (raw, 1m, 5m, 1h, 1d) in days. Null fields inherit the project default. A JSON value: write it with `jsonencode()`.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this service. Leave blank to use the project-wide default.
- `service_color` (String) Color for this service.
- `service_language` (String) Language in which this service is written.
- `tech_stack` (String) Tech stack used in the service. This will help other developers understand the service better. A JSON value: write it with `jsonencode()`.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this service (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the service default, then the project's retention settings. A JSON value: write it with `jsonencode()`.

### Read-Only

- `archived_at` (String) When was this service archived?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `cloud_account_id` (String) Last-seen value of the cloud.account.id OpenTelemetry resource attribute.
- `cloud_platform` (String) Last-seen value of the cloud.platform OpenTelemetry resource attribute, e.g. aws_ecs, gcp_cloud_run, aws_lambda.
- `cloud_provider` (String) Last-seen value of the cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure.
- `cloud_region` (String) Last-seen value of the cloud.region OpenTelemetry resource attribute, e.g. us-east-1.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `deployment_environment` (String) Last-seen value of the deployment.environment.name (or deployment.environment) OpenTelemetry resource attribute, e.g. production, staging.
- `id` (String) Unique identifier for the resource.
- `last_seen_at` (String) When telemetry was last received for this service.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `runtime_name` (String) Last-seen value of the process.runtime.name OpenTelemetry resource attribute, e.g. nodejs, go, OpenJDK Runtime Environment.
- `runtime_version` (String) Last-seen value of the process.runtime.version OpenTelemetry resource attribute.
- `service_namespace` (String) Last-seen value of the service.namespace OpenTelemetry resource attribute.
- `service_version` (String) Last-seen value of the service.version OpenTelemetry resource attribute.
- `slug` (String) Friendly globally unique name for your object.
- `telemetry_sdk_language` (String) Last-seen value of the telemetry.sdk.language OpenTelemetry resource attribute, e.g. java, dotnet, nodejs, python, go. Drives technology-specific golden metrics on the service overview.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing service by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_service.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_service.example <id>
```
