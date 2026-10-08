---
page_title: "oneuptime_on_call_schedule_layer Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  On-Call Schedule Layers
---

# oneuptime_on_call_schedule_layer (Data Source)

On-Call Schedule Layers

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call schedule layer may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_schedule_layer" "example" {
  name = "Example on call schedule layer"
}

# Or by id:
data "oneuptime_on_call_schedule_layer" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description for this layer. This is optional and can be left blank.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Friendly name for this layer.
- `on_call_duty_policy_schedule_id` (String) ID of your On-Call Policy Schedule where this escalation rule belongs. The ID of a `oneuptime_on_call_policy_schedule`.
- `order` (Number) Order / Priority of this layer. Lower the number, higher the priority.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `hand_off_time` (String) Hand off time. When would you like to hand off the duty to the next user in this layer?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `restriction_times` (String) Restrict this layer to these times. A JSON value: write it with `jsonencode()`.
- `rotation` (String) How often would you like to hand off the duty to the next user in this layer? A JSON value: write it with `jsonencode()`.
- `starts_at` (String) Start date and time of this layer.
- `updated_at` (String) Date and Time when the object was updated.
