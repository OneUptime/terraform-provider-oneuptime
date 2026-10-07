---
page_title: "oneuptime_storage_array Data Source - oneuptime"
subcategory: "Other"
description: |-
  Storage arrays (Pure Storage FlashArray and FlashBlade) that are being monitored in this project. Each array is auto-discovered when the OneUptime Storage Array Agent sends metrics, or can be registered by hand.
---

# oneuptime_storage_array (Data Source)

Storage arrays (Pure Storage FlashArray and FlashBlade) that are being monitored in this project. Each array is auto-discovered when the OneUptime Storage Array Agent sends metrics, or can be registered by hand. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_storage_array" "by_name" {
  name = "example-storage_array"
}

data "oneuptime_storage_array" "by_id" {
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
- `description` (String) Friendly description for this storage array.. Computed.
- `storage_system` (String) The storage platform this array runs, normalized from the storage.system OTel resource attribute (or detected from the metric names): purestorage.flasharray or purestorage.flashblade... Computed.
- `reported_name` (String) The name the array reports for itself (array_name on purefa_info / purefb_info). Can differ from the OneUptime name, which comes from the agent configuration... Computed.
- `system_id` (String) The array's own system identifier (system_id on purefa_info / purefb_info)... Computed.
- `os_name` (String) Name of the array's operating system as it reports it (os on purefa_info / purefb_info)... Computed.
- `os_version` (String) Version of the array's operating system as it reports it (version on purefa_info / purefb_info)... Computed.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).. Computed.
- `agent_version` (String) Version of the OneUptime Storage Array Agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `capacity_bytes` (Number) Cached usable capacity of the array in bytes (purefa_array_space_bytes{space="capacity"} / purefb_array_space_bytes{type="array",space="capacity"}). Null until the first array metric batch arrives... Computed.
- `used_bytes` (Number) Cached physical space used on the array in bytes: capacity minus empty space, or capacity times utilization when the empty series is missing. Null until the first array metric batch arrives... Computed.
- `capacity_used_percent` (Number) Cached array space utilization in percent (purefa_array_space_utilization / purefb_array_space_utilization). Stored as decimal so sub-percent precision survives the round trip. Null until the first array metric batch arrives... Computed.
- `data_reduction_ratio` (Number) Cached data reduction ratio of the array (purefa_array_space_data_reduction_ratio / purefb_array_space_data_reduction_ratio), for example 4.2 for 4.2:1. Null until the first array metric batch arrives... Computed.
- `open_alert_count` (Number) Cached count of alerts open on the array itself (purefa_alerts_open / purefb_alerts_open series, hidden alerts excluded)... Computed.
- `critical_alert_count` (Number) Cached count of open array alerts with critical severity... Computed.
- `warning_alert_count` (Number) Cached count of open array alerts with warning severity... Computed.
- `volume_count` (Number) Cached count of volumes on the array (FlashArray)... Computed.
- `host_count` (Number) Cached count of hosts defined on the array (FlashArray)... Computed.
- `pod_count` (Number) Cached count of pods (ActiveCluster / ActiveDR replication containers) on the array (FlashArray)... Computed.
- `file_system_count` (Number) Cached count of file systems on the array (FlashBlade)... Computed.
- `bucket_count` (Number) Cached count of object store buckets on the array (FlashBlade)... Computed.
- `hardware_component_count` (Number) Cached count of hardware components the array reports (chassis, controllers, drive bays, power supplies, fans, ports, blades...)... Computed.
- `unhealthy_hardware_count` (Number) Cached count of hardware components, drives and controllers in a critical, degraded, failed or unknown state... Computed.
- `health_status` (Number) Cached array health derived from the array's open alerts and hardware state: 0 = OK, 1 = Warning (a warning alert, or a degraded or unknown component), 2 = Critical (a critical alert, or a failed or critical component). Rendered as the health pill. Null until the first metric batch arrives... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this storage array archived? Archived storage arrays are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this storage array. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this storage array (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the storage array default, then the project's retention settings... Computed.
