---
page_title: "oneuptime_network_device_auto_import_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Automatically import matching hosts from network device discovery scan results as Network Devices and optionally provision a monitor from a template
---

# oneuptime_network_device_auto_import_rule (Data Source)

Automatically import matching hosts from network device discovery scan results as Network Devices and optionally provision a monitor from a template

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network device auto import rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_device_auto_import_rule" "example" {
  name = "Example network device auto import rule"
}

# Or by id:
data "oneuptime_network_device_auto_import_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this network device auto import rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `include_ping_only_hosts` (Boolean) Also import hosts that answered ping but not SNMP. Off by default: a wrong SNMP credential makes every host on a subnet report as ping-only, and this rule would then import all of them as half-identified devices.
- `ip_match_target` (String) Only trigger for discovered hosts whose IP is inside this CIDR (192.168.1.0/24) or octet range (10.16-22.0-255.51-66) — the same notations a scan target takes. Leave empty to match any address.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `is_exclusion` (Boolean) Invert this rule: matching hosts are NEVER auto-imported, even when another rule matches them. Use it to carve printers, phones, or other unwanted hosts out of a broader rule.
- `monitor_template_id` (String) ID of the optional Network Device monitor template to apply to devices imported by this rule. The ID of a `oneuptime_monitor_template`.
- `name` (String) Name of this network device auto import rule.
- `oid_template_id` (String) ID of the optional OID Collection Template to link to devices imported by this rule. The ID of a `oneuptime_oid_collection_template`.
- `sys_descr_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysDescr. Leave empty to match any description.
- `sys_name_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysName. Leave empty to match any name.
- `sys_object_id_pattern` (String) An OID prefix (1.3.6.1.4.1.9) or a '*' wildcard OID pattern with literal dots (1.3.6.1.4.1.9.* for Cisco) matched against the discovered host's SNMP sysObjectID — the vendor's registered enterprise OID. Not regex: dots match dots, so 1.3.6.1.4.1.9.* can never match enterprise 94. Leave empty to match any vendor. Only hosts reported by probes new enough to carry sysObjectID can match.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
