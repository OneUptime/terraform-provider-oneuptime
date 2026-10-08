---
page_title: "oneuptime_network_site_assignment_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Rules that automatically assign discovered Network Devices and Endpoints to a Network Site by subnet CIDR or hostname pattern. At least one of the two matchers must be set.
---

# oneuptime_network_site_assignment_rule (Data Source)

Rules that automatically assign discovered Network Devices and Endpoints to a Network Site by subnet CIDR or hostname pattern. At least one of the two matchers must be set.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network site assignment rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_site_assignment_rule" "example" {
  site_id = oneuptime_network_site.example.id
}

# Or by id:
data "oneuptime_network_site_assignment_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `hostname_pattern` (String) Devices whose hostname, SNMP system name, display name or DNS name matches this wildcard pattern are assigned to the site.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `priority` (Number) Where this rule sits in the list: when several rules match a device, the one highest in the list wins. The rule at the top has the highest number. A new rule is added to the end of the list, with the lowest number. Setting a number another rule already has puts this rule in that place, and the rules in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `site_id` (String) ID of the Network Site matched resources are assigned to. The ID of a `oneuptime_network_site`.
- `subnet_cidr` (String) Devices and endpoints with an IP in this CIDR are assigned to the site.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
