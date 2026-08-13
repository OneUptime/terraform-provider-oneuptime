---
page_title: "oneuptime_metric_saved_view Resource - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Save and reuse metrics explorer views, including the current search, filters, time range, and page size.
---

# oneuptime_metric_saved_view (Resource)

Save and reuse metrics explorer views, including the current search, filters, time range, and page size.

## Example Usage

```terraform
resource "oneuptime_metric_saved_view" "example" {
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
- `query` (String) Serialized metrics explorer view state (search, filters, time range, page size) for this saved view...
- `is_default` (Bool) Whether this saved metric view should be applied by default...
- `view_type` (String) Which surface this saved view belongs to ('list' or 'explorer'). Null means 'list' — rows created before this column existed all came from the metric list page...

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
terraform import oneuptime_metric_saved_view.example <id>
```
