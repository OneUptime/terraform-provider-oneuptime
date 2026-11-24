---
page_title: "oneuptime_dashboard_private_user Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  People who can sign in to open a dashboard's public link when it asks its visitors to sign in.
---

# oneuptime_dashboard_private_user (Data Source)

People who can sign in to open a dashboard's public link when it asks its visitors to sign in.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one dashboard private user may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_dashboard_private_user" "example" {
  dashboard_id = oneuptime_dashboard.example.id
}

# Or by id:
data "oneuptime_dashboard_private_user" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `dashboard_id` (String) ID of the dashboard this private user may sign in to. The ID of a `oneuptime_dashboard`.
- `has_password` (Boolean) Whether this private user has set a password.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_sso_user` (Boolean) Whether the dashboard's single sign-on added this user.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `email` (String) The address this private user signs in with.
- `last_active` (String) When this private user last signed in to the dashboard.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
