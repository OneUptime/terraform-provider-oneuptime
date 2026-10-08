---
page_title: "oneuptime_log_saved_view Data Source - oneuptime"
subcategory: "Other"
description: |-
  Save and reuse log explorer views, including the current filters, columns, sorting, and page size.
---

# oneuptime_log_saved_view (Data Source)

Save and reuse log explorer views, including the current filters, columns, sorting, and page size.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one log saved view may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_log_saved_view" "example" {
  name = "Example log saved view"
}

# Or by id:
data "oneuptime_log_saved_view" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this saved log view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_default` (Boolean) Whether this saved log view should be applied by default.
- `name` (String) Friendly name for this saved log view.
- `page_size` (Number) Number of logs per page for this saved view.
- `sort_field` (String) Active sort field for this saved log view.
- `sort_order` (String) Sort order for this saved log view.

### Read-Only

- `columns` (String) Selected log table columns for this saved view. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this saved log view belongs to. The ID of a `oneuptime_project`.
- `query` (String) Serialized log query for this saved view. A JSON value: write it with `jsonencode()`.
- `time_range` (String) Time selection for this saved view — the rolling range token (e.g. Past 1 Hour), or an absolute window when the range is Custom. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
