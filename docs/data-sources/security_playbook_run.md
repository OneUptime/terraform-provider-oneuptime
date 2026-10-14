---
page_title: "oneuptime_security_playbook_run Data Source - oneuptime"
subcategory: "Other"
description: |-
  One run of a security playbook on a case: the steps it took, what each did, and who approved the steps that asked.
---

# oneuptime_security_playbook_run (Data Source)

One run of a security playbook on a case: the steps it took, what each did, and who approved the steps that asked.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security playbook run may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_playbook_run" "example" {
  security_case_id = oneuptime_security_case.example.id
}

# Or by id:
data "oneuptime_security_playbook_run" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who started the run. The ID of a `oneuptime_user` (see the data source).
- `current_step_index` (Number) The step the run is on, counting from 0.
- `error_message` (String) Why the run failed, when it did.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `security_case_id` (String) ID of the case the playbook runs on. The ID of a `oneuptime_security_case`.
- `security_playbook_id` (String) ID of the playbook being run. The ID of a `oneuptime_security_playbook`.
- `status` (String) Queued, Running, Waiting for Approval, Succeeded, Failed or Cancelled. Set it to Cancelled to stop a run that has not finished.
- `threat_intel_matched` (Boolean) Whether a Check Threat Intel step of this run found a known-bad entity. Steps set to run only on a match read it.
- `trigger` (String) What started the run: Case Opened, Detection Added, or Manual for a person.
- `triggered_by_user_id` (String) ID of the person who ran the playbook. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `completed_at` (String) When the run finished, failed or was cancelled.
- `created_at` (String) Date and Time when the object was created.
- `playbook_name` (String) The playbook's name when the run started, kept if the playbook is renamed or deleted.
- `project_id` (String) ID of the project this run belongs to. The ID of a `oneuptime_project`.
- `started_at` (String) When the run started its first step.
- `step_results` (String) What each step did: its status, when it started and finished, what it reported, and who decided on its approval. A JSON value: write it with `jsonencode()`.
- `steps` (String) The playbook's steps as they were when the run started. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
