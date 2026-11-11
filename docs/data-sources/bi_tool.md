---
page_title: "oneuptime_bi_tool Data Source - oneuptime"
subcategory: "Other"
description: |-
  A BI tool whose dashboards, and the tables behind them, Data Observability reads, so a table's alerts list the dashboards built on it.
---

# oneuptime_bi_tool (Data Source)

A BI tool whose dashboards, and the tables behind them, Data Observability reads, so a table's alerts list the dashboards built on it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one bi tool may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_bi_tool" "example" {
  name = "Example bi tool"
}

# Or by id:
data "oneuptime_bi_tool" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_on_failed_refreshes` (Boolean) Whether OneUptime keeps a monitor that opens an alert for each dataset whose refresh failed (Power BI), resolved once a later refresh completes.
- `base_url` (String) The tool's own address, such as https://metabase.acme.com or https://acme.cloud.looker.com. Power BI needs none.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `dashboard_count` (Number) How many dashboards the last read found.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether its dashboards are read every hour.
- `last_sync_error` (String) Why the last read failed; empty when it did not.
- `last_sync_warning` (String) What the last read could not read while it read the dashboards - views or refreshes - and why; empty when it read everything.
- `matched_table_count` (Number) How many watched tables the last read found behind its dashboards.
- `name` (String) What the BI tool is called here.
- `refresh_monitor_id` (String) ID of the Metrics monitor that alerts on the tool's failed refreshes.
- `tool` (String) Which BI tool it is.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `failed_refreshes` (String) The datasets whose last refresh failed or was cancelled (Power BI), each with its error, until a later refresh completes - and for an hour after that, marked recovered. A JSON value: read it with `jsondecode()`.
- `last_synced_at` (String) When its dashboards were last read.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `usage_read_from` (String) The start of the hours of the tool's activity log read so far (Power BI): its dashboards' views are counted from them.
- `usage_read_until` (String) The end of the hours of the tool's activity log read so far (Power BI): the next read goes on from here.
