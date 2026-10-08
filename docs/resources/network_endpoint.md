---
page_title: "oneuptime_network_endpoint Resource - oneuptime"
subcategory: "Other"
description: |-
  LAN endpoints (POS terminals, kiosks, cameras, printers) discovered via ARP and FDB walks of Network Devices. Rows are upserted by the server; users can classify them.
---

# oneuptime_network_endpoint (Resource)

LAN endpoints (POS terminals, kiosks, cameras, printers) discovered via ARP and FDB walks of Network Devices. Rows are upserted by the server; users can classify them.

## Example Usage

```terraform
resource "oneuptime_network_endpoint" "example" {
  mac_address = "Example short text"
}
```

## Schema

### Required

- `mac_address` (String) MAC address of this endpoint, colon-separated hex. One row per MAC per project.

### Optional

- `classification` (String) User-editable classification of this endpoint (POS, Kiosk, Camera, Printer, ...).
- `ip_address` (String) Last IP address seen for this endpoint in ARP tables. Managed by the server.
- `site_id` (String) ID of the Network Site this endpoint belongs to. The ID of a `oneuptime_network_site`.
- `vendor` (String) Hardware vendor derived from the MAC OUI prefix. Managed by the server.

### Read-Only

- `attached_interface_index` (Number) SNMP ifIndex of the switch port this endpoint was last seen on. Managed by the server.
- `attached_network_device_id` (String) ID of the Network Device this endpoint was last seen attached to. The ID of a `oneuptime_network_device`.
- `attached_port_name` (String) Name of the switch port this endpoint was last seen on. Managed by the server.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `first_seen_at` (String) When this endpoint was first discovered on the network.
- `id` (String) Unique identifier for the resource.
- `last_seen_at` (String) When this endpoint was last seen in an ARP or FDB walk.
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
