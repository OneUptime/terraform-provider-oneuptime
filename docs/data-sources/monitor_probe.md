---
page_title: "oneuptime_monitor_probe Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Add probes to monitor your resource from multiple locations around the world.
---

# oneuptime_monitor_probe (Data Source)

Add probes to monitor your resource from multiple locations around the world.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor probe may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_probe" "example" {
  probe_id = oneuptime_probe.example.id
}

# Or by id:
data "oneuptime_monitor_probe" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) When on, the probe checks the monitor.
- `monitor_id` (String) ID of your OneUptime Monitor in which this object belongs. The ID of a `oneuptime_monitor`.
- `probe_id` (String) ID of your OneUptime Probe in which this object belongs. The ID of a `oneuptime_probe`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_monitoring_log` (String) A JSON value: read it with `jsondecode()`.
- `last_ping_at` (String) When the probe last checked the monitor. Set by OneUptime.
- `next_ping_at` (String) When the probe checks the monitor next. Set by OneUptime: right away for a new one, then on the monitor's interval.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
