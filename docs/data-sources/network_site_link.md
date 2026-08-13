---
page_title: "oneuptime_network_site_link Data Source - oneuptime"
subcategory: "Other"
description: |-
  Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor.
---

# oneuptime_network_site_link (Data Source)

Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_network_site_link" "by_name" {
  name = "example-network_site_link"
}

data "oneuptime_network_site_link" "by_id" {
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
- `from_site_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `to_site_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
