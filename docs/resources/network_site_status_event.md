---
page_title: "oneuptime_network_site_status_event Resource - oneuptime"
subcategory: "Other"
description: |-
  History of the rolled-up health status of a Network Site (Operational to Offline for example)
---

# oneuptime_network_site_status_event (Resource)

History of the rolled-up health status of a Network Site (Operational to Offline for example)

## Example Usage

```terraform
resource "oneuptime_network_site_status_event" "example" {
  site_id = "123e4567-e89b-12d3-a456-426614174000"
  monitor_status_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `site_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_status_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `starts_at` (String) A date time object..
- `ends_at` (String) A date time object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_site_status_event.example <id>
```
