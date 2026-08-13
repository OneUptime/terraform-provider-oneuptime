---
page_title: "oneuptime_log_saved_view Data Source - oneuptime"
subcategory: "Other"
description: |-
  Save and reuse log explorer views, including the current filters, columns, sorting, and page size.
---

# oneuptime_log_saved_view (Data Source)

Save and reuse log explorer views, including the current filters, columns, sorting, and page size. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_log_saved_view" "by_name" {
  name = "example-log_saved_view"
}

data "oneuptime_log_saved_view" "by_id" {
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
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `query` (String) Serialized log query for this saved view... Computed.
- `columns` (String) Selected log table columns for this saved view... Computed.
- `sort_field` (String) Active sort field for this saved log view... Computed.
- `sort_order` (String) Sort order for this saved log view... Computed.
- `page_size` (Number) Number of logs per page for this saved view... Computed.
- `time_range` (String) Time selection for this saved view — the rolling range token (e.g. Past 1 Hour), or an absolute window when the range is Custom... Computed.
- `is_default` (Bool) Whether this saved log view should be applied by default... Computed.
