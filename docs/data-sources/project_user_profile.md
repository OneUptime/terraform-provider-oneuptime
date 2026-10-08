---
page_title: "oneuptime_project_user_profile Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Stores user profile data including custom fields for each user in a project.
---

# oneuptime_project_user_profile (Data Source)

Stores user profile data including custom fields for each user in a project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one project user profile may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_project_user_profile" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_project_user_profile" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `user_id` (String) ID of User this profile belongs to. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) Custom Fields for this user in this project. A JSON value: write it with `jsonencode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
