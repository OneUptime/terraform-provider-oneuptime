---
page_title: "oneuptime_ai_run_event Data Source - oneuptime"
subcategory: "Other"
description: |-
  An event in an AI run: LLM calls, tool calls with validated arguments, and lifecycle transitions.
---

# oneuptime_ai_run_event (Data Source)

An event in an AI run: LLM calls, tool calls with validated arguments, and lifecycle transitions.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai run event may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_ai_run_event" "example" {
  ai_run_id = data.oneuptime_ai_run.example.id
}

# Or by id:
data "oneuptime_ai_run_event" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_run_id` (String) ID of the run this event belongs to. The ID of a `oneuptime_ai_run` (see the data source).
- `citation_id` (String) ID of the citation this event minted (e.g. C1), if it produced one.
- `event_type` (String) Type of event.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `sequence` (Number) Order of this event within the run.
- `tool_name` (String) Name of the tool for tool-call events.
- `user_id` (String) ID of the user whose run this event belongs to. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this event belongs to. The ID of a `oneuptime_project`.
- `result_summary` (String) Summary of the result: row count, duration, truncation and bytes sent to the LLM. A JSON value: write it with `jsonencode()`.
- `tool_arguments` (String) Validated tool arguments as executed. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
