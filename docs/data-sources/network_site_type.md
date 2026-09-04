---
page_title: "oneuptime_network_site_type Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure the levels of your network site hierarchy (Region, Market, Unit and so on). Choose each type's parent, rename it, or add your own.
---

# oneuptime_network_site_type (Data Source)

Configure the levels of your network site hierarchy (Region, Market, Unit and so on). Choose each type's parent, rename it, or add your own. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_site_type" "by_name" {
  name = "example-network_site_type"
}

data "oneuptime_network_site_type" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `parent_network_site_type_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `order` (Number) Display order among site types that have the same parent... Computed.
- `is_unit_level` (Bool) Sites of this type are the leaf level - the network map opens their device topology, and the health rollup counts them as units... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
