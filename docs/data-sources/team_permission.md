---
page_title: "oneuptime_team_permission Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Permissions for your OneUptime team
---

# oneuptime_team_permission (Data Source)

Permissions for your OneUptime team Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_team_permission" "by_name" {
  name = "example-team_permission"
}

data "oneuptime_team_permission" "by_id" {
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
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `permission` (String) Permission. You can find list of permissions on the Permissions page... Computed.
- `labels` (Set) Relation to Labels Array where this permission is scoped at... Computed.
- `is_block_permission` (Bool) Team permission is_block_permission. Computed.
- `scope` (String) Scope of this permission row. One of: All, Owned, Labels. Defaults to All so new permissions apply to every resource in the project unless explicitly narrowed... Computed.
