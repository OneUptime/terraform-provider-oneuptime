---
page_title: "oneuptime_scheduled_maintenance_note_template Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage scheduled maintenance note templates for your project
---

# oneuptime_scheduled_maintenance_note_template (Resource)

Manage scheduled maintenance note templates for your project

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_note_template" "example" {
  template_name        = "Maintenance Completion Summary"
  template_description = "This template is used to document the completion of scheduled maintenance activities. It provides a structured format for summarizing the work performed, verifying system health, and communicating the maintenance outcome to stakeholders."
}
```

## Schema

### Required

- `template_description` (String) Description of the Incident Template.
- `template_name` (String) Name of the Incident Template.

### Optional

- `note` (String) Note template for public or private notes. This is in markdown.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing scheduled maintenance note template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_scheduled_maintenance_note_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_scheduled_maintenance_note_template.example <id>
```
