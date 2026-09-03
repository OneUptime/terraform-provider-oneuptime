---
page_title: "oneuptime_network_device_role Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure what a device can be on your network (Router, Switch, Firewall and so on), how each role is drawn on the topology map, and which roles sit at the core.
---

# oneuptime_network_device_role (Resource)

Configure what a device can be on your network (Router, Switch, Firewall and so on), how each role is drawn on the topology map, and which roles sit at the core.

## Example Usage

```terraform
resource "oneuptime_network_device_role" "example" {
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
- `topology_shape` (String) The shape devices of this role are drawn with on the network topology map...
- `is_core_layer` (Bool) Devices of this role sit at the top of the network - the tiered and radial topology layouts band them above everything else...
- `is_snmp_walkable` (Bool) Devices of this role usually speak SNMP. Turn it off for roles that only answer a ping - adopting one from the topology map then defaults to a monitor rather than SNMP polling...
- `order` (Number) Where this role appears in the role picker and the topology map legend. Lower numbers come first...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `key` (String) Stable identifier for this role, derived from its name when it is created. SNMP classification and the topology map match on this, so it never changes when the role is renamed...
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_device_role.example <id>
```
