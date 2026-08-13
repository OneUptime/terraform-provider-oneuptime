---
page_title: "oneuptime_network_endpoint Data Source - oneuptime"
subcategory: "Other"
description: |-
  LAN endpoints (POS terminals, kiosks, cameras, printers) discovered via ARP and FDB walks of Network Devices. Rows are upserted by the server; users can classify them.
---

# oneuptime_network_endpoint (Data Source)

LAN endpoints (POS terminals, kiosks, cameras, printers) discovered via ARP and FDB walks of Network Devices. Rows are upserted by the server; users can classify them. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_endpoint" "by_name" {
  name = "example-network_endpoint"
}

data "oneuptime_network_endpoint" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `mac_address` (String) MAC address of this endpoint, colon-separated hex. One row per MAC per project... Computed.
- `ip_address` (String) Last IP address seen for this endpoint in ARP tables. Managed by the server... Computed.
- `vendor` (String) Hardware vendor derived from the MAC OUI prefix. Managed by the server... Computed.
- `classification` (String) User-editable classification of this endpoint (POS, Kiosk, Camera, Printer, ...).. Computed.
- `attached_network_device_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `attached_interface_index` (Number) SNMP ifIndex of the switch port this endpoint was last seen on. Managed by the server... Computed.
- `attached_port_name` (String) Name of the switch port this endpoint was last seen on. Managed by the server... Computed.
- `vlan_id` (Number) VLAN this endpoint was last seen on, from the FDB walk. Managed by the server... Computed.
- `site_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `first_seen_at` (String) A date time object.. Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
