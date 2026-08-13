---
page_title: "oneuptime_trace_drop_filter Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure rules to drop or sample spans before storage to reduce volume and cost.
---

# oneuptime_trace_drop_filter (Resource)

Configure rules to drop or sample spans before storage to reduce volume and cost.

## Example Usage

```terraform
resource "oneuptime_trace_drop_filter" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  filter_query = "This is an example of longer text content that might be stored in this field."
  action = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `filter_query` (String) Filter expression that identifies which spans to drop or sample...
- `action` (String) What to do with matching spans: 'drop' to discard entirely, 'sample' to keep a percentage...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of what this drop filter does...
- `sample_percentage` (Number) When action is 'sample', the percentage of matching spans to keep (1-99)...
- `is_enabled` (Bool) Whether this drop filter is active...
- `sort_order` (Number) Determines the evaluation order of this filter relative to others...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `dropped_count` (Number) Total number of spans this filter has discarded since it was created...
- `last_dropped_at` (String) A date time object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_trace_drop_filter.example <id>
```
