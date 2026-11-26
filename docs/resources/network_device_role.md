---
page_title: "oneuptime_network_device_role Resource - oneuptime"
subcategory: "Network"
description: |-
  Configure what a device can be on your network (Router, Switch, Firewall and so on), how each role is drawn on the topology map, and which roles sit at the core.
---

# oneuptime_network_device_role (Resource)

Configure what a device can be on your network (Router, Switch, Firewall and so on), how each role is drawn on the topology map, and which roles sit at the core.

## Example Usage

```terraform
resource "oneuptime_network_device_role" "example" {
  name        = "Example network device role"
  description = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this network device role's place among the network device roles of its project, lowest first. Setting it places this network device role and moves no other, so a configuration that numbers every network device role of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.network_device_roles`:

```terraform
order = index(local.network_device_roles, each.key) + 1
```

- Give each network device role a number of its own: two network device roles at one number are shown the older first.
- A network device role created without `order` goes after the last one of its project.
- Deleting a network device role leaves a gap in the numbers, which changes nothing: network device roles go by their numbers, not by counting.
- Moving a network device role in the OneUptime Dashboard numbers the project's network device roles 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_core_layer` (Boolean) Devices of this role sit at the top of the network - the tiered and radial topology layouts band them above everything else. Defaults to `false`.
- `is_snmp_walkable` (Boolean) Devices of this role usually speak SNMP. Turn it off for roles that only answer a ping - adopting one from the topology map then defaults to a monitor rather than SNMP polling. Defaults to `true`.
- `order` (Number) Where this role appears in the role picker and the topology map legend, lowest number first. A new role is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).
- `topology_shape` (String) The shape devices of this role are drawn with on the network topology map.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `key` (String) Stable identifier for this role, derived from its name when it is created. SNMP classification and the topology map match on this, so it never changes when the role is renamed.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network device role by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_device_role.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_device_role.example <id>
```
