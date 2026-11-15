---
page_title: "oneuptime_network_interface Data Source - oneuptime"
subcategory: "Network"
description: |-
  Interfaces discovered on Network Devices via SNMP walks. Rows are upserted by the server; users can toggle per-interface monitoring.
---

# oneuptime_network_interface (Data Source)

Interfaces discovered on Network Devices via SNMP walks. Rows are upserted by the server; users can toggle per-interface monitoring.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network interface may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_interface" "example" {
  name = "Example network interface"
}

# Or by id:
data "oneuptime_network_interface" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alias` (String) Interface alias (ifAlias) from SNMP.
- `errors_per_second` (Number) Error rate (in + out errors per second) computed from SNMP counters between two polls. Written again when errors start or stop, when it moves by a tenth or more, and at least every 15 minutes.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `in_rate_mbps` (Number) Inbound throughput in Mbps, computed from SNMP counters between two polls. Written again when it moves by a tenth or more (and by more than 1% of the link speed), and at least every 15 minutes.
- `interface_index` (Number) SNMP ifIndex of this interface on the device.
- `interface_type` (Number) IANAifType number (ifType) from SNMP — 6 = ethernetCsmacd, 24 = softwareLoopback.
- `is_administratively_up` (Boolean) Administrative status (ifAdminStatus) from the last SNMP walk.
- `is_monitored` (Boolean) Include this interface in down/utilization/error alerting.
- `is_operationally_up` (Boolean) Operational status (ifOperStatus) from the last SNMP walk.
- `mac_address` (String) Physical address (ifPhysAddress) from SNMP, colon-separated hex.
- `name` (String) Interface name (ifName / ifDescr) from SNMP.
- `network_device_id` (String) ID of the Network Device this interface was discovered on. The ID of a `oneuptime_network_device`.
- `out_rate_mbps` (Number) Outbound throughput in Mbps, computed from SNMP counters between two polls. Written again when it moves by a tenth or more (and by more than 1% of the link speed), and at least every 15 minutes.
- `speed_in_mbps` (Number) Negotiated interface speed in Mbps. Stored as decimal so 10G+ links don't overflow integers.
- `utilization_percent` (Number) Utilization as a percent of interface speed (max of in/out) between two polls. Written again when it moves by a tenth or more (and by more than a point), crosses 80%, and at least every 15 minutes.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_seen_at` (String) When an SNMP walk last saw this interface, to within 15 minutes: a walk that finds an interface unchanged writes it again only every 15 minutes.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
