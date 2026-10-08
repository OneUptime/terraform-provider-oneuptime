---
page_title: "oneuptime_table_view Data Source - oneuptime"
subcategory: "Other"
description: |-
  Table View is view settings for a table in a project. It contains columns, filters, and other settings.
---

# oneuptime_table_view (Data Source)

Table View is view settings for a table in a project. It contains columns, filters, and other settings.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one table view may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_table_view" "example" {
  name = "Example table view"
}

# Or by id:
data "oneuptime_table_view" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `items_on_page` (Number) Items on page.
- `name` (String) Any friendly name of this object.
- `table_id` (String) ID of the table this view is for.

### Read-Only

- `columns` (String) Which columns are shown, and in what order, for this table view. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `facets` (String) Facet selections (owner, labels, status, etc.) for this table view. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `query` (String) Filters for this table view. A JSON value: write it with `jsonencode()`.
- `sort` (String) Sort for this table view. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
