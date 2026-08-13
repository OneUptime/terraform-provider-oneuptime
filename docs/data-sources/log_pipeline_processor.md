---
page_title: "oneuptime_log_pipeline_processor Data Source - oneuptime"
subcategory: "Other"
description: |-
  Individual processors within a log pipeline that transform log data during ingestion.
---

# oneuptime_log_pipeline_processor (Data Source)

Individual processors within a log pipeline that transform log data during ingestion. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_log_pipeline_processor" "by_name" {
  name = "example-log_pipeline_processor"
}

data "oneuptime_log_pipeline_processor" "by_id" {
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
- `log_pipeline_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `processor_type` (String) The type of processor: GrokParser, AttributeRemapper, SeverityRemapper, or CategoryProcessor... Computed.
- `configuration` (String) Processor-specific configuration as JSON (e.g., grok pattern, source/target fields, mapping rules)... Computed.
- `is_enabled` (Bool) Whether this processor is active... Computed.
- `sort_order` (Number) Determines the execution order of this processor within its pipeline... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
