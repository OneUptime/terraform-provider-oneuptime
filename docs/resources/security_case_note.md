---
page_title: "oneuptime_security_case_note Resource - oneuptime"
subcategory: "Other"
description: |-
  Notes analysts write on a security case: what they checked, what they found and what comes next.
---

# oneuptime_security_case_note (Resource)

Notes analysts write on a security case: what they checked, what they found and what comes next.

## Example Usage

```terraform
resource "oneuptime_security_case_note" "example" {
  security_case_id = oneuptime_security_case.example.id
  note             = "# Heading\n\nThis is **markdown** content"
}
```

## Schema

### Required

- `note` (String) The note, in Markdown.
- `security_case_id` (String) ID of the case this note is on. The ID of a `oneuptime_security_case`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who wrote the note. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this note belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case note by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case_note.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case_note.example <id>
```
