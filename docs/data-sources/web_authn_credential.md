---
page_title: "oneuptime_web_authn_credential Data Source - oneuptime"
subcategory: "Other"
description: |-
  WebAuthn credentials for users (security keys)
---

# oneuptime_web_authn_credential (Data Source)

WebAuthn credentials for users (security keys)

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one web authn credential may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_web_authn_credential" "example" {
  name = "Example web authn credential"
}

# Or by id:
data "oneuptime_web_authn_credential" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_passkey` (Boolean) Whether this credential was registered for passwordless sign-in. Older credentials have no recorded purpose.
- `is_verified` (Boolean) Is this WebAuthn credential verified and validated.
- `name` (String) Name of the WebAuthn credential.
- `user_id` (String) User ID who owns this WebAuthn credential. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `updated_at` (String) Date and Time when the object was updated.
