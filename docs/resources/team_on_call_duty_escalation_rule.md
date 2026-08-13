---
page_title: "oneuptime_team_on_call_duty_escalation_rule Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Manage on-call duty escalation rule for the on-call policy.
---

# oneuptime_team_on_call_duty_escalation_rule (Resource)

Manage on-call duty escalation rule for the on-call policy.

## Example Usage

```terraform
resource "oneuptime_team_on_call_duty_escalation_rule" "example" {
  on_call_duty_policy_id = "123e4567-e89b-12d3-a456-426614174000"
  on_call_duty_policy_escalation_rule_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policy_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `team_id` (String) A unique identifier for an object, represented as a UUID..
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
terraform import oneuptime_team_on_call_duty_escalation_rule.example <id>
```
