---
page_title: "oneuptime_user_override Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call duty user overrides, for example if the user is on leave you can override the on-call duty policy for that user so all the alerts will be routed to the other user.
---

# oneuptime_user_override (Data Source)

Manage on-call duty user overrides, for example if the user is on leave you can override the on-call duty policy for that user so all the alerts will be routed to the other user. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_user_override" "by_name" {
  name = "example-user_override"
}

data "oneuptime_user_override" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `override_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `route_alerts_to_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `starts_at` (String) A date time object.. Computed.
- `ends_at` (String) A date time object.. Computed.
