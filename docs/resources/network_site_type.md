---
page_title: "oneuptime_network_site_type Resource - oneuptime"
subcategory: "Network"
description: |-
  Configure the levels of your network site hierarchy (Region, Market, Unit and so on). Choose each type's parent, rename it, or add your own.
---

# oneuptime_network_site_type (Resource)

Configure the levels of your network site hierarchy (Region, Market, Unit and so on). Choose each type's parent, rename it, or add your own.

## Example Usage

```terraform
resource "oneuptime_network_site_type" "example" {
  name        = "Example network site type"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_unit_level` (Boolean) Sites of this type are the leaf level - the network map opens their device topology, and the health rollup counts them as units. Defaults to `false`.
- `order` (Number) Display order among site types that have the same parent.
- `parent_network_site_type_id` (String) ID of the Network Site Type directly above this type. Empty for top-level types. The ID of a `oneuptime_network_site_type`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network site type by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_site_type.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_site_type.example <id>
```
