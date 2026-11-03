---
page_title: "oneuptime_security_case_note_template Data Source - oneuptime"
subcategory: "Security"
description: |-
  Ready-made text for the notes a team writes on its security cases, picked from a case's note composer.
---

# oneuptime_security_case_note_template (Data Source)

Ready-made text for the notes a team writes on its security cases, picked from a case's note composer.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case note template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_note_template" "example" {
  template_name = "example-template-name"
}

# Or by id:
data "oneuptime_security_case_note_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who added the template. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `note` (String) The text the template puts in a case's note, in Markdown, to edit before it is posted.
- `template_description` (String) What the template is for, and when to use it.
- `template_name` (String) What people pick the template by in a case's note composer.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this note template belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
