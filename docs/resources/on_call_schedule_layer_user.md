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
  on_call_duty_policy_schedule_id       = oneuptime_on_call_policy_schedule.example.id
  on_call_duty_policy_schedule_layer_id = oneuptime_on_call_schedule_layer.example.id
  user_id                               = data.oneuptime_user.example.id
}
```

## Schema

### Required

- `on_call_duty_policy_schedule_id` (String) ID of your On-Call Policy Schedule where this escalation rule belongs. The ID of a `oneuptime_on_call_policy_schedule`.
- `on_call_duty_policy_schedule_layer_id` (String) ID of your On-Call Policy Schedule Layer where this escalation rule belongs. The ID of a `oneuptime_on_call_schedule_layer`.
- `user_id` (String) ID of User who belongs to this team. The ID of a `oneuptime_user` (see the data source).

### Optional

- `order` (Number) Order / Priority of this layer. Lower the number, higher the priority.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call schedule layer user by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_schedule_layer_user.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_schedule_layer_user.example <id>
```
