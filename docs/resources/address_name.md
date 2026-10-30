---
page_title: "oneuptime_address_name Resource - oneuptime"
subcategory: "Other"
description: |-
  Names for IP addresses and subnets that your inventory and DNS do not name. The Traffic pages show an address by its name, with the address beside it.
---

# oneuptime_address_name (Resource)

Names for IP addresses and subnets that your inventory and DNS do not name. The Traffic pages show an address by its name, with the address beside it.

## Example Usage

```terraform
resource "oneuptime_address_name" "example" {
  name    = "Example address name"
  address = "Example short text"
}
```

## Schema

### Required

- `address` (String) One IP address (10.0.0.5) or a subnet of them (10.18.167.0/24), IPv4 or IPv6. Each one is named once.
- `name` (String) What the Traffic pages call the address, or say a subnet's addresses are in: Payroll SaaS, Mumbai office.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing address name by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_address_name.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_address_name.example <id>
```
