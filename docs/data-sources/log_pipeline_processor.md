---
page_title: "oneuptime_log_pipeline_processor Data Source - oneuptime"
subcategory: "Other"
description: |-
  Individual processors within a log pipeline that transform log data during ingestion.
---

# oneuptime_log_pipeline_processor (Data Source)

Individual processors within a log pipeline that transform log data during ingestion.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one log pipeline processor may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_log_pipeline_processor" "example" {
  name = "Example log pipeline processor"
}

# Or by id:
data "oneuptime_log_pipeline_processor" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this log pipeline processor. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this processor is active.
- `log_pipeline_id` (String) ID of the log pipeline this processor belongs to. The ID of a `oneuptime_log_pipeline`.
- `name` (String) Friendly name for this processor.
- `processor_type` (String) The type of processor: GrokParser, KeyValueParser, AttributeRemapper, SeverityRemapper, or CategoryProcessor.
- `sort_order` (Number) Where this processor runs within its pipeline, lowest number first. A new processor is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `configuration` (String) Processor-specific configuration as JSON (e.g., grok pattern, source/target fields, mapping rules). A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this log pipeline processor belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
