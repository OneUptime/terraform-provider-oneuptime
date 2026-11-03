---
page_title: "oneuptime_incoming_call_log_item Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Child log for each escalation attempt / user ring within a call.
---

# oneuptime_incoming_call_log_item (Resource)

Child log for each escalation attempt / user ring within a call.

## Example Usage

```terraform
resource "oneuptime_incoming_call_log_item" "example" {
  incoming_call_log_id    = data.oneuptime_incoming_call_log.example.id
  incoming_call_policy_id = oneuptime_incoming_call_policy.example.id
  status                  = "Initiated"
}
```

## Schema

### Required

- `incoming_call_log_id` (String) ID of the parent Incoming Call Log. The ID of a `oneuptime_incoming_call_log` (see the data source).
- `incoming_call_policy_id` (String) ID of the incoming call policy whose number took the call this attempt is part of. The ID of a `oneuptime_incoming_call_policy`.
- `status` (String) Status of this dial attempt. Allowed values: `Initiated`, `Ringing`, `Connected`, `Escalated`, `NoAnswer`, `Failed`, `Completed`, `CallerHungUp`, `Busy`.

### Optional

- `call_cost_in_usd_cents` (Number) Cost for this dial attempt in USD cents.
- `dial_duration_in_seconds` (Number) How long this dial lasted in seconds.
- `ended_at` (String) When dial ended.
- `incoming_call_policy_escalation_rule_id` (String) ID of the escalation rule used. The ID of a `oneuptime_incoming_call_policy_escalation_rule`.
- `is_answered` (Boolean) Whether this user answered the call. Defaults to `false`.
- `started_at` (String) When dial started.
- `status_message` (String) Additional status information.
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) User ID who was called. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.
- `user_phone_number` (String) Phone number that was dialed.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incoming call log item by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incoming_call_log_item.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incoming_call_log_item.example <id>
```
