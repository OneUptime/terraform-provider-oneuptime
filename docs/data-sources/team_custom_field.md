---
page_title: "oneuptime_team_custom_field Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Manage custom fields for your teams
---

# oneuptime_team_custom_field (Data Source)

Manage custom fields for your teams

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team custom field may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team_custom_field" "example" {
  name = "Example team custom field"
}

# Or by id:
data "oneuptime_team_custom_field" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description of this custom field that will help you remember.
- `dropdown_options` (String) Options and optional colors for dropdown fields, in the order they are listed. Plain one-per-line values remain supported. Records store an option as its text: changing an option here keeps the values records already hold. To rename an option and move those values with it, send the renames with the update, in miscDataProps: {"renamedDropdownOptions": [{"from": "Old text", "to": "New text"}]}. "to" must be one of the options.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from.
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand.
- `name` (String) Any friendly name of this object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `custom_field_type` (String) Is this field Text, Number, Boolean or Dropdown? A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
