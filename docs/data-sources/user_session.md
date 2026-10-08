---
page_title: "oneuptime_user_session Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Active user sessions with refresh tokens and device metadata for enhanced authentication security.
---

# oneuptime_user_session (Data Source)

Active user sessions with refresh tokens and device metadata for enhanced authentication security.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one user session may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_user_session" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_user_session" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `device_browser` (String) Browser or client application used for this session.
- `device_name` (String) Friendly name for the device used to sign in.
- `device_os` (String) Operating system reported for this session.
- `device_type` (String) Type of device (e.g., desktop, mobile).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `ip_address` (String) IP address observed for this session.
- `is_revoked` (Boolean) Marks whether the session has been explicitly revoked.
- `revoked_reason` (String) Optional reason describing why the session was revoked.
- `user_agent` (String) Complete user agent string supplied by the client.
- `user_id` (String) Identifier for the user that owns this session. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `additional_info` (String) Flexible JSON payload for storing structured session metadata. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_active_at` (String) Last time this session was used.
- `revoked_at` (String) Timestamp when the session was revoked, if applicable.
- `updated_at` (String) Date and Time when the object was updated.
