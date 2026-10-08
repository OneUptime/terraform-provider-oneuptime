---
page_title: "oneuptime_on_call_duty_policy_user_owner Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add users as owners to your onCallDutyPolicys.
---

# oneuptime_on_call_duty_policy_user_owner (Resource)

Add users as owners to your onCallDutyPolicys.

## Example Usage

```terraform
resource "oneuptime_on_call_duty_policy_user_owner" "example" {
  user_id                = data.oneuptime_user.example.id
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
}
```

## Schema

### Required

- `on_call_duty_policy_id` (String) ID of your OneUptime OnCallDutyPolicy in which this object belongs. The ID of a `oneuptime_on_call_policy`.
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call duty policy user owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_duty_policy_user_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_duty_policy_user_owner.example <id>
```
