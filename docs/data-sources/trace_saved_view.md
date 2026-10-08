---
page_title: "oneuptime_trace_saved_view Data Source - oneuptime"
subcategory: "Other"
description: |-
  Save and reuse traces explorer views, including the current search, filters, time range, and page size.
---

# oneuptime_trace_saved_view (Data Source)

Save and reuse traces explorer views, including the current search, filters, time range, and page size.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one trace saved view may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_trace_saved_view" "example" {
  name = "Example trace saved view"
}

# Or by id:
data "oneuptime_trace_saved_view" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this saved trace view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_default` (Boolean) Whether this saved trace view should be applied by default.
- `name` (String) Friendly name for this saved trace view.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this saved trace view belongs to. The ID of a `oneuptime_project`.
- `query` (String) Serialized traces explorer view state (search, filters, time range, page size) for this saved view. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
