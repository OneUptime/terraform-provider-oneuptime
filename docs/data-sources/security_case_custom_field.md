---
page_title: "oneuptime_security_case_custom_field Data Source - oneuptime"
subcategory: "Security"
description: |-
  Fields a project adds to its security cases, filled in on each case.
---

# oneuptime_security_case_custom_field (Data Source)

Fields a project adds to its security cases, filled in on each case.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case custom field may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_custom_field" "example" {
  name = "Example security case custom field"
}

# Or by id:
data "oneuptime_security_case_custom_field" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who added the field. The ID of a `oneuptime_user` (see the data source).
- `custom_field_type` (String) What the field holds: text, a number, yes or no, a date, a choice from a list, and the rest.
- `description` (String) What the field is for, shown under it when a case's fields are edited.
- `dropdown_options` (String) Options and optional colors for dropdown fields, in the order they are listed. Plain one-per-line values remain supported. Records store an option as its text: changing an option here keeps the values records already hold. To rename an option and move those values with it, send the renames with the update, in miscDataProps: {"renamedDropdownOptions": [{"from": "Old text", "to": "New text"}]}. "to" must be one of the options.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from.
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand.
- `name` (String) The field's name, as each case shows it and as its value is kept under in the case's customFields.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this custom field belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
