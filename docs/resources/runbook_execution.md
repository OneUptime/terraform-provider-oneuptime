---
page_title: "oneuptime_runbook_execution Resource - oneuptime"
subcategory: "Other"
description: |-
  A single run of a Runbook.
---

# oneuptime_runbook_execution (Resource)

A single run of a Runbook.

## Example Usage

```terraform
resource "oneuptime_runbook_execution" "example" {
  runbook_id = "123e4567-e89b-12d3-a456-426614174000"
  runbook_name_snapshot = "Example short text"
}
```

## Schema

### Required

- `runbook_id` (String) A unique identifier for an object, represented as a UUID..
- `runbook_name_snapshot` (String) Name of the runbook at the time this execution was created (preserved even if the runbook is later renamed or deleted)...

### Optional

- `step_executions` (String) Per-step execution state. Each entry mirrors a step from the runbook with status, output, and timestamps...
- `incident_id` (String) A unique identifier for an object, represented as a UUID..
- `alert_id` (String) A unique identifier for an object, represented as a UUID..
- `scheduled_maintenance_id` (String) A unique identifier for an object, represented as a UUID..
- `triggered_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `status` (String) Current status of this runbook execution...
- `started_at` (String) A date time object..
- `completed_at` (String) A date time object..
- `failure_reason` (String) Reason this runbook execution failed (if it did)...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `project_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_runbook_execution.example <id>
```
