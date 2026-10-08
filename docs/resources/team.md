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
  name        = "Example team"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `custom_fields` (String) Custom Fields on this resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Friendly description that will help you remember.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_permissions_editable` (Boolean) Can you edit team permissions? Teams auto-created for you are uneditable but you should be able to edit permissions on the team you create.
- `is_team_deleteable` (Boolean) Can you delete this team? Teams auto-created for you are not deleteable but you should be able to delete permissions on the team you create.
- `is_team_editable` (Boolean) Can you edit team? Teams auto-created for you are uneditable but you should be able to edit on the team you create.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `should_have_at_least_one_member` (Boolean) Can this team have no members? Owner team should have at least 1 member, other teams can have no members.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing team by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_team.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_team.example <id>
```
