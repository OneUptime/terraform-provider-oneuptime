---
page_title: "oneuptime_docker_host_user_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add users as owners to your Docker hosts.
---

# oneuptime_docker_host_user_owner (Resource)

Add users as owners to your Docker hosts.

## Example Usage

```terraform
resource "oneuptime_docker_host_user_owner" "example" {
  user_id = "123e4567-e89b-12d3-a456-426614174000"
  docker_host_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `user_id` (String) A unique identifier for an object, represented as a UUID..
- `docker_host_id` (String) A unique identifier for an object, represented as a UUID..

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
terraform import oneuptime_docker_host_user_owner.example <id>
```
