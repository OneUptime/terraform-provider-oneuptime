---
page_title: "oneuptime_network_site_assignment_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Rules that automatically assign discovered Network Devices and Endpoints to a Network Site by subnet CIDR or hostname pattern. At least one of the two matchers must be set.
---

# oneuptime_network_site_assignment_rule (Data Source)

Rules that automatically assign discovered Network Devices and Endpoints to a Network Site by subnet CIDR or hostname pattern. At least one of the two matchers must be set. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_site_assignment_rule" "by_name" {
  name = "example-network_site_assignment_rule"
}

data "oneuptime_network_site_assignment_rule" "by_id" {
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
- `site_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `subnet_cidr` (String) Devices and endpoints with an IP in this CIDR are assigned to the site.. Computed.
- `hostname_pattern` (String) Devices whose hostname, SNMP system name or display name matches this wildcard pattern are assigned to the site.. Computed.
- `priority` (Number) Higher priority number wins; ties broken by earlier creation... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
