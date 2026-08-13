---
page_title: "oneuptime_monitor_team_owner Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Add teams as owners to your monitors.
---

# oneuptime_monitor_team_owner (Resource)

Add teams as owners to your monitors.

## Example Usage

```terraform
resource "oneuptime_monitor_team_owner" "example" {
  team_id = "123e4567-e89b-12d3-a456-426614174000"
  monitor_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `team_id` (String) A unique identifier for an object, represented as a UUID..
- `monitor_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_owner_notified` (Bool) Are owners notified of this resource ownership?..

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
terraform import oneuptime_monitor_team_owner.example <id>
```
