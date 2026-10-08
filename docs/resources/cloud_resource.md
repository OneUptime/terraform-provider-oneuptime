---
page_title: "oneuptime_cloud_resource Resource - oneuptime"
subcategory: "Other"
description: |-
  Cloud environments - managed compute auto-discovered from OpenTelemetry cloud.platform (e.g. AWS ECS/Fargate, GCP Cloud Run, Azure Container Apps, Elastic Beanstalk, App Runner) - and cloud resources: the IaaS and PaaS resources (virtual machines, load balancers, buckets, managed databases, queues, ...) discovered from the metrics Azure Monitor, Amazon CloudWatch and Google Cloud Monitoring publish about them.
---

# oneuptime_cloud_resource (Resource)

Cloud environments - managed compute auto-discovered from OpenTelemetry cloud.platform (e.g. AWS ECS/Fargate, GCP Cloud Run, Azure Container Apps, Elastic Beanstalk, App Runner) - and cloud resources: the IaaS and PaaS resources (virtual machines, load balancers, buckets, managed databases, queues, ...) discovered from the metrics Azure Monitor, Amazon CloudWatch and Google Cloud Monitoring publish about them.

## Example Usage

```terraform
resource "oneuptime_cloud_resource" "example" {
  name                = "Example cloud resource"
  resource_identifier = "Example short text"
  description         = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this cloud environment. Ingest names a discovered environment after its platform, region and account.
- `resource_identifier` (String) Environment key: the cloud.platform, cloud.account.id and cloud.region OpenTelemetry resource attributes joined with '|' (e.g. aws_ecs|123456789012|us-east-1; missing parts stay as empty segments). Built by buildCloudEnvironmentKey in Common/Types/Cloud/CloudPlatform. An environment created by hand must carry the same key for ingest to find it instead of creating a duplicate.

### Optional

- `cloud_account_id` (String) Last-seen cloud.account.id OpenTelemetry resource attribute.
- `cloud_platform` (String) Last-seen cloud.platform OpenTelemetry resource attribute, e.g. aws_ecs, gcp_cloud_run, azure_container_apps.
- `cloud_provider` (String) Last-seen cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure.
- `cloud_region` (String) Last-seen cloud.region OpenTelemetry resource attribute, e.g. us-east-1.
- `description` (String) Friendly description that will help you remember.
- `is_archived` (Boolean) Is this cloud resource archived? Archived cloud resources are hidden from lists but keep collecting telemetry. Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this resource. Leave blank to use the project-wide default.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this resource. Unset fields fall back to the resource default, then the project's retention settings. A JSON value: write it with `jsonencode()`.

### Read-Only

- `agent_version` (String) Version of the OneUptime agent reporting this resource.
- `archived_at` (String) When was this cloud resource archived?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `auto_archived_at` (String) For a resource: when it was archived automatically for sending no metrics for the auto-archive period. Cleared when it reports again.
- `cloud_resource_group` (String) For an Azure resource: the resource group it belongs to.
- `cloud_resource_kind` (String) environment: a managed compute environment discovered from the cloud.platform, cloud.account.id and cloud.region resource attributes of a workload's own telemetry. resource: one IaaS or PaaS resource (a virtual machine, a load balancer, a bucket, a managed database, ...) discovered from the metrics Azure Monitor, Amazon CloudWatch or Google Cloud Monitoring publish about it.
- `cloud_resource_type` (String) For a resource: the provider's type for it - the Azure Resource Manager type (Microsoft.Compute/virtualMachines), the AWS CloudFormation type (AWS::EC2::Instance) or the Google Cloud Monitoring resource type (gce_instance).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_seen_at` (String) When telemetry was last received for this resource.
- `otel_collector_status` (String) Whether telemetry is currently being received (connected) or has gone stale (disconnected).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `provider_resource_id` (String) For a resource: the provider's id for it - its Azure resource id, its AWS ARN or its Google Cloud full resource name. Where the metrics do not name the resource completely (an AWS resource whose ARN needs an id no metric reports), a readable composite of what they do name.
- `runtime_name` (String) Last-seen process.runtime.name OpenTelemetry resource attribute.
- `runtime_version` (String) Last-seen process.runtime.version OpenTelemetry resource attribute.
- `slug` (String) Friendly globally unique name for your object.
- `telemetry_attributes` (String) For a resource: the metric attributes, exactly as stored, that select its metrics - for example azuremonitor.resource_id, or the CloudWatch Namespace and identifying Dimensions with the account and region. The resource's pages and the monitors created from them filter on these. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing cloud resource by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_cloud_resource.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_cloud_resource.example <id>
```
