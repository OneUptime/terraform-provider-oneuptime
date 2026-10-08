---
page_title: "oneuptime_network_device_link Data Source - oneuptime"
subcategory: "Other"
description: |-
  Operator-declared links between two Network Devices, for cables LLDP and CDP cannot see: a device with discovery disabled, a device that does not speak either protocol, or one monitored by ping alone. Drawn on the topology map alongside discovered links, and merged with a discovered link between the same pair rather than duplicating it.
---

# oneuptime_network_device_link (Data Source)

Operator-declared links between two Network Devices, for cables LLDP and CDP cannot see: a device with discovery disabled, a device that does not speak either protocol, or one monitored by ping alone. Drawn on the topology map alongside discovered links, and merged with a discovered link between the same pair rather than duplicating it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network device link may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_device_link" "example" {
  name = "Example network device link"
}

# Or by id:
data "oneuptime_network_device_link" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `from_device_id` (String) ID of the Network Device this link starts from. The ID of a `oneuptime_network_device`.
- `from_port_name` (String) Port on the starting device, as free text. Nothing resolves it to an interface row — a hand-drawn link usually exists precisely because the port is not discoverable.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `monitor_id` (String) ID of the Monitor whose status colors this link on the topology map. The ID of a `oneuptime_monitor`.
- `name` (String) Friendly name for this link.
- `parent_device_id` (String) ID of whichever end of this link is the parent. Must be the From Device or the To Device. Empty means the two are peers and the map infers the hierarchy. The ID of a `oneuptime_network_device`.
- `to_device_id` (String) ID of the Network Device this link ends at. The ID of a `oneuptime_network_device`.
- `to_port_name` (String) Port on the ending device, as free text.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
