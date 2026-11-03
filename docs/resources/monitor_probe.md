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
- `last_ping_at` (String) When the probe last checked the monitor. Set by OneUptime.
- `next_ping_at` (String) When the probe checks the monitor next. Set by OneUptime: right away for a new one, then on the monitor's interval.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_monitoring_log` (String) A JSON value: read it with `jsondecode()`.
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
