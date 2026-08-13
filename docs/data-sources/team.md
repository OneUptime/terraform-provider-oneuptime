---
page_title: "oneuptime_team Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Teams lets your organize users of your project into groups and lets you assign different level of permissions.
---

# oneuptime_team (Data Source)

Teams lets your organize users of your project into groups and lets you assign different level of permissions. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_team" "by_name" {
  name = "example-team"
}

data "oneuptime_team" "by_id" {
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
- `description` (String) Friendly description that will help you remember.. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_permissions_editable` (Bool) Can you edit team permissions? Teams auto-created for you are uneditable but you should be able to edit permissions on the team you create.. Computed.
- `is_team_deleteable` (Bool) Can you delete this team? Teams auto-created for you are not deleteable but you should be able to delete permissions on the team you create.. Computed.
- `should_have_at_least_one_member` (Bool) Can this team have no members? Owner team should have at least 1 member, other teams can have no members.. Computed.
- `is_team_editable` (Bool) Can you edit team? Teams auto-created for you are uneditable but you should be able to edit on the team you create.. Computed.
- `custom_fields` (String) Custom Fields on this resource... Computed.
