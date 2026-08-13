---
page_title: "oneuptime_web_authn_credential Data Source - oneuptime"
subcategory: "Other"
description: |-
  WebAuthn credentials for users (security keys)
---

# oneuptime_web_authn_credential (Data Source)

WebAuthn credentials for users (security keys) Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_web_authn_credential" "by_name" {
  name = "example-web_authn_credential"
}

data "oneuptime_web_authn_credential" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
