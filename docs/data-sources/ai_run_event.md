---
page_title: "oneuptime_ai_run_event Data Source - oneuptime"
subcategory: "Other"
description: |-
  An event in an AI run: LLM calls, tool calls with validated arguments, and lifecycle transitions.
---

# oneuptime_ai_run_event (Data Source)

An event in an AI run: LLM calls, tool calls with validated arguments, and lifecycle transitions. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_ai_run_event" "by_name" {
  name = "example-ai_run_event"
}

data "oneuptime_ai_run_event" "by_id" {
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
- `ai_run_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `sequence` (Number) Order of this event within the run... Computed.
- `event_type` (String) Type of event... Computed.
- `tool_name` (String) Name of the tool for tool-call events... Computed.
- `tool_arguments` (String) Validated tool arguments as executed... Computed.
- `result_summary` (String) Summary of the result: row count, duration, truncation and bytes sent to the LLM... Computed.
- `citation_id` (String) ID of the citation this event minted (e.g. C1), if it produced one... Computed.
