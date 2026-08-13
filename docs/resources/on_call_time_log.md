---
page_title: "oneuptime_on_call_time_log Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage on-call duty user overrides, for example if the user is on leave you can override the on-call duty policy for that user so all the alerts will be routed to the other user.
---

# oneuptime_on_call_time_log (Resource)

Manage on-call duty user overrides, for example if the user is on leave you can override the on-call duty policy for that user so all the alerts will be routed to the other user.

## Example Usage

```terraform
resource "oneuptime_on_call_time_log" "example" {
  user_id = "123e4567-e89b-12d3-a456-426614174000"
  starts_at = "2030-01-01T00:00:00Z"
}
```

## Schema

### Required

- `user_id` (String) A unique identifier for an object, represented as a UUID..
- `starts_at` (String) A date time object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID..
- `on_call_duty_policy_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID..
- `team_id` (String) A unique identifier for an object, represented as a UUID..
- `more_info` (String) More information about this log record...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `ends_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_on_call_time_log.example <id>
```
