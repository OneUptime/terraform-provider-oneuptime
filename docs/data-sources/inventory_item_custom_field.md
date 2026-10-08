---
page_title: "oneuptime_inventory_item_custom_field Data Source - oneuptime"
subcategory: "Other"
description: |-
  Manage custom fields on your inventory items
---

# oneuptime_inventory_item_custom_field (Data Source)

Manage custom fields on your inventory items

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one inventory item custom field may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_inventory_item_custom_field" "example" {
  name = "Example inventory item custom field"
}

# Or by id:
data "oneuptime_inventory_item_custom_field" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description of this custom field that will help you remember.
- `dropdown_options` (String) Options and optional colors for dropdown fields. Plain one-per-line values remain supported.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from.
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand.
- `name` (String) Any friendly name of this object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `custom_field_type` (String) Is this field Text, Number or Boolean? A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
