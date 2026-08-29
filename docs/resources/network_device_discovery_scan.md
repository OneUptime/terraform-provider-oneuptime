---
page_title: "oneuptime_network_device_discovery_scan Resource - oneuptime"
subcategory: "Other"
description: |-
  Network discovery scans that sweep an address space — a CIDR subnet or an octet range — from a probe and report the hosts found, so they can be imported as Network Devices. Every sweep pings; scans with Check SNMP on also query each live host over SNMP.
---

# oneuptime_network_device_discovery_scan (Resource)

Network discovery scans that sweep an address space — a CIDR subnet or an octet range — from a probe and report the hosts found, so they can be imported as Network Devices. Every sweep pings; scans with Check SNMP on also query each live host over SNMP.

## Example Usage

```terraform
resource "oneuptime_network_device_discovery_scan" "example" {
  probe_id = "123e4567-e89b-12d3-a456-426614174000"
  cidr = "Example short text"
  name = "Example short text"
}
```

## Schema

### Required

- `probe_id` (String) A unique identifier for an object, represented as a UUID..
- `cidr` (String) Address space to scan, either in CIDR notation (192.168.1.0/24) or octet-range notation where any octet may be an inclusive low-high range (10.16-22.0-255.51-66)..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Optional name for this scan, so it can be told apart from other scans at a glance. Falls back to the scan target when empty...
- `snmp_configs` (String) Ordered list of SNMP credential sets tried against every host in the subnet, first match wins. Each entry carries an id, an optional name, a version, a community string or the v3 credentials, and a port. When empty, the scan uses the single flattened SNMP configuration on this row...
- `is_snmp_enabled` (Bool) Whether hosts that answer the ping sweep are then queried over SNMP. Turn it off for an ICMP-only scan, which reports every host that answers ping and asks nothing else of them...
- `snmp_version` (String) SNMP version tried against every host in the subnet (V1, V2c, V3). Ignored when Check SNMP is off...
- `snmp_community_string` (String) Community string tried against every host in the subnet (SNMP v1/v2c). Ignored when Check SNMP is off...
- `snmp_port` (Number) UDP port tried against every host in the subnet. Ignored when Check SNMP is off...
- `snmp_v3_security_level` (String) SNMP v3 security level tried against every host: noAuthNoPriv, authNoPriv, or authPriv. Ignored when Check SNMP is off...
- `snmp_v3_username` (String) SNMP v3 security name (username) tried against every host. Ignored when Check SNMP is off...
- `snmp_v3_auth_protocol` (String) SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512. Ignored when Check SNMP is off...
- `snmp_v3_auth_key` (String) SNMP v3 authentication passphrase tried against every host. Ignored when Check SNMP is off...
- `snmp_v3_priv_protocol` (String) SNMP v3 privacy (encryption) protocol: DES, AES, or AES256. Ignored when Check SNMP is off...
- `snmp_v3_priv_key` (String) SNMP v3 privacy (encryption) passphrase tried against every host. Ignored when Check SNMP is off...
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
- `responded_host_count` (Number) Number of hosts that answered the check this scan performed: SNMP responders on a scan with Check SNMP on, hosts that answered the ping sweep on an ICMP-only one. Managed by the scanning probe...
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
