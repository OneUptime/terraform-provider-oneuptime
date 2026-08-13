---
page_title: "oneuptime_monitor_log Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  API endpoints for Monitor Log
---

# oneuptime_monitor_log (Data Source)

API endpoints for Monitor Log Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_monitor_log" "by_name" {
  name = "example-monitor_log"
}

data "oneuptime_monitor_log" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `monitor_id` (String) Monitor ID. Computed.
- `time` (String) Time. Computed.
- `log_body` (String) Log Body. Computed.
