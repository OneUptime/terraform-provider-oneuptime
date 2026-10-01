---
page_title: "oneuptime_incident_form_submission Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  The submissions made through this project's incident forms: the form, the incident each one declared, and the name and email the reporter gave.
---

# oneuptime_incident_form_submission (Data Source)

The submissions made through this project's incident forms: the form, the incident each one declared, and the name and email the reporter gave. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_form_submission" "by_name" {
  name = "example-incident_form_submission"
}

data "oneuptime_incident_form_submission" "by_id" {
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
- `incident_form_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `reporter_name` (String) The name the reporter gave, as they typed it. Empty when the form lets people report anonymously and they did... Computed.
- `reporter_email` (String) Email object. Computed.
