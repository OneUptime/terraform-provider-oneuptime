---
page_title: "oneuptime_incident_form Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Forms anyone with the link can fill in to report an incident, without a OneUptime account. Each submission declares an incident in this project.
---

# oneuptime_incident_form (Resource)

Forms anyone with the link can fill in to report an incident, without a OneUptime account. Each submission declares an incident in this project.

## Example Usage

```terraform
resource "oneuptime_incident_form" "example" {
  name = "Example short text"
  incident_severity_id = "123e4567-e89b-12d3-a456-426614174000"
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `name` (String) The form's name, shown as the heading of its public page. Unique within the project...
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Shown at the top of the form's public page, above the questions: what the form is for and what happens after it is sent. Markdown...
- `is_enabled` (Bool) Whether the form's link works. While the form is turned off, its public page shows a not-available message and nothing can be submitted...
- `allow_reporter_to_choose_severity` (Bool) When on, the form asks the reporter to choose a severity from the project's incident severities, with the form's own severity chosen to begin with...
- `incident_template_id` (String) A unique identifier for an object, represented as a UUID..
- `description_setting` (String) Whether the form asks the reporter to describe the incident: Required, Optional or Hidden...
- `custom_field_settings` (String) The incident custom fields the form asks for, keyed by each field's template variable key (variableKey). Required means the reporter must answer it, Optional that they may leave it empty. Only the fields listed as Required or Optional are asked: a field that is not listed, or is Hidden or Default, is not on the form. The answers become the incident's custom field values...
- `is_reporter_details_required` (Bool) When on, the reporter must give their name and email. When off, they may report anonymously...
- `success_message` (String) Shown to the reporter after they submit the form, together with the new incident's number. Markdown...
- `ip_whitelist` (String) The networks the form can be opened and submitted from: one IPv4 or IPv6 address, or one IPv4 range in CIDR notation (such as 10.0.0.0/8), per line. IPv6 ranges are not supported. Leave it empty to allow any network...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `share_key` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_form.example <id>
```
