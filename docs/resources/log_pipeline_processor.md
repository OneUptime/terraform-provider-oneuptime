---
page_title: "oneuptime_log_pipeline_processor Resource - oneuptime"
subcategory: "Other"
description: |-
  Individual processors within a log pipeline that transform log data during ingestion.
---

# oneuptime_log_pipeline_processor (Resource)

Individual processors within a log pipeline that transform log data during ingestion.

## Example Usage

```terraform
resource "oneuptime_log_pipeline_processor" "example" {
  log_pipeline_id = "123e4567-e89b-12d3-a456-426614174000"
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  processor_type = "Example short text"
}
```

## Schema

### Required

- `log_pipeline_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Name object.
- `processor_type` (String) The type of processor: GrokParser, AttributeRemapper, SeverityRemapper, or CategoryProcessor...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `configuration` (String) Processor-specific configuration as JSON (e.g., grok pattern, source/target fields, mapping rules)...
- `is_enabled` (Bool) Whether this processor is active...
- `sort_order` (Number) Where this processor runs within its pipeline, lowest number first. A new processor is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them...

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
terraform import oneuptime_log_pipeline_processor.example <id>
```
