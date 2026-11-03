---
page_title: "oneuptime_on_call_schedule_layer_user Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add users to the layers of your on-call schedules.
---

# oneuptime_on_call_schedule_layer_user (Resource)

Add users to the layers of your on-call schedules.

## Example Usage

```terraform
resource "oneuptime_on_call_schedule_layer_user" "example" {
  on_call_duty_policy_schedule_id       = oneuptime_on_call_policy_schedule.example.id
  on_call_duty_policy_schedule_layer_id = oneuptime_on_call_schedule_layer.example.id
  user_email                            = "alice@example.com"
}
```

## Schema

### Required

- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule the layer is in. The ID of a `oneuptime_on_call_policy_schedule`.
- `on_call_duty_policy_schedule_layer_id` (String) ID of the layer the person takes turns on call in. The ID of a `oneuptime_on_call_schedule_layer`.

### Optional

- `order` (Number) Order / Priority of this layer. Lower the number, higher the priority.
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the person who takes turns on call in this layer. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call schedule layer user by its on call schedule layer and the person's email address (or user ID), with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_schedule_layer_user.example
  id = "<on_call_duty_policy_schedule_layer_id>/alice@example.com"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_schedule_layer_user.example <on_call_duty_policy_schedule_layer_id>/alice@example.com
```

Imported by an address, the person's `user_email` stays as you wrote it. Its own ID works too.
