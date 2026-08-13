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

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `subnet_cidr` (String) Devices and endpoints with an IP in this CIDR are assigned to the site..
- `hostname_pattern` (String) Devices whose hostname, SNMP system name or display name matches this wildcard pattern are assigned to the site..
- `priority` (Number) Higher priority number wins; ties broken by earlier creation...
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
