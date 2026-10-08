---
page_title: "oneuptime_storage_array Resource - oneuptime"
subcategory: "Other"
description: |-
  Storage arrays (Pure Storage FlashArray and FlashBlade) that are being monitored in this project. Each array is auto-discovered when the OneUptime Storage Array Agent sends metrics, or can be registered by hand.
---

# oneuptime_storage_array (Resource)

Storage arrays (Pure Storage FlashArray and FlashBlade) that are being monitored in this project. Each array is auto-discovered when the OneUptime Storage Array Agent sends metrics, or can be registered by hand.

## Example Usage

```terraform
resource "oneuptime_storage_array" "example" {
  name        = "Example storage array"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this storage array in OneUptime. This is the join key — it must match the storage.array.name OTel resource attribute stamped by the OneUptime Storage Array Agent.

### Optional

- `agent_version` (String) Version of the OneUptime Storage Array Agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `bucket_count` (Number) Cached count of object store buckets on the array (FlashBlade).
- `capacity_bytes` (Number) Cached usable capacity of the array in bytes (purefa_array_space_bytes{space="capacity"} / purefb_array_space_bytes{type="array",space="capacity"}). Null until the first array metric batch arrives.
- `capacity_used_percent` (Number) Cached array space utilization in percent (purefa_array_space_utilization / purefb_array_space_utilization). Stored as decimal so sub-percent precision survives the round trip. Null until the first array metric batch arrives.
- `critical_alert_count` (Number) Cached count of open array alerts with critical severity.
- `data_reduction_ratio` (Number) Cached data reduction ratio of the array (purefa_array_space_data_reduction_ratio / purefb_array_space_data_reduction_ratio), for example 4.2 for 4.2:1. Null until the first array metric batch arrives.
- `description` (String) Friendly description for this storage array.
- `file_system_count` (Number) Cached count of file systems on the array (FlashBlade).
- `hardware_component_count` (Number) Cached count of hardware components the array reports (chassis, controllers, drive bays, power supplies, fans, ports, blades...).
- `health_status` (Number) Cached array health derived from the array's open alerts and hardware state: 0 = OK, 1 = Warning (a warning alert, or a degraded or unknown component), 2 = Critical (a critical alert, or a failed or critical component). Rendered as the health pill. Null until the first metric batch arrives.
- `host_count` (Number) Cached count of hosts defined on the array (FlashArray).
- `is_archived` (Boolean) Is this storage array archived? Archived storage arrays are hidden from lists but keep collecting telemetry. Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When metrics were last received from this storage array.
- `open_alert_count` (Number) Cached count of alerts open on the array itself (purefa_alerts_open / purefb_alerts_open series, hidden alerts excluded).
- `os_name` (String) Name of the array's operating system as it reports it (os on purefa_info / purefb_info).
- `os_version` (String) Version of the array's operating system as it reports it (version on purefa_info / purefb_info).
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `pod_count` (Number) Cached count of pods (ActiveCluster / ActiveDR replication containers) on the array (FlashArray).
- `reported_name` (String) The name the array reports for itself (array_name on purefa_info / purefb_info). Can differ from the OneUptime name, which comes from the agent configuration.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this storage array. Leave blank to use the project-wide default.
- `storage_system` (String) The storage platform this array runs, normalized from the storage.system OTel resource attribute (or detected from the metric names): purestorage.flasharray or purestorage.flashblade.
- `system_id` (String) The array's own system identifier (system_id on purefa_info / purefb_info).
- `telemetry_retention_config` (String) Per-pillar retention overrides for this storage array (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the storage array default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `unhealthy_hardware_count` (Number) Cached count of hardware components, drives and controllers in a critical, degraded, failed or unknown state.
- `used_bytes` (Number) Cached physical space used on the array in bytes: capacity minus empty space, or capacity times utilization when the empty series is missing. Null until the first array metric batch arrives.
- `volume_count` (Number) Cached count of volumes on the array (FlashArray).
- `warning_alert_count` (Number) Cached count of open array alerts with warning severity.

### Read-Only

- `archived_at` (String) When was this storage array archived?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing storage array by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_storage_array.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_storage_array.example <id>
```
