---
page_title: "oneuptime_two_factor_backup_code Data Source - oneuptime"
subcategory: "Other"
description: |-
  Single-use backup codes that let a user sign in when their two factor authentication device is unavailable
---

# oneuptime_two_factor_backup_code (Data Source)

Single-use backup codes that let a user sign in when their two factor authentication device is unavailable Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_two_factor_backup_code" "by_name" {
  name = "example-two_factor_backup_code"
}

data "oneuptime_two_factor_backup_code" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
