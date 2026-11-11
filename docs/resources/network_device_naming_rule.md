---
page_title: "oneuptime_network_device_naming_rule Resource - oneuptime"
subcategory: "Network"
description: |-
  Rules that name network devices by your naming convention: which devices, and a name pattern such as UN{site.code}{lookup}. The first enabled rule that matches a device names it when discovery imports it, and the rules can be applied to the devices you already have.
---

# oneuptime_network_device_naming_rule (Resource)

Rules that name network devices by your naming convention: which devices, and a name pattern such as UN{site.code}{lookup}. The first enabled rule that matches a device names it when discovery imports it, and the rules can be applied to the devices you already have.

## Example Usage

```terraform
resource "oneuptime_network_device_naming_rule" "example" {
  name         = "Example network device naming rule"
  name_pattern = "UN{site.code}{lookup}"
  description  = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) A name for this rule, such as Menu boards.
- `name_pattern` (String) What matched devices are called. Text, with placeholders in curly braces filled in from each device: {site.code}, {site}, {discoveredName}, {ip}, {ip.1} to {ip.4}, {mac}, {vendor} and {lookup}, the value of the rule's lookup table. {n}, {nn}, {nnn} or {nnnn} is the lowest number that makes the name unique. A placeholder can be followed by filters: {vendor|first:3|upper}, with upper, lower, first:N, last:N, digits and nospaces.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) What this rule is for.
- `discovered_name_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the name discovery found for the device. Leave empty to match any name.
- `ip_match_target` (String) Only devices whose IP address is in this range: a CIDR such as 10.48.145.0/24, or octet ranges such as 10.0-255.0-255.150-159, where * is any number (10.*.*.51). Leave empty to match any address.
- `is_enabled` (Boolean) Whether this rule names devices. A disabled rule is skipped and the next one is tried. Defaults to `true`.
- `lookup_entries` (String) Rows of a key and the value it stands for, put where the name pattern says {lookup}. A device whose key has no row is left to the next rule. A JSON value: write it with `jsonencode()`.
- `lookup_key` (String) What the rows of the lookup table are keyed by: lastIpOctet (the last part of the device's IP address), ipAddress, macAddress or discoveredName.
- `mac_address_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against the device's MAC address, in any spelling: 00:05:0b*, 00-05-0B* and 00050b* match the same devices. Leave empty to match any MAC address.
- `priority` (Number) Where this rule sits in the list: the first enabled rule from the top that matches a device, and has every value its name pattern needs, names it. The rule at the top has the highest number. A new rule is added to the end of the list, with the lowest number. Setting a number another rule already has puts this rule in that place, and the rules in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `sites` (Set of String) Only devices in one of these sites, or in a site below one of them. Leave empty to match devices in any site, or in none. IDs of `oneuptime_network_site` resources.
- `vendor_pattern` (String) Regex or * wildcard pattern (case-insensitive) matched against who made the device. Leave empty to match any vendor.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network device naming rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_device_naming_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_device_naming_rule.example <id>
```
