---
page_title: "oneuptime_on_call_schedule_layer_user Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add users to the layers of your on-call schedules.
---

# oneuptime_on_call_schedule_layer_user (Data Source)

Add users to the layers of your on-call schedules.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call schedule layer user may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_schedule_layer_user" "example" {
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.example.id
  user_email                      = "alice@example.com"
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
- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule the layer is in. The ID of a `oneuptime_on_call_policy_schedule`.
- `on_call_duty_policy_schedule_layer_id` (String) ID of the layer the person takes turns on call in. The ID of a `oneuptime_on_call_schedule_layer`.
- `order` (Number) Order / Priority of this layer. Lower the number, higher the priority.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the person who takes turns on call in this layer. The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
