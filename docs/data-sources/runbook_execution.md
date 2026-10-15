---
page_title: "oneuptime_runbook_execution Data Source - oneuptime"
subcategory: "Other"
description: |-
  A single run of a Runbook.
---

# oneuptime_runbook_execution (Data Source)

A single run of a Runbook.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one runbook execution may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_runbook_execution" "example" {
  runbook_id = oneuptime_runbook.example.id
}

# Or by id:
data "oneuptime_runbook_execution" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of the alert that triggered or hosts this runbook execution. The ID of a `oneuptime_alert`.
- `failure_reason` (String) Reason this runbook execution failed (if it did).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident that triggered or hosts this runbook execution. The ID of a `oneuptime_incident`.
- `runbook_id` (String) ID of the Runbook this execution belongs to. The ID of a `oneuptime_runbook`.
- `runbook_name_snapshot` (String) Name of the runbook at the time this execution was created (preserved even if the runbook is later renamed or deleted).
- `scheduled_maintenance_id` (String) ID of the scheduled maintenance event that triggered this runbook execution. The ID of a `oneuptime_scheduled_maintenance_event`.
- `status` (String) Current status of this runbook execution.
- `trigger_type` (String) What started this run: "Manual" (a person, or the API), "Rule" (a runbook rule), "AutoRemediation", "AI" (OneUptime AI), "Schedule" (the runbook's schedule) or "Webhook" (a call to the runbook's webhook URL). Empty for runs started before it was recorded.
- `triggered_by_user_id` (String) ID of the User who triggered this runbook execution. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `completed_at` (String) Time at which this runbook execution completed.
- `created_at` (String) Date and Time when the object was created.
- `inputs` (String) The values this run was started with: the JSON body of the webhook call that started it. Steps read them as data - Bash as ONEUPTIME_INPUT_<NAME> environment variables, JavaScript as args.inputs - never as script text. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `started_at` (String) Time at which this runbook execution started.
- `step_executions` (String) Per-step execution state. Each entry mirrors a step from the runbook with status, output, and timestamps. A JSON value: write it with `jsonencode()`.
- `trigger_details` (String) More about what started this run: for a scheduled run, the time it was due (scheduledFor) and whether it made up for times OneUptime missed (isCatchUp) or waited for the run before it (waitedForPreviousRun); for a webhook call, the address it came from (callerIp). A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
