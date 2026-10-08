---
page_title: "oneuptime_form_submission Data Source - oneuptime"
subcategory: "Other"
description: |-
  The submissions made through this project's forms: every answer, the name and email the submitter gave, and the incident or scheduled maintenance event each one created.
---

# oneuptime_form_submission (Data Source)

The submissions made through this project's forms: every answer, the name and email the submitter gave, and the incident or scheduled maintenance event each one created.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one form submission may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_form_submission" "example" {
  form_id = oneuptime_form.example.id
}

# Or by id:
data "oneuptime_form_submission" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `form_id` (String) ID of the form this submission was made through. The ID of a `oneuptime_form`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this submission created. Empty for a form that schedules maintenance, and once that incident is deleted. The ID of a `oneuptime_incident`.
- `scheduled_maintenance_id` (String) ID of the scheduled maintenance event this submission created. Empty for a form that creates incidents, and once that event is deleted. The ID of a `oneuptime_scheduled_maintenance_event`.
- `submitter_name` (String) The name the submitter gave, as they typed it. Empty when the form does not ask for it, or lets people leave it out and they did.
- `target_type` (String) What the submission created: Incident, or ScheduledMaintenance (a scheduled maintenance event).

### Read-Only

- `answers` (String) Every question the submitter answered, in the form's order: the question's id (fieldId), its label when the form was submitted, the stored value, and the value as a person reads it (displayValue). A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `submitter_email` (String) The email address the submitter gave. It is not verified, and nothing is sent to it. Empty when the form does not ask for it, or lets people leave it out and they did.
- `updated_at` (String) Date and Time when the object was updated.
