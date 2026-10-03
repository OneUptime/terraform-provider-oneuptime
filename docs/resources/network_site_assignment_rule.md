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
  site_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `site_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `subnet_cidr` (String) Devices and endpoints with an IP in this CIDR are assigned to the site..
- `hostname_pattern` (String) Devices whose hostname, SNMP system name, display name or DNS name matches this wildcard pattern are assigned to the site..
- `priority` (Number) Where this rule sits in the list: when several rules match a device, the one highest in the list wins. The rule at the top has the highest number. A new rule is added to the end of the list, with the lowest number. Setting a number another rule already has puts this rule in that place, and the rules in the way move one place along to make room. In the dashboard, drag the rows to reorder them...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_site_assignment_rule.example <id>
```
