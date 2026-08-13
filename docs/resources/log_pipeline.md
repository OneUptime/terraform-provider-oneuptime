---
page_title: "oneuptime_log_pipeline Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure server-side log processing pipelines that transform logs at ingest time.
---

# oneuptime_log_pipeline (Resource)

Configure server-side log processing pipelines that transform logs at ingest time.

## Example Usage

```terraform
resource "oneuptime_log_pipeline" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of what this log pipeline does...
- `filter_query` (String) Filter expression that determines which logs this pipeline applies to...
- `is_enabled` (Bool) Whether this log pipeline is active...
- `sort_order` (Number) Determines the execution order of this pipeline relative to others...

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
terraform import oneuptime_log_pipeline.example <id>
```
