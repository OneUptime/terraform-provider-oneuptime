---
page_title: "oneuptime_team_permission Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Permissions for your OneUptime team
---

# oneuptime_team_permission (Resource)

Permissions for your OneUptime team

## Example Usage

```terraform
resource "oneuptime_team_permission" "example" {

}
```

## Schema

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `team_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `permission` (String) Permission. You can find list of permissions on the Permissions page...
- `labels` (Set) Relation to Labels Array where this permission is scoped at...
- `is_block_permission` (Bool) Team permission is_block_permission.
- `scope` (String) Scope of this permission row. One of: All, Owned, Labels. Defaults to All so new permissions apply to every resource in the project unless explicitly narrowed...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_team_permission.example <id>
```
