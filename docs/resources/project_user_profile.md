---
page_title: "oneuptime_project_user_profile Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Stores user profile data including custom fields for each user in a project.
---

# oneuptime_project_user_profile (Resource)

Stores user profile data including custom fields for each user in a project.

## Example Usage

```terraform
resource "oneuptime_project_user_profile" "example" {
  user_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `user_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `custom_fields` (String) Custom Fields for this user in this project...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_project_user_profile.example <id>
```
