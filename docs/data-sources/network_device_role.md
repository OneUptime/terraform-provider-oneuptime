---
page_title: "oneuptime_network_device_role Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure what a device can be on your network (Router, Switch, Firewall and so on), how each role is drawn on the topology map, and which roles sit at the core.
---

# oneuptime_network_device_role (Data Source)

Configure what a device can be on your network (Router, Switch, Firewall and so on), how each role is drawn on the topology map, and which roles sit at the core.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network device role may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_device_role" "example" {
  name = "Example network device role"
}

# Or by id:
data "oneuptime_network_device_role" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_core_layer` (Boolean) Devices of this role sit at the top of the network - the tiered and radial topology layouts band them above everything else.
- `is_snmp_walkable` (Boolean) Devices of this role usually speak SNMP. Turn it off for roles that only answer a ping - adopting one from the topology map then defaults to a monitor rather than SNMP polling.
- `key` (String) Stable identifier for this role, derived from its name when it is created. SNMP classification and the topology map match on this, so it never changes when the role is renamed.
- `name` (String) Any friendly name of this object.
- `order` (Number) Where this role appears in the role picker and the topology map legend, lowest number first. A new role is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `slug` (String) Friendly globally unique name for your object.
- `topology_shape` (String) The shape devices of this role are drawn with on the network topology map.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
