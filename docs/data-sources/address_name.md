---
page_title: "oneuptime_address_name Data Source - oneuptime"
subcategory: "Network"
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
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) What the Traffic pages call the address, or say a subnet's addresses are in: Payroll SaaS, Mumbai office.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
