---
page_title: "oneuptime_data_reconciliation Data Source - oneuptime"
subcategory: "Other"
description: |-
  One measure of two tables - rows, a sum, distinct values, a query's number - compared, in one database or two, with a tolerance, and an alert when they drift.
---

# oneuptime_data_reconciliation (Data Source)

One measure of two tables - rows, a sum, distinct values, a query's number - compared, in one database or two, with a tolerance, and an alert when they drift.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data reconciliation may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_reconciliation" "example" {
  name = "Example data reconciliation"
}

# Or by id:
data "oneuptime_data_reconciliation" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Why the two must agree, and what to do when they do not.
- `difference` (Number) The target's value less the source's, at the latest comparison.
- `difference_percent` (Number) The difference as a percentage of the source's value, at the latest comparison.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `interval_in_minutes` (Number) How often both sides are read and compared: every 15 minutes, hour, 6 hours or day.
- `monitor_id` (String) ID of the Metrics monitor that alerts when the two drift. The ID of a `oneuptime_monitor`.
- `name` (String) What the reconciliation is called, as its alert names it.
- `source_column_name` (String) For Sum or Distinct Count: the column of the source table.
- `source_custom_sql` (String) For Custom SQL: a read-only query of the source's database returning one number, as value.
- `source_data_asset_id` (String) ID of the source table the reconciliation reads. The ID of a `oneuptime_data_asset`.
- `source_data_check_id` (String) ID of the data check that reads the source, made and kept by the reconciliation. The ID of a `oneuptime_data_check`.
- `source_measure` (String) What is measured of the source table: Row Count, Sum, Distinct Count, or Custom SQL.
- `source_value` (Number) The source's value at the latest comparison.
- `status` (String) In Step, Drifted, Waiting for both sides to be read together, or Failed when a side's check could not run.
- `status_message` (String) Why the reconciliation is waiting or failed.
- `target_column_name` (String) For Sum or Distinct Count: the column of the target table.
- `target_custom_sql` (String) For Custom SQL: a read-only query of the target's database returning one number, as value.
- `target_data_asset_id` (String) ID of the target table the reconciliation reads. The ID of a `oneuptime_data_asset`.
- `target_data_check_id` (String) ID of the data check that reads the target, made and kept by the reconciliation. The ID of a `oneuptime_data_check`.
- `target_measure` (String) What is measured of the target table: Row Count, Sum, Distinct Count, or Custom SQL.
- `target_value` (Number) The target's value at the latest comparison.
- `tolerance` (Number) How far apart the two may be before they have drifted: 0 means they must agree exactly.
- `tolerance_type` (String) How the tolerance is read: Absolute, an amount, or Percent, of the source's value.

### Read-Only

- `compared_at` (String) When the later of the two compared values was read.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_changed_at` (String) When the status last changed.
- `updated_at` (String) Date and Time when the object was updated.
