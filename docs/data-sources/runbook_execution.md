---
page_title: "oneuptime_runbook_execution Data Source - oneuptime"
subcategory: "Other"
description: |-
  A single run of a Runbook.
---

# oneuptime_runbook_execution (Data Source)

A single run of a Runbook. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_runbook_execution" "by_name" {
  name = "example-runbook_execution"
}

data "oneuptime_runbook_execution" "by_id" {
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
- `runbook_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `runbook_name_snapshot` (String) Name of the runbook at the time this execution was created (preserved even if the runbook is later renamed or deleted)... Computed.
- `status` (String) Current status of this runbook execution... Computed.
- `step_executions` (String) Per-step execution state. Each entry mirrors a step from the runbook with status, output, and timestamps... Computed.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `scheduled_maintenance_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `triggered_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `started_at` (String) A date time object.. Computed.
- `completed_at` (String) A date time object.. Computed.
- `failure_reason` (String) Reason this runbook execution failed (if it did)... Computed.
