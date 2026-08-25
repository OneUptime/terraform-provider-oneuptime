---
page_title: "oneuptime_network_device_auto_import_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Automatically import matching hosts from network device discovery scan results as Network Devices, with no manual review step
---

# oneuptime_network_device_auto_import_rule (Resource)

Automatically import matching hosts from network device discovery scan results as Network Devices, with no manual review step

## Example Usage

```terraform
resource "oneuptime_network_device_auto_import_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this network device auto import rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this network device auto import rule..
- `is_enabled` (Bool) Whether this rule is enabled..
- `ip_match_target` (String) Only trigger for discovered hosts whose IP is inside this CIDR (192.168.1.0/24) or octet range (10.16-22.0-255.51-66) — the same notations a scan target takes. Leave empty to match any address...
- `sys_name_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysName. Leave empty to match any name...
- `sys_descr_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysDescr. Leave empty to match any description...
- `sys_object_id_pattern` (String) An OID prefix (1.3.6.1.4.1.9) or a '*' wildcard OID pattern with literal dots (1.3.6.1.4.1.9.* for Cisco) matched against the discovered host's SNMP sysObjectID — the vendor's registered enterprise OID. Not regex: dots match dots, so 1.3.6.1.4.1.9.* can never match enterprise 94. Leave empty to match any vendor. Only hosts reported by probes new enough to carry sysObjectID can match...
- `include_ping_only_hosts` (Bool) Also import hosts that answered ping but not SNMP. Off by default: a wrong SNMP credential makes every host on a subnet report as ping-only, and this rule would then import all of them as half-identified devices...
- `is_exclusion` (Bool) Invert this rule: matching hosts are NEVER auto-imported, even when another rule matches them. Use it to carve printers, phones, or other unwanted hosts out of a broader rule...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_device_auto_import_rule.example <id>
```
