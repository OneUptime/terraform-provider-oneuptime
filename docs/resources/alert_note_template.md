---
page_title: "oneuptime_alert_note_template Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Manage alert note templates for your project
---

# oneuptime_alert_note_template (Resource)

Manage alert note templates for your project

## Example Usage

```terraform
resource "oneuptime_alert_note_template" "example" {
  template_name        = "Initial Investigation Template"
  template_description = "This template is used for documenting the initial investigation phase of an alert. It includes sections for analyzing the issue, documenting findings, and outlining next steps for resolution."
}
```

## Schema

### Required

- `template_description` (String) Description of the Alert Template.
- `template_name` (String) Name of the Alert Template.

### Optional

- `note` (String) Note template for public or private notes. This is in markdown.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert note template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_note_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_note_template.example <id>
```
