---
page_title: "oneuptime_data_check Data Source - oneuptime"
subcategory: "Other"
description: |-
  A check on the values of a table - empty or repeated values, values outside a list or pattern, a query of your own - run on a schedule and written as a metric.
---

# oneuptime_data_check (Data Source)

A check on the values of a table - empty or repeated values, values outside a list or pattern, a query of your own - run on a schedule and written as a metric.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data check may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_check" "example" {
  name = "Example data check"
}

# Or by id:
data "oneuptime_data_check" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_condition` (String) When the check's value should alert: unusually high or low for the hour of the week (learned), or above or below a threshold.
- `alert_threshold` (Number) The threshold an Above or Below alert compares with.
- `check_type` (String) What the check measures.
- `column_name` (String) The column the check reads, for a check on one column.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `custom_sql` (String) For Custom SQL: one read-only SELECT returning a number - or, segmented, a segment and a number per row.
- `data_asset_id` (String) ID of the table the check reads. The ID of a `oneuptime_data_asset`.
- `data_reconciliation_id` (String) For a check that reads one side of a reconciliation: the reconciliation's ID. Such a check is changed only from its reconciliation.
- `data_store_id` (String) ID of the data store the table is in, set from the table.
- `description` (String) Why the check matters, and what to do when it fails.
- `external_id` (String) For a dbt test, the test's unique id in the dbt project, which each upload of its artifacts is matched by.
- `format_preset` (String) For Format Mismatch Percent: a format everyone wants, instead of writing the pattern.
- `health` (String) How the check stands now: Failing while its alert is open, Could Not Run when its query failed, Passing, No Alert, Waiting before its first run, or Off.
- `health_rank` (Number) The check's health as a number to sort by, worst first: 0 failing, 1 could not run, 2 passing, 3 no alert, 4 waiting, 5 off.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `interval_in_minutes` (Number) How often the check runs: 15, 60, 360 or 1440 minutes.
- `is_enabled` (Boolean) Whether the check runs.
- `is_managed_by_file` (Boolean) Whether the check came from a file in a repository (oneuptime data-checks apply).
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `last_run_message` (String) Why the last run failed, when it did.
- `last_run_status` (String) How the last run went: Ran, or Failed.
- `last_segment_count` (Number) How many segments the last run read.
- `last_value` (Number) The value the last run read - the largest segment's, for a segmented check.
- `monitor_id` (String) ID of the Metrics monitor that alerts on the check. The ID of a `oneuptime_monitor`.
- `name` (String) What the check is called, as its metric and alerts name it.
- `pattern` (String) For Format Mismatch Percent: the regular expression every value must match.
- `reference_column_name` (String) For Orphan Count: the column of that table it refers to.
- `reference_table_key` (String) For Orphan Count: the table the column refers to, as schema.table.
- `sample_percent` (Number) Read this percent of the table's rows instead of all of them, for a check that reads a rate.
- `segment_column` (String) Run the check once per value of this column - a tenant, a region - one series each, the largest 100.

### Read-Only

- `accepted_values` (String) For Invalid Percent: the values the column may hold. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `failing_rows_sample` (String) A masked sample of the rows the last run found failing, when the project keeps them. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_run_at` (String) When the check last ran.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
