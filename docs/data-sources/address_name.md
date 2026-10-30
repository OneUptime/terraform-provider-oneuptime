---
page_title: "oneuptime_address_name Data Source - oneuptime"
subcategory: "Other"
description: |-
  Names for IP addresses and subnets that your inventory and DNS do not name. The Traffic pages show an address by its name, with the address beside it.
---

# oneuptime_address_name (Data Source)

Names for IP addresses and subnets that your inventory and DNS do not name. The Traffic pages show an address by its name, with the address beside it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one address name may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_address_name" "example" {
  name = "Example address name"
}

# Or by id:
data "oneuptime_address_name" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `address` (String) One IP address (10.0.0.5) or a subnet of them (10.18.167.0/24), IPv4 or IPv6. Each one is named once.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) What the Traffic pages call the address, or say a subnet's addresses are in: Payroll SaaS, Mumbai office.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
