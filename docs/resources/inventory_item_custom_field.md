---
page_title: "oneuptime_inventory_item_custom_field Resource - oneuptime"
subcategory: "Other"
description: |-
  Manage custom fields on your inventory items
---

# oneuptime_inventory_item_custom_field (Resource)

Manage custom fields on your inventory items

## Example Usage

```terraform
resource "oneuptime_inventory_item_custom_field" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description of this custom field that will help you remember..
- `custom_field_type` (String) Is this field Text, Number or Boolean?..
- `dropdown_options` (String) Options and optional colors for dropdown fields. Plain one-per-line values remain supported...
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand...
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_inventory_item_custom_field.example <id>
```
