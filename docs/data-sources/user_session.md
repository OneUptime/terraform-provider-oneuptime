---
page_title: "oneuptime_user_session Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Active user sessions with refresh tokens and device metadata for enhanced authentication security.
---

# oneuptime_user_session (Data Source)

Active user sessions with refresh tokens and device metadata for enhanced authentication security. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_user_session" "by_name" {
  name = "example-user_session"
}

data "oneuptime_user_session" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
