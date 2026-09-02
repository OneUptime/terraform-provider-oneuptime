---
page_title: "oneuptime_network_device_auto_import_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Automatically import matching hosts from network device discovery scan results as Network Devices and optionally provision a monitor from a template
---

# oneuptime_network_device_auto_import_rule (Data Source)

Automatically import matching hosts from network device discovery scan results as Network Devices and optionally provision a monitor from a template Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_device_auto_import_rule" "by_name" {
  name = "example-network_device_auto_import_rule"
}

data "oneuptime_network_device_auto_import_rule" "by_id" {
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
- `description` (String) Description of this network device auto import rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `ip_match_target` (String) Only trigger for discovered hosts whose IP is inside this CIDR (192.168.1.0/24) or octet range (10.16-22.0-255.51-66) — the same notations a scan target takes. Leave empty to match any address... Computed.
- `sys_name_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysName. Leave empty to match any name... Computed.
- `sys_descr_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysDescr. Leave empty to match any description... Computed.
- `sys_object_id_pattern` (String) An OID prefix (1.3.6.1.4.1.9) or a '*' wildcard OID pattern with literal dots (1.3.6.1.4.1.9.* for Cisco) matched against the discovered host's SNMP sysObjectID — the vendor's registered enterprise OID. Not regex: dots match dots, so 1.3.6.1.4.1.9.* can never match enterprise 94. Leave empty to match any vendor. Only hosts reported by probes new enough to carry sysObjectID can match... Computed.
- `include_ping_only_hosts` (Bool) Also import hosts that answered ping but not SNMP. Off by default: a wrong SNMP credential makes every host on a subnet report as ping-only, and this rule would then import all of them as half-identified devices... Computed.
- `is_exclusion` (Bool) Invert this rule: matching hosts are NEVER auto-imported, even when another rule matches them. Use it to carve printers, phones, or other unwanted hosts out of a broader rule... Computed.
- `monitor_template_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `oid_template_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
