---
page_title: "oneuptime_iot_fleet Data Source - oneuptime"
subcategory: "Other"
description: |-
  IoT device fleets that are being monitored in this project. Each fleet is auto-discovered when an IoT device or gateway sends metrics with the iot.fleet.name OTel resource attribute, or can be manually registered.
---

# oneuptime_iot_fleet (Data Source)

IoT device fleets that are being monitored in this project. Each fleet is auto-discovered when an IoT device or gateway sends metrics with the iot.fleet.name OTel resource attribute, or can be manually registered.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one iot fleet may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_io_t_fleet` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_iot_fleet" "example" {
  name = "Example iot fleet"
}

# Or by id:
data "oneuptime_iot_fleet" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `agent_version` (String) Version of the IoT agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description for this IoT fleet.
- `device_count` (Number) Cached count of devices in this fleet.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_archived` (Boolean) Is this IoT fleet archived? Archived IoT fleets are hidden from lists but keep collecting telemetry.
- `name` (String) Name of this IoT fleet. This is the join key — it must match the iot.fleet.name OTel resource attribute stamped by the IoT device or gateway.
- `online_device_count` (Number) Cached count of devices currently online (iot_device_up == 1) in this fleet. Rendered as 'Devices X/Y online' next to deviceCount.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this IoT fleet. Leave blank to use the project-wide default.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `archived_at` (String) When was this IoT fleet archived?
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_seen_at` (String) When metrics were last received from this fleet.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this IoT fleet (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the IoT fleet default, then the project's retention settings. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
