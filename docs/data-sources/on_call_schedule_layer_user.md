---
page_title: "oneuptime_on_call_schedule_layer_user Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  On-Call Schedule Layer Users
---

# oneuptime_on_call_schedule_layer_user (Data Source)

On-Call Schedule Layer Users

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call schedule layer user may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_schedule_layer_user" "example" {
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.example.id
}

# Or by id:
data "oneuptime_on_call_schedule_layer_user" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `on_call_duty_policy_schedule_id` (String) ID of your On-Call Policy Schedule where this escalation rule belongs. The ID of a `oneuptime_on_call_policy_schedule`.
- `on_call_duty_policy_schedule_layer_id` (String) ID of your On-Call Policy Schedule Layer where this escalation rule belongs. The ID of a `oneuptime_on_call_schedule_layer`.
- `order` (Number) Order / Priority of this layer. Lower the number, higher the priority.
- `user_id` (String) ID of User who belongs to this team. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
