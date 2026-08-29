---
page_title: "oneuptime_cloud_resource Resource - oneuptime"
subcategory: "Other"
description: |-
  Managed cloud compute auto-discovered from OpenTelemetry cloud.platform (e.g. AWS ECS/Fargate, GCP Cloud Run, Azure Container Apps, Elastic Beanstalk, App Runner).
---

# oneuptime_cloud_resource (Resource)

Managed cloud compute auto-discovered from OpenTelemetry cloud.platform (e.g. AWS ECS/Fargate, GCP Cloud Run, Azure Container Apps, Elastic Beanstalk, App Runner).

## Example Usage

```terraform
resource "oneuptime_cloud_resource" "example" {
  name = "Example short text"
  resource_identifier = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Friendly name for this cloud resource..
- `resource_identifier` (String) Stable identifier for this managed-compute workload (service.name, falling back to host.name). Identity key for this resource...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this resource. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this resource. Unset fields fall back to the resource default, then the project's retention settings...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this cloud resource archived? Archived cloud resources are hidden from lists but keep collecting telemetry...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `cloud_platform` (String) Last-seen cloud.platform OpenTelemetry resource attribute, e.g. aws_ecs, gcp_cloud_run, azure_container_apps...
- `cloud_provider` (String) Last-seen cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure...
- `cloud_region` (String) Last-seen cloud.region OpenTelemetry resource attribute, e.g. us-east-1...
- `cloud_account_id` (String) Last-seen cloud.account.id OpenTelemetry resource attribute...
- `runtime_name` (String) Last-seen process.runtime.name OpenTelemetry resource attribute...
- `runtime_version` (String) Last-seen process.runtime.version OpenTelemetry resource attribute...
- `otel_collector_status` (String) Whether telemetry is currently being received (connected) or has gone stale (disconnected)...
- `agent_version` (String) Version of the OneUptime agent reporting this resource...
- `last_seen_at` (String) A date time object..
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_cloud_resource.example <id>
```
