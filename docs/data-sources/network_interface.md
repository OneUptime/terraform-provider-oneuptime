---
page_title: "oneuptime_network_interface Data Source - oneuptime"
subcategory: "Other"
description: |-
  Interfaces discovered on Network Devices via SNMP walks. Rows are upserted by the server; users can toggle per-interface monitoring.
---

# oneuptime_network_interface (Data Source)

Interfaces discovered on Network Devices via SNMP walks. Rows are upserted by the server; users can toggle per-interface monitoring. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_interface" "by_name" {
  name = "example-network_interface"
}

data "oneuptime_network_interface" "by_id" {
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
- `network_device_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `interface_index` (Number) SNMP ifIndex of this interface on the device.. Computed.
- `alias` (String) Interface alias (ifAlias) from SNMP.. Computed.
- `mac_address` (String) Physical address (ifPhysAddress) from SNMP, colon-separated hex.. Computed.
- `interface_type` (Number) IANAifType number (ifType) from SNMP — 6 = ethernetCsmacd, 24 = softwareLoopback.. Computed.
- `is_monitored` (Bool) Include this interface in down/utilization/error alerting... Computed.
- `is_operationally_up` (Bool) Operational status (ifOperStatus) from the last SNMP walk.. Computed.
- `is_administratively_up` (Bool) Administrative status (ifAdminStatus) from the last SNMP walk.. Computed.
- `speed_in_mbps` (Number) Negotiated interface speed in Mbps. Stored as decimal so 10G+ links don't overflow integers... Computed.
- `in_rate_mbps` (Number) Most recent inbound throughput in Mbps, computed from SNMP counters... Computed.
- `out_rate_mbps` (Number) Most recent outbound throughput in Mbps, computed from SNMP counters... Computed.
- `utilization_percent` (Number) Most recent utilization as a percent of interface speed (max of in/out)... Computed.
- `errors_per_second` (Number) Most recent error rate (in + out errors per second) computed from SNMP counters... Computed.
- `last_seen_at` (String) A date time object.. Computed.
