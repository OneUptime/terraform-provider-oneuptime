---
page_title: "oneuptime_network_device_team_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your network devices.
---

# oneuptime_network_device_team_owner (Resource)

Add teams as owners to your network devices.

## Example Usage

```terraform
resource "oneuptime_network_device_team_owner" "example" {
  team_id = "123e4567-e89b-12d3-a456-426614174000"
  network_device_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `team_id` (String) A unique identifier for an object, represented as a UUID..
- `network_device_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_owner_notified` (Bool) Are owners notified of this resource ownership?..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_network_device_team_owner.example <id>
```
