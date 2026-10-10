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
  probe_id = oneuptime_probe.example.id
  cidr     = "Example short text"
  name     = "Example network device discovery scan"
}
```

## Schema

### Required

- `cidr` (String) Address space to scan, either in CIDR notation (192.168.1.0/24) or octet-range notation where any octet may be an inclusive low-high range (10.16-22.0-255.51-66).
- `probe_id` (String) ID of the Probe that runs this discovery scan. The ID of a `oneuptime_probe`.

### Optional

- `is_netbios_lookup_enabled` (Boolean) Whether hosts with no SNMP name are asked for their NetBIOS name over UDP 137, including hosts reverse DNS already named. A NetBIOS name is the name the host reports for itself, so it names the device ahead of its reverse DNS name. Best-effort: Windows/Samba hosts that allow UDP 137 from the probe. Private addresses only; never done by global probes. Defaults to `false`.
- `is_recurring` (Boolean) Re-run this scan automatically every Rescan Interval minutes to keep discovery continuous. Defaults to `false`.
- `is_snmp_enabled` (Boolean) Whether hosts that answer the ping sweep are then queried over SNMP. Turn it off for an ICMP-only scan, which reports every host that answers ping and asks nothing else of them. Defaults to `true`.
- `name` (String) Optional name for this scan, so it can be told apart from other scans at a glance. Falls back to the scan target when empty.
- `rescan_interval_in_minutes` (Number) How often a recurring scan re-runs, in minutes. Ignored unless Is Recurring is on.
- `snmp_community_string` (String) Community string tried against every host in the subnet (SNMP v1/v2c). Ignored when Check SNMP is off.
- `snmp_configs` (String) Ordered list of SNMP credential sets tried against every host in the subnet, first match wins. Each entry carries an id, an optional name, a version, a community string or the v3 credentials, and a port. When empty, the scan uses the single flattened SNMP configuration on this row. A JSON value: write it with `jsonencode()`.
- `snmp_port` (Number) UDP port tried against every host in the subnet. Ignored when Check SNMP is off.
- `snmp_v3_auth_key` (String) SNMP v3 authentication passphrase tried against every host. Ignored when Check SNMP is off.
- `snmp_v3_auth_protocol` (String) SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512. Ignored when Check SNMP is off.
- `snmp_v3_priv_key` (String) SNMP v3 privacy (encryption) passphrase tried against every host. Ignored when Check SNMP is off.
- `snmp_v3_priv_protocol` (String) SNMP v3 privacy (encryption) protocol: DES, AES, or AES256. Ignored when Check SNMP is off.
- `snmp_v3_security_level` (String) SNMP v3 security level tried against every host: noAuthNoPriv, authNoPriv, or authPriv. Ignored when Check SNMP is off.
- `snmp_v3_username` (String) SNMP v3 security name (username) tried against every host. Ignored when Check SNMP is off.
- `snmp_version` (String) SNMP version tried against every host in the subnet (V1, V2c, V3). Ignored when Check SNMP is off.
- `use_short_device_names` (Boolean) Name imported devices by their short hostname (the first label of a fully qualified name, e.g. 'core-sw-01' rather than 'core-sw-01.corp.example.com'). The full reverse-DNS name is still stored on the device as its DNS Name. Defaults to `false`.

### Read-Only

- `auto_import_processed_at` (String) When auto-import rules last processed this scan's results. Managed by the server: cleared when new results arrive, stamped by the worker that evaluates the rules. NULL means the current results have not been processed yet.
- `completed_at` (String) When the scanning probe completed (or failed) this scan. Managed by the scanning probe.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `discovered_devices` (String) Devices found by this scan: array of {ipAddress, sysName, sysDescr, isAlreadyRegistered}. Managed by the scanning probe. A JSON value: write it with `jsonencode()`.
- `id` (String) Unique identifier for the resource.
- `next_scan_at` (String) When a recurring scan is next due to run. Managed by the server.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `responded_host_count` (Number) Number of hosts that answered the check this scan performed: SNMP responders on a scan with Check SNMP on, hosts that answered the ping sweep on an ICMP-only one. Managed by the scanning probe.
- `scanned_host_count` (Number) Total number of host addresses swept in the subnet. Managed by the scanning probe.
- `started_at` (String) When the scanning probe started this scan. Managed by the scanning probe.
- `status` (String) Status of this discovery scan: "Pending", "In Progress", "Completed" or "Failed". Managed by the scanning probe.
- `status_message` (String) Details about the current status of this scan, e.g. the failure reason. Managed by the scanning probe.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network device discovery scan by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_device_discovery_scan.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_device_discovery_scan.example <id>
```
