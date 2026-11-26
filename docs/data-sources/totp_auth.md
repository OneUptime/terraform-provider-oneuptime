---
page_title: "oneuptime_totp_auth Data Source - oneuptime"
subcategory: "Teams & Access"
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

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_verified` (Boolean) Is this TOTP authentication verified and validated (has user entered the token to verify it).
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Name of the TOTP authentication.
- `two_factor_otp_url` (String) OTP URL of the TOTP authentication.
- `user_id` (String) User ID who deleted this object (if this object was deleted by a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `updated_at` (String) Date and Time when the object was updated.
