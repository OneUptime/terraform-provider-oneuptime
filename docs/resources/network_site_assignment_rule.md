---
page_title: "oneuptime_network_site_assignment_rule Resource - oneuptime"
subcategory: "Other"
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

## Schema

### Required

- `site_id` (String) ID of the Network Site matched resources are assigned to. The ID of a `oneuptime_network_site`.

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `hostname_pattern` (String) Devices whose hostname, SNMP system name, display name or DNS name matches this wildcard pattern are assigned to the site.
- `priority` (Number) Where this rule sits in the list: when several rules match a device, the one highest in the list wins. The rule at the top has the highest number. A new rule is added to the end of the list, with the lowest number. Setting a number another rule already has puts this rule in that place, and the rules in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `subnet_cidr` (String) Devices and endpoints with an IP in this CIDR are assigned to the site.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
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
