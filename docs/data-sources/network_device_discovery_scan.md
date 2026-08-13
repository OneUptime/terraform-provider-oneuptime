---
page_title: "oneuptime_network_device_discovery_scan Data Source - oneuptime"
subcategory: "Other"
description: |-
  Network discovery scans that sweep an address space — a CIDR subnet or an octet range — via SNMP from a probe and report devices found, so they can be imported as Network Devices.
---

# oneuptime_network_device_discovery_scan (Data Source)

Network discovery scans that sweep an address space — a CIDR subnet or an octet range — via SNMP from a probe and report devices found, so they can be imported as Network Devices. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_device_discovery_scan" "by_name" {
  name = "example-network_device_discovery_scan"
}

data "oneuptime_network_device_discovery_scan" "by_id" {
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
- `probe_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `cidr` (String) Address space to scan, either in CIDR notation (192.168.1.0/24) or octet-range notation where any octet may be an inclusive low-high range (10.16-22.0-255.51-66).. Computed.
- `snmp_version` (String) SNMP version tried against every host in the subnet (V1, V2c, V3).. Computed.
- `snmp_community_string` (String) Community string tried against every host in the subnet (SNMP v1/v2c).. Computed.
- `snmp_port` (Number) UDP port tried against every host in the subnet.. Computed.
- `snmp_v3_security_level` (String) SNMP v3 security level tried against every host: noAuthNoPriv, authNoPriv, or authPriv.. Computed.
- `snmp_v3_username` (String) SNMP v3 security name (username) tried against every host.. Computed.
- `snmp_v3_auth_protocol` (String) SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512.. Computed.
- `snmp_v3_auth_key` (String) SNMP v3 authentication passphrase tried against every host.. Computed.
- `snmp_v3_priv_protocol` (String) SNMP v3 privacy (encryption) protocol: DES, AES, or AES256.. Computed.
- `snmp_v3_priv_key` (String) SNMP v3 privacy (encryption) passphrase tried against every host.. Computed.
- `status` (String) Status of this discovery scan: "Pending", "In Progress", "Completed" or "Failed". Managed by the scanning probe... Computed.
- `status_message` (String) Details about the current status of this scan, e.g. the failure reason. Managed by the scanning probe... Computed.
- `discovered_devices` (String) Devices found by this scan: array of {ipAddress, sysName, sysDescr, isAlreadyRegistered}. Managed by the scanning probe... Computed.
- `scanned_host_count` (Number) Total number of host addresses swept in the subnet. Managed by the scanning probe... Computed.
- `responded_host_count` (Number) Number of hosts that responded to SNMP during the sweep. Managed by the scanning probe... Computed.
- `started_at` (String) A date time object.. Computed.
- `completed_at` (String) A date time object.. Computed.
- `is_recurring` (Bool) Re-run this scan automatically every Rescan Interval minutes to keep discovery continuous... Computed.
- `rescan_interval_in_minutes` (Number) How often a recurring scan re-runs, in minutes. Ignored unless Is Recurring is on... Computed.
- `next_scan_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
