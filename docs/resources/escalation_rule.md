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
  on_call_duty_policy_id = "123e4567-e89b-12d3-a456-426614174000"
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID..
- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `escalate_after_in_minutes` (Number) How long should we wait before we execute the next escalation rule?..
- `order` (Number) Order of this rule..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_escalation_rule.example <id>
```
