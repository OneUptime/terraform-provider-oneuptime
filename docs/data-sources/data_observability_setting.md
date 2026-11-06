---
page_title: "oneuptime_data_observability_setting Data Source - oneuptime"
subcategory: "Other"
description: |-
  How Data Observability treats a project's databases: how much a check may read, whether checks wait when a database is busy, and whether failing rows are kept.
---

# oneuptime_data_observability_setting (Data Source)

How Data Observability treats a project's databases: how much a check may read, whether checks wait when a database is busy, and whether failing rows are kept.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data observability setting may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_observability_setting" "example" {
  max_gigabytes_scanned_per_check = 1
}

# Or by id:
data "oneuptime_data_observability_setting" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `keep_failing_rows` (Boolean) A failing check keeps up to 20 of the rows that broke it, with sensitive columns masked. They are stored in OneUptime. Turning this off deletes them.
- `max_gigabytes_scanned_per_check` (Number) BigQuery and Athena bill by bytes read. A check that would read more is stopped, and says so.
- `wait_when_database_is_busy` (Boolean) Checks wait for the next collection when a database is close to its connection limit or has many statements running.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project these settings belong to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
