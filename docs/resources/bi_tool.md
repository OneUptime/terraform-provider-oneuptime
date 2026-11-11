---
page_title: "oneuptime_bi_tool Resource - oneuptime"
subcategory: "Other"
description: |-
  A BI tool whose dashboards, and the tables behind them, Data Observability reads, so a table's alerts list the dashboards built on it.
---

# oneuptime_bi_tool (Resource)

A BI tool whose dashboards, and the tables behind them, Data Observability reads, so a table's alerts list the dashboards built on it.

## Example Usage

```terraform
resource "oneuptime_bi_tool" "example" {
  name = "Example bi tool"
  tool = "Metabase"
}
```

## Schema

### Required

- `name` (String) What the BI tool is called here.
- `tool` (String) Which BI tool it is. Allowed values: `Metabase`, `Superset`, `Tableau`, `Looker`, `Power BI`, `Hightouch`, `Census`.

### Optional

- `alert_on_failed_refreshes` (Boolean) Whether OneUptime keeps a monitor that opens an alert for each dataset whose refresh failed (Power BI), resolved once a later refresh completes. Defaults to `false`.
- `base_url` (String) The tool's own address, such as https://metabase.acme.com or https://acme.cloud.looker.com. Power BI needs none.
- `credentials` (String, Sensitive) JSON credentials that may read the tool's dashboards and what they query. Encrypted at rest and never returned by the API.
- `is_enabled` (Boolean) Whether its dashboards are read every hour. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `dashboard_count` (Number) How many dashboards the last read found.
- `failed_refreshes` (String) The datasets whose last refresh failed or was cancelled (Power BI), each with its error, until a later refresh completes - and for an hour after that, marked recovered. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `last_sync_error` (String) Why the last read failed; empty when it did not.
- `last_sync_warning` (String) What the last read could not read while it read the dashboards - views or refreshes - and why; empty when it read everything.
- `last_synced_at` (String) When its dashboards were last read.
- `matched_table_count` (Number) How many watched tables the last read found behind its dashboards.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `refresh_monitor_id` (String) ID of the Metrics monitor that alerts on the tool's failed refreshes.
- `updated_at` (String) Date and Time when the object was updated.
- `usage_read_from` (String) The start of the hours of the tool's activity log read so far (Power BI): its dashboards' views are counted from them.
- `usage_read_until` (String) The end of the hours of the tool's activity log read so far (Power BI): the next read goes on from here.

## Import

Import an existing bi tool by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_bi_tool.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_bi_tool.example <id>
```
