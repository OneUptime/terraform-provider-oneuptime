---
page_title: "oneuptime_escalation_rule Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call duty escalation rule for the on-call policy.
---

# oneuptime_escalation_rule (Resource)

Manage on-call duty escalation rule for the on-call policy.

## Example Usage

```terraform
resource "oneuptime_escalation_rule" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
  name                   = "Example escalation rule"
  description            = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.
- `on_call_duty_policy_id` (String) ID of your On-Call Policy where this escalation rule belongs. The ID of a `oneuptime_on_call_policy`.

### Optional

- `description` (String) Friendly description that will help you remember.
- `escalate_after_in_minutes` (Number) How long should we wait before we execute the next escalation rule?
- `order` (Number) Order of this rule.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing escalation rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_escalation_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_escalation_rule.example <id>
```
