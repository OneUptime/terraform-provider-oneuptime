---
page_title: "oneuptime_on_call_schedule_layer Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  On-Call Schedule Layers
---

# oneuptime_on_call_schedule_layer (Resource)

On-Call Schedule Layers

## Example Usage

```terraform
resource "oneuptime_on_call_schedule_layer" "example" {
  on_call_duty_policy_schedule_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  starts_at = "2030-01-01T00:00:00Z"
  hand_off_time = "2030-01-01T00:00:00Z"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Friendly name for this layer..
- `starts_at` (String) A date time object..
- `hand_off_time` (String) A date time object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description for this layer. This is optional and can be left blank...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `order` (Number) Order / Priority of this layer. Lower the number, higher the priority...
- `rotation` (String) How often would you like to hand off the duty to the next user in this layer?..
- `restriction_times` (String) Restrict this layer to these times..

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
terraform import oneuptime_on_call_schedule_layer.example <id>
```
