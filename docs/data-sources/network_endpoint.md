---
page_title: "oneuptime_network_endpoint Data Source - oneuptime"
subcategory: "Other"
description: |-
  LAN endpoints (POS terminals, kiosks, cameras, printers) discovered via ARP and FDB walks of Network Devices. Rows are upserted by the server; users can classify them.
---

# oneuptime_network_endpoint (Data Source)

LAN endpoints (POS terminals, kiosks, cameras, printers) discovered via ARP and FDB walks of Network Devices. Rows are upserted by the server; users can classify them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network endpoint may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_endpoint" "example" {
  mac_address = "example-mac-address"
}

# Or by id:
data "oneuptime_network_endpoint" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `attached_interface_index` (Number) SNMP ifIndex of the switch port this endpoint was last seen on. Managed by the server.
- `attached_network_device_id` (String) ID of the Network Device this endpoint was last seen attached to. The ID of a `oneuptime_network_device`.
- `attached_port_name` (String) Name of the switch port this endpoint was last seen on. Managed by the server.
- `classification` (String) User-editable classification of this endpoint (POS, Kiosk, Camera, Printer, ...).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `ip_address` (String) Last IP address seen for this endpoint in ARP tables. Managed by the server.
- `mac_address` (String) MAC address of this endpoint, colon-separated hex. One row per MAC per project.
- `site_id` (String) ID of the Network Site this endpoint belongs to. The ID of a `oneuptime_network_site`.
- `vendor` (String) Hardware vendor derived from the MAC OUI prefix. Managed by the server.
- `vlan_id` (Number) VLAN this endpoint was last seen on, from the FDB walk. Managed by the server.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `first_seen_at` (String) When this endpoint was first discovered on the network.
- `last_seen_at` (String) When this endpoint was last seen in an ARP or FDB walk.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
