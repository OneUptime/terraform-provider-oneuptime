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
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.example.id
  name                            = "Example on call schedule layer"
  starts_at                       = "2030-01-01T00:00:00Z"
  hand_off_time                   = "2030-01-01T00:00:00Z"
  description                     = "Managed by Terraform"
}
```

## Schema

### Required

- `hand_off_time` (String) Hand off time. When would you like to hand off the duty to the next user in this layer?
- `name` (String) Friendly name for this layer.
- `on_call_duty_policy_schedule_id` (String) ID of your On-Call Policy Schedule where this escalation rule belongs. The ID of a `oneuptime_on_call_policy_schedule`.
- `starts_at` (String) Start date and time of this layer.

### Optional

- `description` (String) Description for this layer. This is optional and can be left blank.
- `order` (Number) Order / Priority of this layer. Lower the number, higher the priority.
- `restriction_times` (String) Restrict this layer to these times. A JSON value: write it with `jsonencode()`.
- `rotation` (String) How often would you like to hand off the duty to the next user in this layer? A JSON value: write it with `jsonencode()`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call schedule layer by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_schedule_layer.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_schedule_layer.example <id>
```
