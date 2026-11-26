---
page_title: "oneuptime_iot_fleet Resource - oneuptime"
subcategory: "IoT"
description: |-
  IoT device fleets that are being monitored in this project. Each fleet is auto-discovered when an IoT device or gateway sends metrics with the iot.fleet.name OTel resource attribute, or can be manually registered.
---

# oneuptime_iot_fleet (Resource)

IoT device fleets that are being monitored in this project. Each fleet is auto-discovered when an IoT device or gateway sends metrics with the iot.fleet.name OTel resource attribute, or can be manually registered.

~> **Renamed:** this resource was called `oneuptime_io_t_fleet` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing iot fleet:

```terraform
moved {
  from = oneuptime_io_t_fleet.example
  to   = oneuptime_iot_fleet.example
}
```

## Example Usage

```terraform
resource "oneuptime_iot_fleet" "example" {
  name        = "Example iot fleet"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this IoT fleet. This is the join key — it must match the iot.fleet.name OTel resource attribute stamped by the IoT device or gateway.

### Optional

- `description` (String) Friendly description for this IoT fleet.
- `is_archived` (Boolean) Is this IoT fleet archived? Archived IoT fleets are hidden from lists but keep collecting telemetry. Defaults to `false`.
- `labels` (Set of String) Labels to group and filter by. OneUptime adds to it too, on its own: the entries listed here are added and taken off, and every other entry stays. Left out, it keeps the entries OneUptime gives it. IDs of `oneuptime_label` resources.
- `retain_telemetry_data_for_days` (Number) Number of days to retain telemetry data for this IoT fleet. Leave blank to use the project-wide default.
- `telemetry_retention_config` (String) Per-pillar retention overrides for this IoT fleet (logs by severity, traces by status, metrics, profiles). Unset fields fall back to the IoT fleet default, then the project's retention settings. A JSON value: write it with `jsonencode()`.

### Read-Only

- `agent_version` (String) Version of the IoT agent reporting telemetry, as self-reported via the oneuptime.agent.version resource attribute.
- `archived_at` (String) When was this IoT fleet archived?
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `device_count` (Number) Cached count of devices in this fleet.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `last_seen_at` (String) When metrics were last received from this fleet.
- `online_device_count` (Number) Cached count of devices currently online (iot_device_up == 1) in this fleet. Rendered as 'Devices X/Y online' next to deviceCount.
- `otel_collector_status` (String) Connection status of the OTel Collector agent (connected or disconnected).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing iot fleet by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_iot_fleet.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_iot_fleet.example <id>
```
