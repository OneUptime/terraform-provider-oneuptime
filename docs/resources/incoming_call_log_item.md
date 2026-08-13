---
page_title: "oneuptime_incoming_call_log_item Resource - oneuptime"
subcategory: "Other"
description: |-
  Child log for each escalation attempt / user ring within a call.
---

# oneuptime_incoming_call_log_item (Resource)

Child log for each escalation attempt / user ring within a call.

## Example Usage

```terraform
resource "oneuptime_incoming_call_log_item" "example" {
  incoming_call_log_id = "123e4567-e89b-12d3-a456-426614174000"
  status = "Example short text"
}
```

## Schema

### Required

- `incoming_call_log_id` (String) A unique identifier for an object, represented as a UUID..
- `status` (String) Status of this dial attempt..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `incoming_call_policy_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID..
- `user_id` (String) A unique identifier for an object, represented as a UUID..
- `user_phone_number` (String) Phone object.
- `status_message` (String) Additional status information..
- `dial_duration_in_seconds` (Number) How long this dial lasted in seconds..
- `call_cost_in_usd_cents` (Number) Cost for this dial attempt in USD cents..
- `started_at` (String) A date time object..
- `ended_at` (String) A date time object..
- `is_answered` (Bool) Whether this user answered the call..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incoming_call_log_item.example <id>
```
