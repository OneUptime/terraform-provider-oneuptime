---
page_title: "oneuptime_bi_dashboard Data Source - oneuptime"
subcategory: "Other"
description: |-
  A dashboard of a connected BI tool, with the watched tables its charts query.
---

# oneuptime_bi_dashboard (Data Source)

A dashboard of a connected BI tool, with the watched tables its charts query.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one bi dashboard may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_bi_dashboard" "example" {
  name = "Example bi dashboard"
}

# Or by id:
data "oneuptime_bi_dashboard" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_bi_connection_id` (String) ID of the BI tool the dashboard was read from. The ID of a `oneuptime_bi_tool`.
- `external_id` (String) The dashboard's own ID in its tool.
- `folder` (String) The collection, project, folder or workspace it is in.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_refresh_status` (String) How the last refresh of the data it charts went (Power BI): Completed, Failed, Cancelled, Disabled, In Progress or Unknown.
- `name` (String) The dashboard's name, as its tool shows it.
- `owner` (String) Who owns it in its tool.
- `table_count` (Number) How many watched tables its charts query.
- `tool` (String) Which BI tool the dashboard is in.
- `url` (String) Where the dashboard opens in its tool.
- `view_count` (Number) How many times it was opened in the last 28 days, from its tool's activity log (Power BI).
- `viewer_count` (Number) How many people opened it in the last 28 days.

### Read-Only

- `changed_at` (String) When the dashboard last changed in its tool.
- `created_at` (String) Date and Time when the object was created.
- `data_assets` (Set of String) The watched tables the dashboard's charts query. IDs of `oneuptime_data_asset` resources.
- `last_refresh_at` (String) When the last refresh of the data it charts ended, or started while it runs.
- `last_viewed_at` (String) When someone last opened it.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `refresh_history` (String) The latest refreshes of the data it charts, each with how it went and a failed one's error, and how many ran and failed in the last seven days. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
- `usage_summary` (String) Who opened it most in the last 28 days, with how often and when last, and since when its views are counted while the 28 days are still being read. A JSON value: read it with `jsondecode()`.
