---
page_title: "oneuptime_incoming_call_policy_user_owner Resource - oneuptime"
subcategory: "Other"
description: |-
  Add users as owners to your incoming call policies.
---

# oneuptime_incoming_call_policy_user_owner (Resource)

Add users as owners to your incoming call policies.

## Example Usage

```terraform
resource "oneuptime_incoming_call_policy_user_owner" "example" {
  user_id                 = data.oneuptime_user.example.id
  incoming_call_policy_id = oneuptime_incoming_call_policy.example.id
}
```

## Schema

### Required

- `incoming_call_policy_id` (String) ID of your OneUptime Incoming Call Policy in which this object belongs. The ID of a `oneuptime_incoming_call_policy`.
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incoming call policy user owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incoming_call_policy_user_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incoming_call_policy_user_owner.example <id>
```
