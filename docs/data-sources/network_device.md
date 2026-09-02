---
page_title: "oneuptime_network_device Data Source - oneuptime"
subcategory: "Other"
description: |-
  Network Devices (routers, switches, firewalls) that are being monitored in this project via SNMP polling and traps.
---

# oneuptime_network_device (Data Source)

Network Devices (routers, switches, firewalls) that are being monitored in this project via SNMP polling and traps. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_device" "by_name" {
  name = "example-network_device"
}

data "oneuptime_network_device" "by_id" {
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
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `description` (String) Friendly description for this network device.. Computed.
- `hostname` (String) IP address or hostname the probe polls; also matches SNMP trap sources.. Computed.
- `probe_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `site_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `oid_template_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `current_monitor_status_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitoring_method` (String) How this device's health is established: SNMP (an assigned probe walks it on a schedule) or Monitor (no polling — the linked monitor's status is the device's status). Devices created before this existed are SNMP... Computed.
- `device_role` (String) What this device does on the network — router, switch, access point and so on. Left empty, the role is worked out from the device's own SNMP identity. Set it when there is no SNMP to read: a ping-only device has no identity to classify, and the role decides both the shape it is drawn with and where it sits in the topology hierarchy... Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `snmp_version` (String) SNMP version to use when polling this device (V1, V2c, V3).. Computed.
- `snmp_community_string` (String) Community string used for SNMP v1/v2c polling.. Computed.
- `snmp_port` (Number) UDP port used for SNMP polling.. Computed.
- `snmp_v3_auth` (String) Deprecated: SNMP v3 auth is now stored in the snmpV3* columns below. Retained for reading legacy devices... Computed.
- `snmp_v3_security_level` (String) SNMP v3 security level: noAuthNoPriv, authNoPriv, or authPriv.. Computed.
- `snmp_v3_username` (String) Security name (username) used for SNMP v3 polling.. Computed.
- `snmp_v3_auth_protocol` (String) SNMP v3 authentication protocol: MD5, SHA, SHA256, or SHA512.. Computed.
- `snmp_v3_auth_key` (String) SNMP v3 authentication passphrase.. Computed.
- `snmp_v3_priv_protocol` (String) SNMP v3 privacy (encryption) protocol: DES, AES, or AES256.. Computed.
- `snmp_v3_priv_key` (String) SNMP v3 privacy (encryption) passphrase.. Computed.
- `is_polling_enabled` (Bool) Whether the assigned probe polls this device on a schedule. Disable to pause SNMP polling without deleting the device... Computed.
- `polling_interval_in_minutes` (Number) How often, in minutes, the assigned probe polls this device via SNMP.. Computed.
- `walk_interfaces` (Bool) Walk the IF-MIB interface tables on each poll to inventory interfaces, bandwidth, and errors. Also collects LLDP/CDP neighbors for the topology graph... Computed.
- `collect_endpoints` (Bool) Also walk the device's ARP cache and bridge forwarding database on each poll to discover endpoints (laptops, printers, POS terminals) attached to it. Strictly opt-in: costs extra SNMP table walks per poll. Only meaningful when Walk Interfaces is on... Computed.
- `snmp_oids` (String) SNMP OIDs collected on each poll for this device ALONE, on top of whatever its OID Collection Template collects. Values are recorded as metrics and can be alerted on through monitor criteria. If several devices need the same OID, put it on a template instead... Computed.
- `auto_apply_vendor_health_template` (Bool) When the device's vendor is fingerprinted from its SNMP sysObjectID and no Health OIDs are configured yet, apply the matching vendor health template automatically on the next poll. Off by default for hand-made devices — the vendor template banner stays the manual path; auto-imported devices enable it so the zero-touch pipeline ends with health metrics, not an empty OID list... Computed.
- `next_poll_at` (String) A date time object.. Computed.
- `last_walk_log` (String) The previous poll's interface counters. Kept so interface rates (bandwidth, utilization, errors/sec) can be computed as counter deltas between polls, and stores nothing else - the rest of the walk response has no reader and this column is rewritten on every poll of every device. Managed by the server... Computed.
- `sys_descr` (String) System description (sysDescr) enriched from SNMP walks of this device.. Computed.
- `sys_name` (String) System name (sysName) enriched from SNMP walks of this device.. Computed.
- `sys_object_id` (String) sysObjectID — the vendor's registered OID for this device model, enriched from SNMP walks. Used to fingerprint the vendor and suggest an OID template... Computed.
- `sys_location` (String) System location (sysLocation) enriched from SNMP walks of this device.. Computed.
- `sys_contact` (String) System contact (sysContact) enriched from SNMP walks of this device.. Computed.
- `vendor` (String) Hardware vendor, from ENTITY-MIB or derived from sysObjectID. Managed by the probe... Computed.
- `device_model` (String) Hardware model from ENTITY-MIB (entPhysicalModelName). Managed by the probe... Computed.
- `serial_number` (String) Chassis serial number from ENTITY-MIB (entPhysicalSerialNum). Managed by the probe... Computed.
- `firmware_version` (String) Firmware revision from ENTITY-MIB (entPhysicalFirmwareRev). Managed by the probe... Computed.
- `software_version` (String) Operating system / software revision from ENTITY-MIB (entPhysicalSoftwareRev). Managed by the probe... Computed.
- `last_rebooted_at` (String) A date time object.. Computed.
- `cdp_neighbors` (String) CDP neighbors discovered on the last SNMP walk, complementing LLDP for the topology graph. Managed by the probe... Computed.
- `lldp_neighbors` (String) LLDP neighbors discovered on the last SNMP walk, used to build the network topology graph. Managed by the probe... Computed.
- `last_seen_at` (String) A date time object.. Computed.
- `last_polled_at` (String) A date time object.. Computed.
- `is_reachable` (Bool) Whether the most recent SNMP walk reached this device. NULL means it has never been polled. This — not the age of lastSeenAt — is what the device list, the topology graph and the site rollup read, so a device whose last poll succeeded is never shown as down just because the probe is behind schedule. Managed by the probe... Computed.
- `interfaces_total` (Number) Cached total count of interfaces on this device.. Computed.
- `interfaces_up` (Number) Cached count of operationally up interfaces on this device.. Computed.
- `interfaces_down` (Number) Cached count of operationally down interfaces on this device.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_archived` (Bool) Is this network device archived? Archived network devices are hidden from lists but keep collecting telemetry... Computed.
- `archived_at` (String) A date time object.. Computed.
- `archived_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
