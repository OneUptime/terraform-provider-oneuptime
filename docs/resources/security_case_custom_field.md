---
page_title: "oneuptime_security_case_custom_field Resource - oneuptime"
subcategory: "Security"
description: |-
  Fields a project adds to its security cases, filled in on each case.
---

# oneuptime_security_case_custom_field (Resource)

Fields a project adds to its security cases, filled in on each case.

## Example Usage

```terraform
resource "oneuptime_security_case_custom_field" "example" {
  name        = "Example security case custom field"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) The field's name, as each case shows it and as its value is kept under in the case's customFields.

### Optional

- `custom_field_type` (String) What the field holds: text, a number, yes or no, a date, a choice from a list, and the rest. Allowed values: `Text`, `Number`, `Boolean`, `Dropdown`, `MultiSelectDropdown`, `Date`, `DateTime`, `LongText`, `Markdown`, `File`, `MultipleFiles`.
- `description` (String) What the field is for, shown under it when a case's fields are edited.
- `dropdown_options` (String) Options and optional colors for dropdown fields, in the order they are listed. Plain one-per-line values remain supported. Records store an option as its text: changing an option here keeps the values records already hold. To rename an option and move those values with it, send the renames with the update, in miscDataProps: {"renamedDropdownOptions": [{"from": "Old text", "to": "New text"}]}. "to" must be one of the options.
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from.
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand. Allowed values: `Monitor`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who added the field. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this custom field belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case custom field by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case_custom_field.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case_custom_field.example <id>
```
