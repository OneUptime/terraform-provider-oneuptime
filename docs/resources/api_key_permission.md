---
page_title: "oneuptime_api_key_permission Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Permissions for your API Keys
---

# oneuptime_api_key_permission (Resource)

Permissions for your API Keys

## Example Usage

```terraform
resource "oneuptime_api_key_permission" "example" {
  api_key_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `api_key_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `permission` (String) Permission. You can find list of permissions on the Permissions page...
- `labels` (Set) Relation to Labels Array where this permission is scoped at...
- `is_block_permission` (Bool) Api key permission is_block_permission.

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_api_key_permission.example <id>
```
