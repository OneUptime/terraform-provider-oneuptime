---
page_title: "oneuptime_incident_form Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Forms anyone with the link can fill in to report an incident, without a OneUptime account. Each submission declares an incident in this project.
---

# oneuptime_incident_form (Data Source)

Forms anyone with the link can fill in to report an incident, without a OneUptime account. Each submission declares an incident in this project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_form" "by_name" {
  name = "example-incident_form"
}

data "oneuptime_incident_form" "by_id" {
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
- `description` (String) Shown at the top of the form's public page, above the questions: what the form is for and what happens after it is sent. Markdown... Computed.
- `is_enabled` (Bool) Whether the form's link works. While the form is turned off, its public page shows a not-available message and nothing can be submitted... Computed.
- `share_key` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `allow_reporter_to_choose_severity` (Bool) When on, the form asks the reporter to choose a severity from the project's incident severities, with the form's own severity chosen to begin with... Computed.
- `incident_template_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description_setting` (String) Whether the form asks the reporter to describe the incident: Required, Optional or Hidden... Computed.
- `custom_field_settings` (String) The incident custom fields the form asks for, keyed by each field's template variable key (variableKey). Required means the reporter must answer it, Optional that they may leave it empty. Only the fields listed as Required or Optional are asked: a field that is not listed, or is Hidden or Default, is not on the form. The answers become the incident's custom field values... Computed.
- `is_reporter_details_required` (Bool) When on, the reporter must give their name and email. When off, they may report anonymously... Computed.
- `success_message` (String) Shown to the reporter after they submit the form, together with the new incident's number. Markdown... Computed.
- `ip_whitelist` (String) The networks the form can be opened and submitted from: one IPv4 or IPv6 address, or one IPv4 range in CIDR notation (such as 10.0.0.0/8), per line. IPv6 ranges are not supported. Leave it empty to allow any network... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
