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
  runbook_id            = oneuptime_runbook.example.id
  runbook_name_snapshot = "Example short text"
}
```

## Schema

### Required

- `runbook_id` (String) ID of the Runbook this execution belongs to. The ID of a `oneuptime_runbook`.
- `runbook_name_snapshot` (String) Name of the runbook at the time this execution was created (preserved even if the runbook is later renamed or deleted).

### Optional

- `alert_id` (String) ID of the alert that triggered or hosts this runbook execution. The ID of a `oneuptime_alert`.
- `completed_at` (String) Time at which this runbook execution completed.
- `failure_reason` (String) Reason this runbook execution failed (if it did).
- `incident_id` (String) ID of the incident that triggered or hosts this runbook execution. The ID of a `oneuptime_incident`.
- `scheduled_maintenance_id` (String) ID of the scheduled maintenance event that triggered this runbook execution. The ID of a `oneuptime_scheduled_maintenance_event`.
- `started_at` (String) Time at which this runbook execution started.
- `step_executions` (String) Per-step execution state. Each entry mirrors a step from the runbook with status, output, and timestamps. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `id` (String) Unique identifier for the resource.
- `inputs` (String) The values this run was started with: the JSON body of the webhook call that started it. Steps read them as data - Bash as ONEUPTIME_INPUT_<NAME> environment variables, JavaScript as args.inputs - never as script text. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status` (String) Current status of this runbook execution.
- `trigger_details` (String) More about what started this run: for a scheduled run, the time it was due (scheduledFor) and whether it made up for times OneUptime missed (isCatchUp) or waited for the run before it (waitedForPreviousRun); for a webhook call, the address it came from (callerIp). A JSON value: write it with `jsonencode()`.
- `trigger_type` (String) What started this run: "Manual" (a person, or the API), "Rule" (a runbook rule), "AutoRemediation", "AI" (OneUptime AI), "Schedule" (the runbook's schedule) or "Webhook" (a call to the runbook's webhook URL). Empty for runs started before it was recorded.
- `triggered_by_user_id` (String) ID of the User who triggered this runbook execution. The ID of a `oneuptime_user` (see the data source).
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing runbook execution by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_runbook_execution.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_runbook_execution.example <id>
```
