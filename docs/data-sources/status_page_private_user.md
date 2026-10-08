---
page_title: "oneuptime_status_page_private_user Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage private users on your status page
---

# oneuptime_status_page_private_user (Data Source)

Manage private users on your status page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page private user may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_private_user" "example" {
  status_page_id = oneuptime_status_page.example.id
}

# Or by id:
data "oneuptime_status_page_private_user" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_sso_user` (Boolean) Did this user sign up via SSO?
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `email` (String) Email.
- `password` (String, Sensitive) Password.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
