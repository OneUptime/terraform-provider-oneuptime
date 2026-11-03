---
page_title: "oneuptime_security_case_note_template Resource - oneuptime"
subcategory: "Security"
description: |-
  Ready-made text for the notes a team writes on its security cases, picked from a case's note composer.
---

# oneuptime_security_case_note_template (Resource)

Ready-made text for the notes a team writes on its security cases, picked from a case's note composer.

## Example Usage

```terraform
resource "oneuptime_security_case_note_template" "example" {
  template_name = "Containment checklist"
  note          = "# Heading\n\nThis is **markdown** content"
}
```

## Schema

### Required

- `note` (String) The text the template puts in a case's note, in Markdown, to edit before it is posted.
- `template_name` (String) What people pick the template by in a case's note composer.

### Optional

- `template_description` (String) What the template is for, and when to use it.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who added the template. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this note template belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case note template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case_note_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case_note_template.example <id>
```
