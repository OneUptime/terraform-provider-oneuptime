---
page_title: "oneuptime_workflow_log Data Source - oneuptime"
subcategory: "Workflows"
description: |-
  Logs of the workflows executed
---

# oneuptime_workflow_log (Data Source)

Logs of the workflows executed Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_workflow_log" "by_name" {
  name = "example-workflow_log"
}

data "oneuptime_workflow_log" "by_id" {
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
- `workflow_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `logs` (String) Logs.. Computed.
- `workflow_status` (String) Status of this workflow.. Computed.
- `started_at` (String) A date time object.. Computed.
- `completed_at` (String) A date time object.. Computed.
- `resume_at` (String) A date time object.. Computed.
- `step_trace` (String) Structured per-step record of this run: arguments, return values, the port taken and timing... Computed.
