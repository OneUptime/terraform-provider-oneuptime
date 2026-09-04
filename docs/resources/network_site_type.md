---
page_title: "oneuptime_network_site_type Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure the levels of your network site hierarchy (Region, Market, Unit and so on). Choose each type's parent, rename it, or add your own.
---

# oneuptime_network_site_type (Resource)

Configure the levels of your network site hierarchy (Region, Market, Unit and so on). Choose each type's parent, rename it, or add your own.

## Example Usage

```terraform
resource "oneuptime_network_site_type" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `parent_network_site_type_id` (String) A unique identifier for an object, represented as a UUID..
- `order` (Number) Display order among site types that have the same parent...
- `is_unit_level` (Bool) Sites of this type are the leaf level - the network map opens their device topology, and the health rollup counts them as units...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_site_type.example <id>
```
