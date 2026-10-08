---
page_title: "oneuptime_totp_auth Data Source - oneuptime"
subcategory: "Other"
description: |-
  TOTP Authentication for users
---

# oneuptime_totp_auth (Data Source)

TOTP Authentication for users

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one totp auth may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_totp_auth" "example" {
  name = "Example totp auth"
}

# Or by id:
data "oneuptime_totp_auth" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_verified` (Boolean) Is this TOTP authentication verified and validated (has user entered the token to verify it).
- `name` (String) Name of the TOTP authentication.
- `two_factor_otp_url` (String) OTP URL of the TOTP authentication.
- `user_id` (String) User ID who deleted this object (if this object was deleted by a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `updated_at` (String) Date and Time when the object was updated.
