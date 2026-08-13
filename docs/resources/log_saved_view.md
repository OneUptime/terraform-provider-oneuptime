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
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
}
```

## Schema

### Required

- `name` (String) Name object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `query` (String) Serialized log query for this saved view...
- `columns` (String) Selected log table columns for this saved view...
- `sort_field` (String) Active sort field for this saved log view...
- `sort_order` (String) Sort order for this saved log view...
- `page_size` (Number) Number of logs per page for this saved view...
- `time_range` (String) Time selection for this saved view — the rolling range token (e.g. Past 1 Hour), or an absolute window when the range is Custom...
- `is_default` (Bool) Whether this saved log view should be applied by default...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_log_saved_view.example <id>
```
