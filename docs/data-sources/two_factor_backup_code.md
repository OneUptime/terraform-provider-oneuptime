---
page_title: "oneuptime_two_factor_backup_code Data Source - oneuptime"
subcategory: "Other"
description: |-
  Single-use backup codes that let a user sign in when their two factor authentication device is unavailable
---

# oneuptime_two_factor_backup_code (Data Source)

Single-use backup codes that let a user sign in when their two factor authentication device is unavailable

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one two factor backup code may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_two_factor_backup_code" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_two_factor_backup_code" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `user_id` (String) User ID who owns this backup code. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `updated_at` (String) Date and Time when the object was updated.
- `used_at` (String) When this backup code was used to sign in. Null while the code is still unused.
