---
page_title: "oneuptime_incident_postmortem_template Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage postmortem templates for your incidents
---

# oneuptime_incident_postmortem_template (Resource)

Manage postmortem templates for your incidents

## Example Usage

```terraform
resource "oneuptime_incident_postmortem_template" "example" {
  template_name = "Example short text"
  template_description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `template_name` (String) Name of the Postmortem Template..
- `template_description` (String) Description of the Postmortem Template..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `postmortem_note` (String) Markdown template used when documenting an incident postmortem...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_postmortem_template.example <id>
```
