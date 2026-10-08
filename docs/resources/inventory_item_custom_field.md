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
  name        = "Example inventory item custom field"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `custom_field_type` (String) Is this field Text, Number or Boolean? A JSON value: write it with `jsonencode()`.
- `description` (String) Friendly description of this custom field that will help you remember.
- `dropdown_options` (String) Options and optional colors for dropdown fields. Plain one-per-line values remain supported.
- `map_from_custom_field_name` (String) Name of the custom field on the related resource this field copies its value from.
- `map_from_resource_type` (String) Related resource this field copies its value from. Empty means values are entered by hand.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing inventory item custom field by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_inventory_item_custom_field.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_inventory_item_custom_field.example <id>
```
