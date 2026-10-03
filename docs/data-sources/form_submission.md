---
page_title: "oneuptime_form_submission Data Source - oneuptime"
subcategory: "Other"
description: |-
  The submissions made through this project's forms: every answer, the name and email the submitter gave, and the incident or scheduled maintenance event each one created.
---

# oneuptime_form_submission (Data Source)

The submissions made through this project's forms: every answer, the name and email the submitter gave, and the incident or scheduled maintenance event each one created. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_form_submission" "by_name" {
  name = "example-form_submission"
}

data "oneuptime_form_submission" "by_id" {
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
- `form_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `answers` (String) Every question the submitter answered, in the form's order: the question's id (fieldId), its label when the form was submitted, the stored value, and the value as a person reads it (displayValue)... Computed.
- `submitter_name` (String) The name the submitter gave, as they typed it. Empty when the form does not ask for it, or lets people leave it out and they did... Computed.
- `submitter_email` (String) Email object. Computed.
- `target_type` (String) What the submission created: Incident, or ScheduledMaintenance (a scheduled maintenance event)... Computed.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `scheduled_maintenance_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
