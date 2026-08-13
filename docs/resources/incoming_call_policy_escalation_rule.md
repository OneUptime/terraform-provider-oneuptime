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
  incoming_call_policy_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `incoming_call_policy_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Rule name (e.g., 'Primary On-Call', 'Backup Engineer')..
- `description` (String) Optional description of this escalation rule..
- `order` (Number) Execution order (1, 2, 3...)..
- `escalate_after_seconds` (Number) Seconds before escalating to next rule..
- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID..
- `user_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incoming_call_policy_escalation_rule.example <id>
```
