---
page_title: "oneuptime_table_view Data Source - oneuptime"
subcategory: "Other"
description: |-
  Table View is view settings for a table in a project. It contains columns, filters, and other settings.
---

# oneuptime_table_view (Data Source)

Table View is view settings for a table in a project. It contains columns, filters, and other settings. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_table_view" "by_name" {
  name = "example-table_view"
}

data "oneuptime_table_view" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `table_id` (String) ID of the table this view is for.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `query` (String) Filters for this table view.. Computed.
- `sort` (String) Sort for this table view.. Computed.
- `items_on_page` (Number) Items on page.. Computed.
- `facets` (String) Facet selections (owner, labels, status, etc.) for this table view.. Computed.
- `columns` (String) Which columns are shown, and in what order, for this table view.. Computed.
