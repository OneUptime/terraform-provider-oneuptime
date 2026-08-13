---
page_title: "oneuptime_team_member Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  This model connects users and teams
---

# oneuptime_team_member (Data Source)

This model connects users and teams Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_team_member" "by_name" {
  name = "example-team_member"
}

data "oneuptime_team_member" "by_id" {
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
- `team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `has_accepted_invitation` (Bool) Has this team member accepted invitation.. Computed.
- `invitation_accepted_at` (String) A date time object.. Computed.
