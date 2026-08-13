---
page_title: "oneuptime_on_call_schedule_layer_user Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  On-Call Schedule Layer Users
---

# oneuptime_on_call_schedule_layer_user (Resource)

On-Call Schedule Layer Users

## Example Usage

```terraform
resource "oneuptime_on_call_schedule_layer_user" "example" {
  on_call_duty_policy_schedule_id = "123e4567-e89b-12d3-a456-426614174000"
  on_call_duty_policy_schedule_layer_id = "123e4567-e89b-12d3-a456-426614174000"
  user_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policy_schedule_layer_id` (String) A unique identifier for an object, represented as a UUID..
- `user_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `order` (Number) Order / Priority of this layer. Lower the number, higher the priority...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_on_call_schedule_layer_user.example <id>
```
