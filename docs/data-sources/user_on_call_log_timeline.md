---
page_title: "oneuptime_user_on_call_log_timeline Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  Timeline events for user on-call log.
---

# oneuptime_user_on_call_log_timeline (Data Source)

Timeline events for user on-call log. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_user_on_call_log_timeline" "by_name" {
  name = "example-user_on_call_log_timeline"
}

data "oneuptime_user_on_call_log_timeline" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
