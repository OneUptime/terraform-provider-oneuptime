---
page_title: "oneuptime_storage_array_team_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add teams as owners to your storage arrays.
---

# oneuptime_storage_array_team_owner (Resource)

Add teams as owners to your storage arrays.

## Example Usage

```terraform
resource "oneuptime_storage_array_team_owner" "example" {
  team_id = "123e4567-e89b-12d3-a456-426614174000"
  storage_array_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `team_id` (String) A unique identifier for an object, represented as a UUID..
- `storage_array_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_owner_notified` (Bool) Are owners notified of this resource ownership?..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_storage_array_team_owner.example <id>
```
