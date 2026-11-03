---
page_title: "oneuptime_escalation_rule_user Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add users to your escalation rules.
---

# oneuptime_escalation_rule_user (Resource)

Add users to your escalation rules.

~> **Renamed:** this resource was called `oneuptime_users_on_call_duty_escalation_rule` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing escalation rule user:

```terraform
moved {
  from = oneuptime_users_on_call_duty_escalation_rule.example
  to   = oneuptime_escalation_rule_user.example
}
```

## Example Usage

```terraform
resource "oneuptime_escalation_rule_user" "example" {
  on_call_duty_policy_id                 = oneuptime_on_call_policy.example.id
  on_call_duty_policy_escalation_rule_id = oneuptime_escalation_rule.example.id
  user_email                             = "alice@example.com"
}
```

## Schema

### Required

- `on_call_duty_policy_escalation_rule_id` (String) ID of the escalation rule that pages the person. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of the on-call policy the escalation rule is in. The ID of a `oneuptime_on_call_policy`.

### Optional

- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of the person this escalation rule pages. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing escalation rule user by its escalation rule and the person's email address (or user ID), with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_escalation_rule_user.example
  id = "<on_call_duty_policy_escalation_rule_id>/alice@example.com"
}
```

or on the command line:

```shell
terraform import oneuptime_escalation_rule_user.example <on_call_duty_policy_escalation_rule_id>/alice@example.com
```

Imported by an address, the person's `user_email` stays as you wrote it. Its own ID works too.
