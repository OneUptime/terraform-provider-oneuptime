---
page_title: "oneuptime_network_site_type Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure the levels of your network site hierarchy (Region, Market, Unit and so on). Choose each type's parent, rename it, or add your own.
---

# oneuptime_network_site_type (Data Source)

Configure the levels of your network site hierarchy (Region, Market, Unit and so on). Choose each type's parent, rename it, or add your own.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network site type may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_site_type" "example" {
  name = "Example network site type"
}

# Or by id:
data "oneuptime_network_site_type" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_unit_level` (Boolean) Sites of this type are the leaf level - the network map opens their device topology, and the health rollup counts them as units.
- `name` (String) Any friendly name of this object.
- `order` (Number) Display order among site types that have the same parent.
- `parent_network_site_type_id` (String) ID of the Network Site Type directly above this type. Empty for top-level types. The ID of a `oneuptime_network_site_type`.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
