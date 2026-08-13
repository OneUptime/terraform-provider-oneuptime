---
page_title: "oneuptime_team_member Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  This model connects users and teams
---

# oneuptime_team_member (Resource)

This model connects users and teams

## Example Usage

```terraform
resource "oneuptime_team_member" "example" {
  user_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `user_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `team_id` (String) A unique identifier for an object, represented as a UUID..
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `has_accepted_invitation` (Bool) Has this team member accepted invitation..
- `invitation_accepted_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_team_member.example <id>
```
