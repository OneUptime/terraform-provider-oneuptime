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
  trace_pipeline_id = "123e4567-e89b-12d3-a456-426614174000"
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  processor_type = "Example short text"
}
```

## Schema

### Required

- `trace_pipeline_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Name object.
- `processor_type` (String) The type of processor: AttributeRemapper, SpanNameRemapper, StatusRemapper, SpanKindRemapper, or CategoryProcessor...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `configuration` (String) Processor-specific configuration as JSON (e.g., source/target fields, mapping rules)...
- `is_enabled` (Bool) Whether this processor is active...
- `sort_order` (Number) Determines the execution order of this processor within its pipeline...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_trace_pipeline_processor.example <id>
```
