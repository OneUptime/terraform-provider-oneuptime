---
page_title: "oneuptime_host Resource - oneuptime"
subcategory: "Other"
description: |-
  Hosts that are being monitored in this project. Each host is auto-discovered when an OTel Collector reports the host.name resource attribute, or can be manually registered.
---

# oneuptime_host (Resource)

Hosts that are being monitored in this project. Each host is auto-discovered when an OTel Collector reports the host.name resource attribute, or can be manually registered.

## Example Usage

```terraform
resource "oneuptime_host" "example" {
  name = "Example short text"
  host_identifier = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Friendly name for this host..
- `host_identifier` (String) Unique identifier for this host, sourced from the host.name OTel resource attribute..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this host..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this host archived? Archived hosts are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this host. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the host default, then the project's retention settings...
- `otel_collector_status` (String) Connection status of the OTel Collector reporting on this host (connected or disconnected)..
- `agent_version` (String) Version of the OneUptime agent reporting telemetry on this host, as self-reported via the oneuptime.agent.version resource attribute..
- `last_seen_at` (String) A date time object..
- `os_type` (String) Operating system type of the host..
- `os_version` (String) Operating system version of the host..
- `host_id` (String) Stable host identifier reported by the OTel host.id resource attribute..
- `host_arch` (String) CPU architecture from the OTel host.arch resource attribute..
- `host_type` (String) Cloud-instance class reported by the OTel host.type resource attribute..
- `host_ip_addresses` (String) Comma-separated list of every IP address reported by the OTel host.ip resource attribute, in the order the collector reported them, deduplicated. The Hosts list shows the most routable one (IPv4, non-loopback, non-link-local) first; the host detail page groups them all by category...
- `cpu_cores` (Number) Logical CPU core count, sourced from system.cpu.logical.count metric..
- `total_memory_bytes` (Number) Total physical memory in bytes, sourced from system.memory.usage metric (sum of all states)...
- `process_count` (Number) Most recent process count from system.processes.count metric..
- `container_runtime` (String) Container runtime detected on this host, if any (e.g. docker, containerd)..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `docker_host_id` (String) A unique identifier for an object, represented as a UUID..
- `kubernetes_cluster_id` (String) A unique identifier for an object, represented as a UUID..
- `proxmox_cluster_id` (String) A unique identifier for an object, represented as a UUID..
- `archived_at` (String) A date time object..
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deployment_environment` (String) Last-seen value of the deployment.environment.name (or deployment.environment) OpenTelemetry resource attribute, e.g. production, staging...
- `runtime_name` (String) Last-seen value of the process.runtime.name OpenTelemetry resource attribute...
- `runtime_version` (String) Last-seen value of the process.runtime.version OpenTelemetry resource attribute...
- `cloud_provider` (String) Last-seen value of the cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure...
- `cloud_platform` (String) Last-seen value of the cloud.platform OpenTelemetry resource attribute, e.g. aws_ec2, gcp_compute_engine...
- `cloud_region` (String) Last-seen value of the cloud.region OpenTelemetry resource attribute, e.g. us-east-1...
- `cloud_account_id` (String) Last-seen value of the cloud.account.id OpenTelemetry resource attribute...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_host.example <id>
```
