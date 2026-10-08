---
page_title: "oneuptime_user Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  A signed up or invited OneUptime user.
---

# oneuptime_user (Data Source)

A signed up or invited OneUptime user.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one user may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_user" "example" {
  name = "Example user"
}

# Or by id:
data "oneuptime_user" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `company_name` (String)
- `company_size` (String)
- `enable_two_factor_auth` (Boolean) Is two factor authentication enabled?
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_blocked` (Boolean)
- `is_disabled` (Boolean)
- `is_email_verified` (Boolean)
- `job_role` (String)
- `name` (String)
- `profile_picture_id` (String) The ID of a `oneuptime_file`.
- `referral` (String)
- `timezone` (String)
- `two_factor_auth_enabled` (Boolean)

### Read-Only

- `company_phone_number` (String)
- `created_at` (String) Date and Time when the object was created.
- `email` (String) Email.
- `new_unverified_temporary_email` (String)
- `password` (String, Sensitive) Password.
- `updated_at` (String) Date and Time when the object was updated.
