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
  created_by_user_id = data.oneuptime_user.example.id
  user_email         = "alice@example.com"
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
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of User this profile belongs to. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `custom_fields` (String) Custom Fields for this user in this project. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
