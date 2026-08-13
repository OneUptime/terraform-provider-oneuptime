---
page_title: "oneuptime_network_site_link Resource - oneuptime"
subcategory: "Other"
description: |-
  Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor.
---

# oneuptime_network_site_link (Resource)

Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor.

## Example Usage

```terraform
resource "oneuptime_network_site_link" "example" {
  from_site_id = "123e4567-e89b-12d3-a456-426614174000"
  to_site_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
}
```

## Schema

### Required

- `from_site_id` (String) A unique identifier for an object, represented as a UUID..
- `to_site_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Friendly name for this link..
- `monitor_id` (String) A unique identifier for an object, represented as a UUID..
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
terraform import oneuptime_network_site_link.example <id>
```
