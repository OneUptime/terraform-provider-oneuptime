---
page_title: "oneuptime_user_notification_log Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Log events for user notifications
---

# oneuptime_user_notification_log (Data Source)

Log events for user notifications Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_user_notification_log" "by_name" {
  name = "example-user_notification_log"
}

data "oneuptime_user_notification_log" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
