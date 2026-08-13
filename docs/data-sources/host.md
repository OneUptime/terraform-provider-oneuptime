---
page_title: "oneuptime_host Data Source - oneuptime"
subcategory: "Other"
description: |-
  Hosts that are being monitored in this project. Each host is auto-discovered when an OTel Collector reports the host.name resource attribute, or can be manually registered.
---

# oneuptime_host (Data Source)

Hosts that are being monitored in this project. Each host is auto-discovered when an OTel Collector reports the host.name resource attribute, or can be manually registered. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_host" "by_name" {
  name = "example-host"
}

data "oneuptime_host" "by_id" {
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
- `description` (String) Friendly description for this host.. Computed.
- `host_identifier` (String) Unique identifier for this host, sourced from the host.name OTel resource attribute.. Computed.
- `otel_collector_status` (String) Connection status of the OTel Collector reporting on this host (connected or disconnected).. Computed.
- `agent_version` (String) Version of the OneUptime agent reporting telemetry on this host, as self-reported via the oneuptime.agent.version resource attribute.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `os_type` (String) Operating system type of the host.. Computed.
- `os_version` (String) Operating system version of the host.. Computed.
- `host_id` (String) Stable host identifier reported by the OTel host.id resource attribute.. Computed.
- `host_arch` (String) CPU architecture from the OTel host.arch resource attribute.. Computed.
- `host_type` (String) Cloud-instance class reported by the OTel host.type resource attribute.. Computed.
- `host_ip_addresses` (String) Comma-separated list of every IP address reported by the OTel host.ip resource attribute, in the order the collector reported them, deduplicated. The Hosts list shows the most routable one (IPv4, non-loopback, non-link-local) first; the host detail page groups them all by category... Computed.
- `cpu_cores` (Number) Logical CPU core count, sourced from system.cpu.logical.count metric.. Computed.
- `total_memory_bytes` (Number) Total physical memory in bytes, sourced from system.memory.usage metric (sum of all states)... Computed.
- `process_count` (Number) Most recent process count from system.processes.count metric.. Computed.
- `container_runtime` (String) Container runtime detected on this host, if any (e.g. docker, containerd).. Computed.
- `docker_host_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `kubernetes_cluster_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `proxmox_cluster_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this host archived? Archived hosts are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this host. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this host (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the host default, then the project's retention settings... Computed.
- `deployment_environment` (String) Last-seen value of the deployment.environment.name (or deployment.environment) OpenTelemetry resource attribute, e.g. production, staging... Computed.
- `runtime_name` (String) Last-seen value of the process.runtime.name OpenTelemetry resource attribute... Computed.
- `runtime_version` (String) Last-seen value of the process.runtime.version OpenTelemetry resource attribute... Computed.
- `cloud_provider` (String) Last-seen value of the cloud.provider OpenTelemetry resource attribute, e.g. aws, gcp, azure... Computed.
- `cloud_platform` (String) Last-seen value of the cloud.platform OpenTelemetry resource attribute, e.g. aws_ec2, gcp_compute_engine... Computed.
- `cloud_region` (String) Last-seen value of the cloud.region OpenTelemetry resource attribute, e.g. us-east-1... Computed.
- `cloud_account_id` (String) Last-seen value of the cloud.account.id OpenTelemetry resource attribute... Computed.
