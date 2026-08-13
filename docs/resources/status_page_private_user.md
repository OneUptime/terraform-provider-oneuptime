---
page_title: "oneuptime_status_page_private_user Resource - oneuptime"
subcategory: "Status Pages"
description: |-
   Manage private users on your status page
---

# oneuptime_status_page_private_user (Resource)

 Manage private users on your status page

## Example Usage

```terraform
resource "oneuptime_status_page_private_user" "example" {
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `email` (String) Email object.
- `password` (String, Sensitive) Password..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_sso_user` (Bool) Did this user sign up via SSO?..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page_private_user.example <id>
```
