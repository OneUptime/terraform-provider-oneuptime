---
page_title: "oneuptime_incident_postmortem_template Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage postmortem templates for your incidents
---

# oneuptime_incident_postmortem_template (Data Source)

Manage postmortem templates for your incidents Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_postmortem_template" "by_name" {
  name = "example-incident_postmortem_template"
}

data "oneuptime_incident_postmortem_template" "by_id" {
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
- `postmortem_note` (String) Markdown template used when documenting an incident postmortem... Computed.
- `template_name` (String) Name of the Postmortem Template.. Computed.
- `template_description` (String) Description of the Postmortem Template.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
