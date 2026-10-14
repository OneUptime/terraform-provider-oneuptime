---
page_title: "oneuptime_security_case_note Data Source - oneuptime"
subcategory: "Other"
description: |-
  Notes analysts write on a security case: what they checked, what they found and what comes next.
---

# oneuptime_security_case_note (Data Source)

Notes analysts write on a security case: what they checked, what they found and what comes next.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case note may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_note" "example" {
  security_case_id = oneuptime_security_case.example.id
}

# Or by id:
data "oneuptime_security_case_note" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who wrote the note. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `note` (String) The note, in Markdown.
- `security_case_id` (String) ID of the case this note is on. The ID of a `oneuptime_security_case`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this note belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
