---
page_title: "oneuptime_monitor_probe Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Add probes to monitor your resource from multiple locations around the world.
---

# oneuptime_monitor_probe (Data Source)

Add probes to monitor your resource from multiple locations around the world. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_monitor_probe" "by_name" {
  name = "example-monitor_probe"
}

data "oneuptime_monitor_probe" "by_id" {
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
- `probe_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `last_ping_at` (String) A date time object.. Computed.
- `next_ping_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_enabled` (Bool) Monitor probe is_enabled. Computed.
- `last_monitoring_log` (String) Monitor probe last_monitoring_log. Computed.
