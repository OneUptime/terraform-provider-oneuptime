---
page_title: "oneuptime_monitor_log Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  API endpoints for Monitor Log
---

# oneuptime_monitor_log (Data Source)

API endpoints for Monitor Log

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_log" "example" {
  monitor_id = "example-monitor-id"
}

# Or by id:
data "oneuptime_monitor_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `monitor_id` (String) Monitor ID.
- `time` (String) Time.

### Read-Only

- `log_body` (String) Log Body. A JSON value: write it with `jsonencode()`.
- `project_id` (String) Project ID.
