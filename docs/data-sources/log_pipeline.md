---
page_title: "oneuptime_log_pipeline Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure server-side log processing pipelines that transform logs at ingest time.
---

# oneuptime_log_pipeline (Data Source)

Configure server-side log processing pipelines that transform logs at ingest time.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one log pipeline may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_log_pipeline" "example" {
  name = "Example log pipeline"
}

# Or by id:
data "oneuptime_log_pipeline" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this log pipeline. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of what this log pipeline does.
- `filter_query` (String) Filter expression that determines which logs this pipeline applies to.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this log pipeline is active.
- `name` (String) Friendly name for this log pipeline.
- `sort_order` (Number) Where this pipeline runs among the project's log pipelines, lowest number first. A new pipeline is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this log pipeline belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
