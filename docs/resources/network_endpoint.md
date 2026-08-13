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

- `mac_address` (String) MAC address of this endpoint, colon-separated hex. One row per MAC per project...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `classification` (String) User-editable classification of this endpoint (POS, Kiosk, Camera, Printer, ...)..
- `site_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `ip_address` (String) Last IP address seen for this endpoint in ARP tables. Managed by the server...
- `vendor` (String) Hardware vendor derived from the MAC OUI prefix. Managed by the server...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `attached_network_device_id` (String) A unique identifier for an object, represented as a UUID..
- `attached_interface_index` (Number) SNMP ifIndex of the switch port this endpoint was last seen on. Managed by the server...
- `attached_port_name` (String) Name of the switch port this endpoint was last seen on. Managed by the server...
- `vlan_id` (Number) VLAN this endpoint was last seen on, from the FDB walk. Managed by the server...
- `first_seen_at` (String) A date time object..
- `last_seen_at` (String) A date time object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_endpoint.example <id>
```
