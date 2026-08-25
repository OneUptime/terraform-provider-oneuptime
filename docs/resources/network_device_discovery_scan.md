---
page_title: "oneuptime_network_device_discovery_scan Resource - oneuptime"
subcategory: "Other"
description: |-
  Network discovery scans that sweep an address space — a CIDR subnet or an octet range — via SNMP from a probe and report devices found, so they can be imported as Network Devices.
---

# oneuptime_network_device_discovery_scan (Resource)

Network discovery scans that sweep an address space — a CIDR subnet or an octet range — via SNMP from a probe and report devices found, so they can be imported as Network Devices.

## Example Usage

```terraform
resource "oneuptime_network_device_discovery_scan" "example" {
  probe_id = "123e4567-e89b-12d3-a456-426614174000"
  cidr = "Example short text"
}
```

## Schema

### Required

- `probe_id` (String) A unique identifier for an object, represented as a UUID..
- `cidr` (String) Address space to scan, either in CIDR notation (192.168.1.0/24) or octet-range notation where any octet may be an inclusive low-high range (10.16-22.0-255.51-66)..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `snmp_version` (String) SNMP version tried against every host in the subnet (V1, V2c, V3)..
- `snmp_community_string` (String) Community string tried against every host in the subnet (SNMP v1/v2c)..
- `snmp_port` (Number) UDP port tried against every host in the subnet..
- `snmp_v3_security_level` (String) SNMP v3 security level tried against every host: noAuthNoPriv, authNoPriv, or authPriv..
- `snmp_v3_username` (String) SNMP v3 security name (username) tried against every host..
- `snmp_v3_auth_protocol` (String) SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512..
- `snmp_v3_auth_key` (String) SNMP v3 authentication passphrase tried against every host..
- `snmp_v3_priv_protocol` (String) SNMP v3 privacy (encryption) protocol: DES, AES, or AES256..
- `snmp_v3_priv_key` (String) SNMP v3 privacy (encryption) passphrase tried against every host..
- `is_recurring` (Bool) Re-run this scan automatically every Rescan Interval minutes to keep discovery continuous...
- `rescan_interval_in_minutes` (Number) How often a recurring scan re-runs, in minutes. Ignored unless Is Recurring is on...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `status` (String) Status of this discovery scan: "Pending", "In Progress", "Completed" or "Failed". Managed by the scanning probe...
- `status_message` (String) Details about the current status of this scan, e.g. the failure reason. Managed by the scanning probe...
- `discovered_devices` (String) Devices found by this scan: array of {ipAddress, sysName, sysDescr, isAlreadyRegistered}. Managed by the scanning probe...
- `scanned_host_count` (Number) Total number of host addresses swept in the subnet. Managed by the scanning probe...
- `responded_host_count` (Number) Number of hosts that responded to SNMP during the sweep. Managed by the scanning probe...
- `started_at` (String) A date time object..
- `completed_at` (String) A date time object..
- `next_scan_at` (String) A date time object..
- `auto_import_processed_at` (String) A date time object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_device_discovery_scan.example <id>
```
