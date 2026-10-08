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
  template_name        = "Example short text"
  template_description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `template_description` (String) Description of the Postmortem Template.
- `template_name` (String) Name of the Postmortem Template.

### Optional

- `postmortem_note` (String) Markdown template used when documenting an incident postmortem.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident postmortem template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_postmortem_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_postmortem_template.example <id>
```
