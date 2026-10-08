---
page_title: "oneuptime_log_saved_view Resource - oneuptime"
subcategory: "Other"
description: |-
  Save and reuse log explorer views, including the current filters, columns, sorting, and page size.
---

# oneuptime_log_saved_view (Resource)

Save and reuse log explorer views, including the current filters, columns, sorting, and page size.

## Example Usage

```terraform
resource "oneuptime_log_saved_view" "example" {
  name = "Example log saved view"
}
```

## Schema

### Required

- `name` (String) Friendly name for this saved log view.

### Optional

- `columns` (String) Selected log table columns for this saved view. A JSON value: write it with `jsonencode()`.
- `is_default` (Boolean) Whether this saved log view should be applied by default. Defaults to `false`.
- `page_size` (Number) Number of logs per page for this saved view. Defaults to `100`.
- `query` (String) Serialized log query for this saved view. A JSON value: write it with `jsonencode()`.
- `sort_field` (String) Active sort field for this saved log view.
- `sort_order` (String) Sort order for this saved log view.
- `time_range` (String) Time selection for this saved view — the rolling range token (e.g. Past 1 Hour), or an absolute window when the range is Custom. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this saved log view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this saved log view belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing log saved view by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_log_saved_view.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_log_saved_view.example <id>
```
