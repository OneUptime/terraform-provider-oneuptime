---
page_title: "oneuptime_data_check Resource - oneuptime"
subcategory: "Other"
description: |-
  A check on the values of a table - empty or repeated values, values outside a list or pattern, a query of your own - run on a schedule and written as a metric.
---

# oneuptime_data_check (Resource)

A check on the values of a table - empty or repeated values, values outside a list or pattern, a query of your own - run on a schedule and written as a metric.

## Example Usage

```terraform
resource "oneuptime_data_check" "example" {
  data_asset_id = oneuptime_data_asset.example.id
  name          = "Example data check"
  check_type    = "Null Percent"
  description   = "Managed by Terraform"
}
```

## Schema

### Required

- `check_type` (String) What the check measures. Allowed values: `Null Percent`, `Duplicate Percent`, `Distinct Count`, `Minimum`, `Maximum`, `Average`, `Sum`, `Invalid Percent`, `Format Mismatch Percent`, `Orphan Count`, `Row Count`, `Custom SQL`, `Failing Rows`, `dbt Test`.
- `data_asset_id` (String) ID of the table the check reads. The ID of a `oneuptime_data_asset`.
- `name` (String) What the check is called, as its metric and alerts name it.

### Optional

- `accepted_values` (String) For Invalid Percent: the values the column may hold. A JSON value: write it with `jsonencode()`.
- `alert_condition` (String) When the check's value should alert: unusually high or low for the hour of the week (learned), or above or below a threshold. Allowed values: `Unusually High`, `Unusually Low`, `Above`, `Below`.
- `alert_threshold` (Number) The threshold an Above or Below alert compares with.
- `column_name` (String) The column the check reads, for a check on one column.
- `custom_sql` (String) For Custom SQL: one read-only SELECT returning a number - or, segmented, a segment and a number per row.
- `description` (String) Why the check matters, and what to do when it fails.
- `format_preset` (String) For Format Mismatch Percent: a format everyone wants, instead of writing the pattern.
- `interval_in_minutes` (Number) How often the check runs: 15, 60, 360 or 1440 minutes. Defaults to `60`.
- `is_enabled` (Boolean) Whether the check runs. Defaults to `true`.
- `pattern` (String) For Format Mismatch Percent: the regular expression every value must match.
- `reference_column_name` (String) For Orphan Count: the column of that table it refers to.
- `reference_table_key` (String) For Orphan Count: the table the column refers to, as schema.table.
- `sample_percent` (Number) Read this percent of the table's rows instead of all of them, for a check that reads a rate.
- `segment_column` (String) Run the check once per value of this column - a tenant, a region - one series each, the largest 100.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_reconciliation_id` (String) For a check that reads one side of a reconciliation: the reconciliation's ID. Such a check is changed only from its reconciliation.
- `data_store_id` (String) ID of the data store the table is in, set from the table.
- `external_id` (String) For a dbt test, the test's unique id in the dbt project, which each upload of its artifacts is matched by.
- `failing_rows_sample` (String) A masked sample of the rows the last run found failing, when the project keeps them. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `is_managed_by_file` (Boolean) Whether the check came from a file in a repository (oneuptime data-checks apply).
- `last_run_at` (String) When the check last ran.
- `last_run_message` (String) Why the last run failed, when it did.
- `last_run_status` (String) How the last run went: Ran, or Failed.
- `last_segment_count` (Number) How many segments the last run read.
- `last_value` (Number) The value the last run read - the largest segment's, for a segmented check.
- `monitor_id` (String) ID of the Metrics monitor that alerts on the check. The ID of a `oneuptime_monitor`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing data check by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_check.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_data_check.example <id>
```
