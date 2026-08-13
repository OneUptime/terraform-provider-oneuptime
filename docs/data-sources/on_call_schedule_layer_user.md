---
page_title: "oneuptime_on_call_schedule_layer_user Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  On-Call Schedule Layer Users
---

# oneuptime_on_call_schedule_layer_user (Data Source)

On-Call Schedule Layer Users Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_on_call_schedule_layer_user" "by_name" {
  name = "example-on_call_schedule_layer_user"
}

data "oneuptime_on_call_schedule_layer_user" "by_id" {
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
- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_schedule_layer_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `order` (Number) Order / Priority of this layer. Lower the number, higher the priority... Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
