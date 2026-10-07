---
page_title: "oneuptime_cloud_resource Data Source - oneuptime"
subcategory: "Other"
description: |-
  Cloud environments - managed compute auto-discovered from OpenTelemetry cloud.platform (e.g. AWS ECS/Fargate, GCP Cloud Run, Azure Container Apps, Elastic Beanstalk, App Runner) - and cloud resources: the IaaS and PaaS resources (virtual machines, load balancers, buckets, managed databases, queues, ...) discovered from the metrics Azure Monitor, Amazon CloudWatch and Google Cloud Monitoring publish about them.
---

# oneuptime_cloud_resource (Data Source)

Cloud environments - managed compute auto-discovered from OpenTelemetry cloud.platform (e.g. AWS ECS/Fargate, GCP Cloud Run, Azure Container Apps, Elastic Beanstalk, App Runner) - and cloud resources: the IaaS and PaaS resources (virtual machines, load balancers, buckets, managed databases, queues, ...) discovered from the metrics Azure Monitor, Amazon CloudWatch and Google Cloud Monitoring publish about them. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_cloud_resource" "by_name" {
  name = "example-cloud_resource"
}

data "oneuptime_cloud_resource" "by_id" {
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
- `resource_identifier` (String) Environment key: the cloud.platform, cloud.account.id and cloud.region OpenTelemetry resource attributes joined with '|' (e.g. aws_ecs|123456789012|us-east-1; missing parts stay as empty segments). Built by buildCloudEnvironmentKey in Common/Types/Cloud/CloudPlatform. An environment created by hand must carry the same key for ingest to find it instead of creating a duplicate... Computed.
- `cloud_platform` (String) Last-seen cloud.platform OpenTelemetry resource attribute, e.g. aws_ecs, gcp_cloud_run, azure_container_apps... Computed.
- `cloud_provider` (String) Last-seen cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure... Computed.
- `cloud_region` (String) Last-seen cloud.region OpenTelemetry resource attribute, e.g. us-east-1... Computed.
- `cloud_account_id` (String) Last-seen cloud.account.id OpenTelemetry resource attribute... Computed.
- `cloud_resource_kind` (String) environment: a managed compute environment discovered from the cloud.platform, cloud.account.id and cloud.region resource attributes of a workload's own telemetry. resource: one IaaS or PaaS resource (a virtual machine, a load balancer, a bucket, a managed database, ...) discovered from the metrics Azure Monitor, Amazon CloudWatch or Google Cloud Monitoring publish about it... Computed.
- `cloud_resource_type` (String) For a resource: the provider's type for it - the Azure Resource Manager type (Microsoft.Compute/virtualMachines), the AWS CloudFormation type (AWS::EC2::Instance) or the Google Cloud Monitoring resource type (gce_instance)... Computed.
- `provider_resource_id` (String) For a resource: the provider's id for it - its Azure resource id, its AWS ARN or its Google Cloud full resource name. Where the metrics do not name the resource completely (an AWS resource whose ARN needs an id no metric reports), a readable composite of what they do name... Computed.
- `cloud_resource_group` (String) For an Azure resource: the resource group it belongs to... Computed.
- `telemetry_attributes` (String) For a resource: the metric attributes, exactly as stored, that select its metrics - for example azuremonitor.resource_id, or the CloudWatch Namespace and identifying Dimensions with the account and region. The resource's pages and the monitors created from them filter on these... Computed.
- `runtime_name` (String) Last-seen process.runtime.name OpenTelemetry resource attribute... Computed.
- `runtime_version` (String) Last-seen process.runtime.version OpenTelemetry resource attribute... Computed.
- `otel_collector_status` (String) Whether telemetry is currently being received (connected) or has gone stale (disconnected)... Computed.
- `agent_version` (String) Version of the OneUptime agent reporting this resource... Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this resource. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this resource. Unset fields fall back to the resource default, then the project's retention settings... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this cloud resource archived? Archived cloud resources are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `auto_archived_at` (String) A date time object.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
