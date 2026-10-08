---
page_title: "oneuptime_metric_saved_view Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Save and reuse metrics explorer views, including the current search, filters, time range, and page size.
---

# oneuptime_metric_saved_view (Data Source)

Save and reuse metrics explorer views, including the current search, filters, time range, and page size.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one metric saved view may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_metric_saved_view" "example" {
  name = "Example metric saved view"
}

# Or by id:
data "oneuptime_metric_saved_view" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this saved metric view. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_default` (Boolean) Whether this saved metric view should be applied by default.
- `name` (String) Friendly name for this saved metric view.
- `view_type` (String) Which surface this saved view belongs to ('list' or 'explorer'). Null means 'list' — rows created before this column existed all came from the metric list page.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this saved metric view belongs to. The ID of a `oneuptime_project`.
- `query` (String) Serialized metrics explorer view state (search, filters, time range, page size) for this saved view. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
