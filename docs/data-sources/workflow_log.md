---
page_title: "oneuptime_workflow_log Data Source - oneuptime"
subcategory: "Workflows"
description: |-
  Logs of the workflows executed
---

# oneuptime_workflow_log (Data Source)

Logs of the workflows executed

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one workflow log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_workflow_log" "example" {
  workflow_id = oneuptime_workflow.example.id
}

# Or by id:
data "oneuptime_workflow_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `logs` (String) Logs.
- `workflow_id` (String) ID of Workflow this logs belong to. The ID of a `oneuptime_workflow`.
- `workflow_status` (String) Status of this workflow.

### Read-Only

- `completed_at` (String) When did this workflow complete.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `resume_at` (String) When this workflow run is scheduled to resume after a Sleep step (only set while the run is Waiting).
- `started_at` (String) When did this workflow start.
- `step_trace` (String) Structured per-step record of this run: arguments, return values, the port taken and timing. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
