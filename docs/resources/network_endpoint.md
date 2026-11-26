---
page_title: "oneuptime_network_endpoint Resource - oneuptime"
subcategory: "Network"
description: |-
  LAN endpoints (POS terminals, kiosks, cameras, printers) discovered via ARP and FDB walks of Network Devices, and from a network controller's client list. Rows are upserted by the server; users can classify them.
---

# oneuptime_network_endpoint (Resource)

LAN endpoints (POS terminals, kiosks, cameras, printers) discovered via ARP and FDB walks of Network Devices, and from a network controller's client list. Rows are upserted by the server; users can classify them.

## Example Usage

```terraform
resource "oneuptime_network_endpoint" "example" {
  mac_address = "aa:bb:cc:dd:ee:ff"
}
```

## Schema

### Required

- `mac_address` (String) MAC address of this endpoint, colon-separated hex. One row per MAC per project.

### Optional

- `classification` (String) User-editable classification of this endpoint (POS, Kiosk, Camera, Printer, ...).
- `site_id` (String) ID of the Network Site this endpoint belongs to. OneUptime changes it too, on its own: left out, it keeps the value OneUptime gives it. The ID of a `oneuptime_network_site`.
- `vendor` (String) Hardware vendor derived from the MAC OUI prefix. Managed by the server.

### Read-Only

- `attached_interface_index` (Number) SNMP ifIndex of the switch port this endpoint was last seen on. Managed by the server.
- `attached_network_device_id` (String) ID of the Network Device this endpoint was last seen attached to. The ID of a `oneuptime_network_device`.
- `attached_port_name` (String) Name of the switch port this endpoint was last seen on. Managed by the server.
- `attachment_source` (String) The network controller that reported where this endpoint is attached, such as meraki. Empty when a walk of a switch or router found it. Managed by the server.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `edited_at` (String) When a person last added this endpoint or changed its vendor, IP address or site. Managed by the server.
- `first_seen_at` (String) When this endpoint was first discovered on the network.
- `hostname` (String) What a network controller calls this endpoint, such as the name in its Meraki client list. Managed by the server.
- `id` (String) Unique identifier for the resource.
- `ip_address` (String) Last IP address seen for this endpoint in ARP tables. Managed by the server.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `last_seen_at` (String) When this endpoint was last seen in an ARP or FDB walk, or in a network controller's client list. While a walk keeps seeing it unchanged, this is refreshed every 15 minutes, so it is accurate to within 15 minutes and one poll.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `vlan_id` (Number) VLAN this endpoint was last seen on, from the FDB walk. Managed by the server.

## Import

Import an existing network endpoint by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_endpoint.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_endpoint.example <id>
```
