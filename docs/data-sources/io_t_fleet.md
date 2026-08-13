---
page_title: "oneuptime_io_t_fleet Data Source - oneuptime"
subcategory: "Other"
description: |-
  IoT device fleets that are being monitored in this project. Each fleet is auto-discovered when an IoT device or gateway sends metrics with the iot.fleet.name OTel resource attribute, or can be manually registered.
---

# oneuptime_io_t_fleet (Data Source)

IoT device fleets that are being monitored in this project. Each fleet is auto-discovered when an IoT device or gateway sends metrics with the iot.fleet.name OTel resource attribute, or can be manually registered. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_io_t_fleet" "by_name" {
  name = "example-io_t_fleet"
}

data "oneuptime_io_t_fleet" "by_id" {
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
- `description` (String) Friendly description for this IoT fleet.. Computed.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).. Computed.
- `agent_version` (String) Version of the IoT agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `device_count` (Number) Cached count of devices in this fleet.. Computed.
- `online_device_count` (Number) Cached count of devices currently online (iot_device_up == 1) in this fleet. Rendered as 'Devices X/Y online' next to deviceCount... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this IoT fleet archived? Archived IoT fleets are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this IoT fleet. Leave blank to use the project-wide default... Computed.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this IoT fleet (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the IoT fleet default, then the project's retention settings... Computed.
