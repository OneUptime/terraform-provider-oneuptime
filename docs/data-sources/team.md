---
page_title: "oneuptime_team Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Teams lets your organize users of your project into groups and lets you assign different level of permissions.
---

# oneuptime_team (Data Source)

Teams lets your organize users of your project into groups and lets you assign different level of permissions.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one team may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_team" "example" {
  name = "Example team"
}

# Or by id:
data "oneuptime_team" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_permissions_editable` (Boolean) Can you edit team permissions? Teams auto-created for you are uneditable but you should be able to edit permissions on the team you create.
- `is_team_deleteable` (Boolean) Can you delete this team? Teams auto-created for you are not deleteable but you should be able to delete permissions on the team you create.
- `is_team_editable` (Boolean) Can you edit team? Teams auto-created for you are uneditable but you should be able to edit on the team you create.
- `name` (String) Any friendly name of this object.
- `should_have_at_least_one_member` (Boolean) Can this team have no members? Owner team should have at least 1 member, other teams can have no members.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
