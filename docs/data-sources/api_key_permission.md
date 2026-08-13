---
page_title: "oneuptime_api_key_permission Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Permissions for your API Keys
---

# oneuptime_api_key_permission (Data Source)

Permissions for your API Keys Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_api_key_permission" "by_name" {
  name = "example-api_key_permission"
}

data "oneuptime_api_key_permission" "by_id" {
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
- `api_key_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `permission` (String) Permission. You can find list of permissions on the Permissions page... Computed.
- `labels` (Set) Relation to Labels Array where this permission is scoped at... Computed.
- `is_block_permission` (Bool) Api key permission is_block_permission. Computed.
