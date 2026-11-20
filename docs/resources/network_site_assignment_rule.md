---
page_title: "oneuptime_network_site_assignment_rule Resource - oneuptime"
subcategory: "Network"
description: |-
  Rules that automatically assign discovered Network Devices and Endpoints to a Network Site by subnet CIDR or hostname pattern. At least one of the two matchers must be set.
---

# oneuptime_network_site_assignment_rule (Resource)

Rules that automatically assign discovered Network Devices and Endpoints to a Network Site by subnet CIDR or hostname pattern. At least one of the two matchers must be set.

## Example Usage

```terraform
resource "oneuptime_network_site_assignment_rule" "example" {
  site_id = oneuptime_network_site.example.id
}
```

<a id="order"></a>
## Order

`priority` is this network site assignment rule's place among the network site assignment rules of its project, the highest number first. Setting it places this network site assignment rule and moves no other, so a configuration that numbers every network site assignment rule of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.network_site_assignment_rules`:

```terraform
priority = length(local.network_site_assignment_rules) - index(local.network_site_assignment_rules, each.key)
```

- Give each network site assignment rule a number of its own: two network site assignment rules at one number are shown the older first.
- A network site assignment rule created without `priority` goes below the last one of its project, one number lower.
- Deleting a network site assignment rule leaves a gap in the numbers, which changes nothing: network site assignment rules go by their numbers, not by counting.
- Moving a network site assignment rule in the OneUptime Dashboard numbers the project's network site assignment rules n..1 again; the next plan shows that as changes to `priority`.


## Schema

### Required

- `site_id` (String) ID of the Network Site matched resources are assigned to. The ID of a `oneuptime_network_site`.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `hostname_pattern` (String) Devices whose hostname, SNMP system name, display name or DNS name matches this wildcard pattern are assigned to the site.
- `priority` (Number) Where this rule sits in the list: when several rules match a device, the one highest in the list wins. The rule at the top has the highest number. A new rule is added to the end of the list, with the lowest number. Setting a number places this rule there and moves no other rule: two rules at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).
- `subnet_cidr` (String) Devices and endpoints with an IP in this CIDR are assigned to the site.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network site assignment rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_site_assignment_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_site_assignment_rule.example <id>
```
