---
page_title: "oneuptime_incoming_call_policy_escalation_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Manage escalation rules for incoming call policies that define who to call and in what order
---

# oneuptime_incoming_call_policy_escalation_rule (Resource)

Manage escalation rules for incoming call policies that define who to call and in what order

## Example Usage

```terraform
resource "oneuptime_incoming_call_policy_escalation_rule" "example" {
  incoming_call_policy_id = oneuptime_incoming_call_policy.example.id
  name                    = "Example incoming call policy escalation rule"
  description             = "Managed by Terraform"
}
```

## Schema

### Required

- `incoming_call_policy_id` (String) ID of the parent Incoming Call Policy. The ID of a `oneuptime_incoming_call_policy`.

### Optional

- `description` (String) Optional description of this escalation rule.
- `escalate_after_seconds` (Number) How long, in seconds, the phone rings before the call moves on to the next rule. 20 when left out; a time below 5 or above 600 rings for 5 or 600, the limits Twilio takes. Defaults to `20`.
- `name` (String) Rule name (e.g., 'Primary On-Call', 'Backup Engineer').
- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule to route to (mutually exclusive with userId). The ID of a `oneuptime_on_call_policy_schedule`.
- `order` (Number) Where this rule sits in the escalation, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `user_id` (String) ID of the user to route to directly (mutually exclusive with onCallDutyPolicyScheduleId). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incoming call policy escalation rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incoming_call_policy_escalation_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incoming_call_policy_escalation_rule.example <id>
```
