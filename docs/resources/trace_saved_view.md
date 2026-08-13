---
page_title: "oneuptime_trace_saved_view Resource - oneuptime"
subcategory: "Other"
description: |-
  Save and reuse traces explorer views, including the current search, filters, time range, and page size.
---

# oneuptime_trace_saved_view (Resource)

Save and reuse traces explorer views, including the current search, filters, time range, and page size.

## Example Usage

```terraform
resource "oneuptime_trace_saved_view" "example" {
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
- `query` (String) Serialized traces explorer view state (search, filters, time range, page size) for this saved view...
- `is_default` (Bool) Whether this saved trace view should be applied by default...

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
terraform import oneuptime_trace_saved_view.example <id>
```
