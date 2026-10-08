---
page_title: "oneuptime_trace_pipeline_processor Resource - oneuptime"
subcategory: "Other"
description: |-
  Individual processors within a trace pipeline that transform span data during ingestion.
---

# oneuptime_trace_pipeline_processor (Resource)

Individual processors within a trace pipeline that transform span data during ingestion.

## Example Usage

```terraform
resource "oneuptime_trace_pipeline_processor" "example" {
  trace_pipeline_id = oneuptime_trace_pipeline.example.id
  name              = "Example trace pipeline processor"
  processor_type    = "Example short text"
}
```

## Schema

### Required

- `name` (String) Friendly name for this processor.
- `processor_type` (String) The type of processor: AttributeRemapper, SpanNameRemapper, StatusRemapper, SpanKindRemapper, or CategoryProcessor.
- `trace_pipeline_id` (String) ID of the trace pipeline this processor belongs to. The ID of a `oneuptime_trace_pipeline`.

### Optional

- `configuration` (String) Processor-specific configuration as JSON (e.g., source/target fields, mapping rules). A JSON value: write it with `jsonencode()`.
- `is_enabled` (Boolean) Whether this processor is active. Defaults to `true`.
- `sort_order` (Number) Where this processor runs within its pipeline, lowest number first. A new processor is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this trace pipeline processor. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this trace pipeline processor belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing trace pipeline processor by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_trace_pipeline_processor.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_trace_pipeline_processor.example <id>
```
