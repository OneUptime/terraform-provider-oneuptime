---
page_title: "oneuptime_network_device Data Source - oneuptime"
subcategory: "Other"
description: |-
  Network Devices (routers, switches, firewalls) that are being monitored in this project via SNMP polling and traps.
---

# oneuptime_network_device (Data Source)

Network Devices (routers, switches, firewalls) that are being monitored in this project via SNMP polling and traps.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network device may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_device" "example" {
  name = "Example network device"
}

# Or by id:
data "oneuptime_network_device" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `auto_apply_vendor_health_template` (Boolean) When the device's vendor is fingerprinted from its SNMP sysObjectID and no Health OIDs are configured yet, apply the matching vendor health template automatically on the next poll. Off by default for hand-made devices — the vendor template banner stays the manual path; auto-imported devices enable it so the zero-touch pipeline ends with health metrics, not an empty OID list.
- `collect_endpoints` (Boolean) Also walk the device's ARP cache and bridge forwarding database on each poll to discover endpoints (laptops, printers, POS terminals) attached to it. Strictly opt-in: costs extra SNMP table walks per poll. Only meaningful when Walk Interfaces is on.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_monitor_status_id` (String) Whats the current status ID of this network device? Stamped from the monitor that polls it. The ID of a `oneuptime_monitor_status`.
- `description` (String) Friendly description for this network device.
- `device_model` (String) Hardware model from ENTITY-MIB (entPhysicalModelName). Managed by the probe.
- `device_role` (String) Deprecated legacy device role key. Use the Network Device Role relation instead; this column exists only for the backfill migration and will be removed.
- `dns_name` (String) Fully qualified DNS name of this device, from its reverse-DNS (PTR) record when it was discovered, or its previous full name when its name was shortened to the hostname. Kept so the device can still be found, and matched by site-assignment hostname patterns, by the name DNS gives it.
- `firmware_version` (String) Firmware revision from ENTITY-MIB (entPhysicalFirmwareRev). Managed by the probe.
- `hostname` (String) IP address or hostname the probe polls; also matches SNMP trap sources.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `interfaces_down` (Number) Cached count of operationally down interfaces on this device.
- `interfaces_total` (Number) Cached total count of interfaces on this device.
- `interfaces_up` (Number) Cached count of operationally up interfaces on this device.
- `is_archived` (Boolean) Is this network device archived? Archived network devices are hidden from lists but keep collecting telemetry.
- `is_mac_address_learned` (Boolean) True when the MAC Address was filled in from a walked device's ARP table rather than typed. A learned MAC is corrected when a later walk binds the device's address to a different MAC; a typed one is never touched.
- `is_polling_enabled` (Boolean) Whether the assigned probe polls this device on a schedule. Disable to pause SNMP polling without deleting the device.
- `is_reachable` (Boolean) Whether the most recent SNMP walk reached this device. NULL means it has never been polled. This — not the age of lastSeenAt — is what the device list, the topology graph and the site rollup read, so a device whose last poll succeeded is never shown as down just because the probe is behind schedule. Managed by the probe.
- `is_snmp_reachable` (Boolean) Whether the most recent SNMP walk of this device succeeded. Separate from isReachable, which is the ping verdict: a device that answers ping but not SNMP is Up with a failing SNMP walk, which almost always means wrong credentials or SNMP disabled on the device. NULL means no walk was attempted — the device has no usable SNMP credentials (it is pinged only) or has never been polled. Managed by the probe.
- `mac_address` (String) MAC address of this device. Lets the topology map find the switch port it is plugged into from the forwarding tables of walked switches, for a device that speaks neither LLDP nor CDP (one monitored by ping alone). Optional: a device whose hostname is an IP address that a walked router's ARP table resolves is matched by address, and the MAC learned that way is stored here.
- `monitor_id` (String) ID of the monitor that reports this device's health when its monitoring method is Monitor. The ID of a `oneuptime_monitor`.
- `monitoring_method` (String) How this device's health is established: SNMP (an assigned probe walks it on a schedule) or Monitor (no polling — the linked monitor's status is the device's status). Devices created before this existed are SNMP.
- `name` (String) Friendly name for this network device.
- `network_device_role_id` (String) ID of the Network Device Role this device is assigned. The ID of a `oneuptime_network_device_role`.
- `oid_template_id` (String) ID of the OID Collection Template this device collects. The ID of a `oneuptime_oid_collection_template`.
- `polling_interval_in_minutes` (Number) How often, in minutes, the assigned probe polls this device via SNMP.
- `probe_id` (String) ID of the Probe that polls this network device. The ID of a `oneuptime_probe`.
- `serial_number` (String) Chassis serial number from ENTITY-MIB (entPhysicalSerialNum). Managed by the probe.
- `site_id` (String) ID of the Network Site this network device belongs to. The ID of a `oneuptime_network_site`.
- `slug` (String) Friendly globally unique name for your object.
- `snmp_community_string` (String) Community string used for SNMP v1/v2c polling.
- `snmp_credential_profile_id` (String) ID of the SNMP Credential Profile this device is walked with. The ID of a `oneuptime_snmp_credential_profile`.
- `snmp_port` (Number) UDP port used for SNMP polling.
- `snmp_v3_auth_key` (String) SNMP v3 authentication passphrase.
- `snmp_v3_auth_protocol` (String) SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512.
- `snmp_v3_priv_key` (String) SNMP v3 privacy (encryption) passphrase.
- `snmp_v3_priv_protocol` (String) SNMP v3 privacy (encryption) protocol: DES, AES, or AES256.
- `snmp_v3_security_level` (String) SNMP v3 security level: noAuthNoPriv, authNoPriv, or authPriv.
- `snmp_v3_username` (String) Security name (username) used for SNMP v3 polling.
- `snmp_version` (String) SNMP version to use when polling this device (V1, V2c, V3).
- `software_version` (String) Operating system / software revision from ENTITY-MIB (entPhysicalSoftwareRev). Managed by the probe.
- `sys_contact` (String) System contact (sysContact) enriched from SNMP walks of this device.
- `sys_descr` (String) System description (sysDescr) enriched from SNMP walks of this device.
- `sys_location` (String) System location (sysLocation) enriched from SNMP walks of this device.
- `sys_name` (String) System name (sysName) enriched from SNMP walks of this device.
- `sys_object_id` (String) sysObjectID — the vendor's registered OID for this device model, enriched from SNMP walks. Used to fingerprint the vendor and suggest an OID template.
- `vendor` (String) Hardware vendor, from ENTITY-MIB or derived from sysObjectID. Managed by the probe.
- `walk_interfaces` (Boolean) Walk the IF-MIB interface tables on each poll to inventory interfaces, bandwidth, and errors. Also collects LLDP/CDP neighbors for the topology graph.

### Read-Only

- `archived_at` (String) When was this network device archived?
- `cdp_neighbors` (String) CDP neighbors discovered on the last SNMP walk, complementing LLDP for the topology graph. Managed by the probe. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `last_polled_at` (String) When the assigned probe last ATTEMPTED an SNMP walk of this device, whether or not the device answered. Paired with lastSeenAt (which only moves on a successful walk) this is what separates "the device did not answer" from "we have not asked recently". Managed by the probe.
- `last_rebooted_at` (String) Device boot time computed from sysUpTime on each walk. Managed by the probe.
- `last_seen_at` (String) When SNMP data was last received from this device.
- `last_snmp_seen_at` (String) When the last SUCCESSFUL SNMP walk of this device completed — the moment its interfaces, inventory and health OIDs were last refreshed. Only moves on a successful walk, so it stays honest while lastSeenAt keeps moving on ping alone. Managed by the probe.
- `last_walk_log` (String) The previous poll's interface counters. Kept so interface rates (bandwidth, utilization, errors/sec) can be computed as counter deltas between polls, and stores nothing else - the rest of the walk response has no reader and this column is rewritten on every poll of every device. Managed by the server. A JSON value: write it with `jsonencode()`.
- `lldp_neighbors` (String) LLDP neighbors discovered on the last SNMP walk, used to build the network topology graph. Managed by the probe. A JSON value: write it with `jsonencode()`.
- `next_poll_at` (String) When the assigned probe should next poll this device. Advanced by the claim cycle; a device created now is due immediately.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `snmp_oids` (String) SNMP OIDs collected on each poll for this device ALONE, on top of whatever its OID Collection Template collects. Values are recorded as metrics and can be alerted on through monitor criteria. If several devices need the same OID, put it on a template instead. A JSON value: write it with `jsonencode()`.
- `snmp_table_snapshot` (String) The rows of every SNMP table collected on the last successful walk - tunnels, radios, neighbours and so on - with their values. Managed by the probe. A JSON value: write it with `jsonencode()`.
- `snmp_tables` (String) SNMP tables walked on each poll for this device alone, on top of its OID Collection Template's tables. A table with the same key as a template table replaces it on this device. A JSON value: write it with `jsonencode()`.
- `snmp_v3_auth` (String) Deprecated: SNMP v3 auth is now stored in the snmpV3* columns below. Retained for reading legacy devices. A JSON value: write it with `jsonencode()`.
- `transceiver_snapshot` (String) The transceivers (SFP, SFP+, QSFP and similar optics) in this device's ports: who made each one, its temperature, supply voltage, bias current and transmit and receive power against the device's own warning and alarm thresholds, its health, whether it is still detected, and a month of daily received power averages. Managed by the probe. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
