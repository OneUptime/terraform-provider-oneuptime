---
page_title: "oneuptime_io_t_fleet Resource - oneuptime"
subcategory: "Other"
description: |-
  IoT device fleets that are being monitored in this project. Each fleet is auto-discovered when an IoT device or gateway sends metrics with the iot.fleet.name OTel resource attribute, or can be manually registered.
---

# oneuptime_io_t_fleet (Resource)

IoT device fleets that are being monitored in this project. Each fleet is auto-discovered when an IoT device or gateway sends metrics with the iot.fleet.name OTel resource attribute, or can be manually registered.

## Example Usage

```terraform
resource "oneuptime_io_t_fleet" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this IoT fleet. This is the join key — it must match the iot.fleet.name OTel resource attribute stamped by the IoT device or gateway...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description for this IoT fleet..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_archived` (Bool) Is this IoT fleet archived? Archived IoT fleets are hidden from lists but keep collecting telemetry...
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this IoT fleet. Leave blank to use the project-wide default...
- `telemetry_retention_config` (String) Per-pillar retention overrides for this IoT fleet (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the IoT fleet default, then the project's retention settings...
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected)..
- `agent_version` (String) Version of the IoT agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute..
- `last_seen_at` (String) A date time object..
- `device_count` (Number) Cached count of devices in this fleet..
- `online_device_count` (Number) Cached count of devices currently online (iot_device_up == 1) in this fleet. Rendered as 'Devices X/Y online' next to deviceCount...

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

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_io_t_fleet.example <id>
```
