---
page_title: "oneuptime_network_device_auto_import_rule Resource - oneuptime"
subcategory: "Network"
description: |-
  Automatically import matching hosts from network device discovery scan results as Network Devices and optionally provision a monitor from a template
---

# oneuptime_network_device_auto_import_rule (Resource)

Automatically import matching hosts from network device discovery scan results as Network Devices and optionally provision a monitor from a template

## Example Usage

```terraform
resource "oneuptime_network_device_auto_import_rule" "example" {
  name        = "Example network device auto import rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this network device auto import rule.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this network device auto import rule.
- `include_ping_only_hosts` (Boolean) Also import hosts that answered ping but not SNMP. Off by default: a wrong SNMP credential makes every host on a subnet report as ping-only, and this rule would then import all of them as half-identified devices. Defaults to `false`.
- `ip_match_target` (String) Only trigger for discovered hosts whose IP is inside this CIDR (192.168.1.0/24) or octet range (10.16-22.0-255.51-66) — the same notations a scan target takes. Leave empty to match any address.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `is_exclusion` (Boolean) Invert this rule: matching hosts are NEVER auto-imported, even when another rule matches them. Use it to carve printers, phones, or other unwanted hosts out of a broader rule. Defaults to `false`.
- `monitor_template_id` (String) ID of the optional monitor template (Network Device, Ping, IP, Port or NTP) each device this rule imports gets a monitor from. The ID of a `oneuptime_monitor_template`.
- `oid_template_id` (String) ID of the optional OID Collection Template to link to devices imported by this rule. The ID of a `oneuptime_oid_collection_template`.
- `polling_interval_in_minutes` (Number) How often the scan's probe polls each device this rule imports, in minutes: 1, 2, 5, 10, 15, 30 or 60. The probe pings every device and walks the ones that answered SNMP. Leave it empty for the devices' default, every 5 minutes. The same choice as Check every when importing from Review Results.
- `sys_descr_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysDescr. Leave empty to match any description.
- `sys_name_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the discovered host's SNMP sysName. Leave empty to match any name.
- `sys_object_id_pattern` (String) An OID prefix (1.3.6.1.4.1.9) or a '*' wildcard OID pattern with literal dots (1.3.6.1.4.1.9.* for Cisco) matched against the discovered host's SNMP sysObjectID — the vendor's registered enterprise OID. Not regex: dots match dots, so 1.3.6.1.4.1.9.* can never match enterprise 94. Leave empty to match any vendor. Only hosts reported by probes new enough to carry sysObjectID can match.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network device auto import rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_device_auto_import_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_device_auto_import_rule.example <id>
```
