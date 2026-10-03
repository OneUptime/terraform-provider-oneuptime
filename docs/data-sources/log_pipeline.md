---
page_title: "oneuptime_log_pipeline Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure server-side log processing pipelines that transform logs at ingest time.
---

# oneuptime_log_pipeline (Data Source)

Configure server-side log processing pipelines that transform logs at ingest time. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_log_pipeline" "by_name" {
  name = "example-log_pipeline"
}

data "oneuptime_log_pipeline" "by_id" {
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
- `description` (String) Description of what this log pipeline does... Computed.
- `filter_query` (String) Filter expression that determines which logs this pipeline applies to... Computed.
- `is_enabled` (Bool) Whether this log pipeline is active... Computed.
- `sort_order` (Number) Where this pipeline runs among the project's log pipelines, lowest number first. A new pipeline is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
