---
page_title: "oneuptime_exception_saved_view Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Save and reuse views of the exceptions list, including the current search, filters, time range, status, sort and page size.
---

# oneuptime_exception_saved_view (Data Source)

Save and reuse views of the exceptions list, including the current search, filters, time range, status, sort and page size.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one exception saved view may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_exception_saved_view" "example" {
  name = "Example exception saved view"
}

# Or by id:
data "oneuptime_exception_saved_view" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this saved exceptions view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_default` (Boolean) Whether this saved exceptions view should be applied by default.
- `name` (String) Friendly name for this saved exceptions view.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this saved exceptions view belongs to. The ID of a `oneuptime_project`.
- `query` (String) Serialized exceptions list view state (search, filters, time range, page size, status, lens, assignee filter and sort) for this saved view. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
