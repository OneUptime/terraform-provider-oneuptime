---
page_title: "oneuptime_trace_pipeline_processor Data Source - oneuptime"
subcategory: "Other"
description: |-
  Individual processors within a trace pipeline that transform span data during ingestion.
---

# oneuptime_trace_pipeline_processor (Data Source)

Individual processors within a trace pipeline that transform span data during ingestion.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one trace pipeline processor may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_trace_pipeline_processor" "example" {
  name = "Example trace pipeline processor"
}

# Or by id:
data "oneuptime_trace_pipeline_processor" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the user who created this trace pipeline processor. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this processor is active.
- `name` (String) Friendly name for this processor.
- `processor_type` (String) The type of processor: AttributeRemapper, SpanNameRemapper, StatusRemapper, SpanKindRemapper, or CategoryProcessor.
- `sort_order` (Number) Where this processor runs within its pipeline, lowest number first. A new processor is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `trace_pipeline_id` (String) ID of the trace pipeline this processor belongs to. The ID of a `oneuptime_trace_pipeline`.

### Read-Only

- `configuration` (String) Processor-specific configuration as JSON (e.g., source/target fields, mapping rules). A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this trace pipeline processor belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
