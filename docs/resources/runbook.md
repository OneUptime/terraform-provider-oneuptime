---
page_title: "oneuptime_runbook Resource - oneuptime"
subcategory: "Other"
description: |-
  Reusable response procedures (manual checklists or scripts) that can be attached to incidents, alerts, or scheduled maintenance.
---

# oneuptime_runbook (Resource)

Reusable response procedures (manual checklists or scripts) that can be attached to incidents, alerts, or scheduled maintenance.

## Example Usage

```terraform
resource "oneuptime_runbook" "example" {
  name        = "Example runbook"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_enabled` (Boolean) Is this runbook enabled? Defaults to `true`.
- `is_schedule_enabled` (Boolean) Whether the runbook runs on its schedule (Schedule Cron, read in Schedule Time Zone). Turning it on takes permission to run the runbook as well as to edit it. Defaults to `false`.
- `is_webhook_enabled` (Boolean) Whether a POST request to the runbook's webhook URL starts a run. The URL holds a secret key OneUptime generates when this is first turned on, shown on the runbook's Triggers page to people who may edit and run the runbook. Turning it on takes permission to run the runbook as well as to edit it. Defaults to `false`.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `schedule_cron` (String) When the runbook runs on its schedule: a five-field cron expression (minute, hour, day of month, month, day of week), read on the wall clock of Schedule Time Zone. Runbooks run at most once a minute.
- `schedule_missed_run_policy` (String) What happens to scheduled times OneUptime could not start a run at, because it was down: "RunOnce" (the default) starts one run when it is back, "Skip" starts none.
- `schedule_overlap_policy` (String) What a scheduled time does while a run of the runbook is still queued, running or waiting for someone: "Skip" (the default) starts nothing for that time, "RunAfter" starts one run as soon as the current run ends.
- `schedule_timezone` (String) The IANA time zone the schedule is read in, such as Europe/London or America/New_York. A 09:00 schedule runs at 09:00 there all year, also after the clocks change for daylight saving time. Left out, it is UTC.
- `steps` (String) Ordered list of steps to run for this runbook. Each step is one of Manual, JavaScript, HTTP request, Bash or AI. A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `schedule_next_run_at` (String) When the runbook's schedule starts its next run. Worked out by OneUptime; empty while the schedule is off.
- `schedule_queued_run_at` (String) The scheduled time a run is waiting for the current run to end for, when "If a Run Is Still Going" is "RunAfter". Set and cleared by OneUptime.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing runbook by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_runbook.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_runbook.example <id>
```
