---
page_title: "oneuptime_monitor_probe Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Add probes to monitor your resource from multiple locations around the world.
---

# oneuptime_monitor_probe (Resource)

Add probes to monitor your resource from multiple locations around the world.

## Example Usage

```terraform
resource "oneuptime_monitor_probe" "example" {
  probe_id   = oneuptime_probe.example.id
  monitor_id = oneuptime_monitor.example.id
}
```

## Schema

### Required

- `monitor_id` (String) ID of your OneUptime Monitor in which this object belongs. The ID of a `oneuptime_monitor`.
- `probe_id` (String) ID of your OneUptime Probe in which this object belongs. The ID of a `oneuptime_probe`.

### Optional

- `is_enabled` (Boolean) When on, the probe checks the monitor. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `last_monitoring_log` (String) A JSON value: read it with `jsondecode()`.
- `last_ping_at` (String) When the probe last checked the monitor. Set by OneUptime.
- `next_ping_at` (String) When the probe checks the monitor next. Set by OneUptime: right away for a new one, then on the monitor's interval.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing monitor probe by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_monitor_probe.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_monitor_probe.example <id>
```
