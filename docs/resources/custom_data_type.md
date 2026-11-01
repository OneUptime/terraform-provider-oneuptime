---
page_title: "oneuptime_custom_data_type Resource - oneuptime"
subcategory: "Other"
description: |-
  Your own kinds of sensitive data, such as employee IDs or policy numbers: a pattern, the words and field names that confirm it, and an optional checksum.
---

# oneuptime_custom_data_type (Resource)

Your own kinds of sensitive data, such as employee IDs or policy numbers: a pattern, the words and field names that confirm it, and an optional checksum.

## Example Usage

```terraform
resource "oneuptime_custom_data_type" "example" {
  name        = "Example custom data type"
  category    = "Example short text"
  sensitivity = "Example short text"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `category` (String) The category it belongs to.
- `name` (String) The data type's name, such as Employee ID.
- `sensitivity` (String) How sensitive it is: Restricted, Confidential, Internal or Public.

### Optional

- `description` (String) What the data type is.
- `ignore_case` (Boolean) Whether the pattern ignores upper and lower case. Defaults to `false`.
- `is_enabled` (Boolean) Whether scans look for this data type. Defaults to `true`.
- `keywords` (String) Words that, near a value, confirm it is this type. A JSON value: write it with `jsonencode()`.
- `name_hints` (String) Field, column and attribute names that hold this type, such as employee id. A JSON value: write it with `jsonencode()`.
- `pattern` (String) A JavaScript regular expression for one value. Leave empty to find it from field names alone.
- `requires_keyword` (Boolean) Count a value only when a keyword is nearby or the field's name hints at the type. Defaults to `false`.
- `validator` (String) A checksum a value must pass: None, Luhn, Mod97 or Verhoeff. Defaults to `None`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this custom data type. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this custom data type belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing custom data type by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_custom_data_type.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_custom_data_type.example <id>
```
