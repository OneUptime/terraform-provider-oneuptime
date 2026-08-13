---
page_title: "oneuptime_team Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Teams lets your organize users of your project into groups and lets you assign different level of permissions.
---

# oneuptime_team (Resource)

Teams lets your organize users of your project into groups and lets you assign different level of permissions.

## Example Usage

```terraform
resource "oneuptime_team" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `custom_fields` (String) Custom Fields on this resource...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `is_permissions_editable` (Bool) Can you edit team permissions? Teams auto-created for you are uneditable but you should be able to edit permissions on the team you create..
- `is_team_deleteable` (Bool) Can you delete this team? Teams auto-created for you are not deleteable but you should be able to delete permissions on the team you create..
- `should_have_at_least_one_member` (Bool) Can this team have no members? Owner team should have at least 1 member, other teams can have no members..
- `is_team_editable` (Bool) Can you edit team? Teams auto-created for you are uneditable but you should be able to edit on the team you create..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_team.example <id>
```
