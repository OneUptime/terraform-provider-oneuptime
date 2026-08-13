---
page_title: "oneuptime_metric_saved_view Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Save and reuse metrics explorer views, including the current search, filters, time range, and page size.
---

# oneuptime_metric_saved_view (Data Source)

Save and reuse metrics explorer views, including the current search, filters, time range, and page size. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_metric_saved_view" "by_name" {
  name = "example-metric_saved_view"
}

data "oneuptime_metric_saved_view" "by_id" {
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
- `query` (String) Serialized metrics explorer view state (search, filters, time range, page size) for this saved view... Computed.
- `is_default` (Bool) Whether this saved metric view should be applied by default... Computed.
- `view_type` (String) Which surface this saved view belongs to ('list' or 'explorer'). Null means 'list' — rows created before this column existed all came from the metric list page... Computed.
