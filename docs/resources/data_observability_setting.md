---
page_title: "oneuptime_data_observability_setting Resource - oneuptime"
subcategory: "Other"
description: |-
  How Data Observability treats a project's databases: how much a check may read, whether checks wait when a database is busy, and whether failing rows are kept.
---

# oneuptime_data_observability_setting (Resource)

How Data Observability treats a project's databases: how much a check may read, whether checks wait when a database is busy, and whether failing rows are kept.

## Example Usage

```terraform
resource "oneuptime_data_observability_setting" "example" {

}
```

## Schema

### Optional

- `keep_failing_rows` (Boolean) A failing check keeps up to 20 of the rows that broke it, with sensitive columns masked. They are stored in OneUptime. Turning this off deletes them. Defaults to `false`.
- `max_gigabytes_scanned_per_check` (Number) BigQuery and Athena bill by bytes read. A check that would read more is stopped, and says so. Defaults to `10`.
- `wait_when_database_is_busy` (Boolean) Checks wait for the next collection when a database is close to its connection limit or has many statements running. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project these settings belong to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing data observability setting by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_data_observability_setting.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_data_observability_setting.example <id>
```
