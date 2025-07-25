---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_team_data Data Source - oneuptime"
subcategory: ""
description: |-
  Team data data source
---

# oneuptime_team_data (Data Source)

Team data data source

## Example Usage

```terraform
data "oneuptime_team_data" "example" {
  name = "example-team_data"
}
```

## Schema

- `id` (String) Identifier to filter by. Optional.
- `name` (String) Name to filter by. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Permissions - Create: [Project Owner, Project Admin, Create Team], Read: [Project Owner, Project Admin, Project Member, Read Teams], Update: [Project Owner, Project Admin, Edit Team]. Computed.
- `slug` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Teams], Update: [No access - you don't have permission for this operation]. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_permissions_editable` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Edit Team, Edit Team Permissions], Update: [No access - you don't have permission for this operation]. Computed.
- `is_team_deleteable` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Edit Team, Edit Team Permissions], Update: [No access - you don't have permission for this operation]. Computed.
- `should_have_at_least_one_member` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Edit Team, Edit Team Permissions], Update: [No access - you don't have permission for this operation]. Computed.
- `is_team_editable` (Bool) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Edit Team, Edit Team Permissions], Update: [No access - you don't have permission for this operation]. Computed.
