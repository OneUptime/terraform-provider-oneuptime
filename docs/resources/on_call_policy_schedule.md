---
page_title: "oneuptime_on_call_policy_schedule Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage schedules and rotations for your on-call duty policy.
---

# oneuptime_on_call_policy_schedule (Resource)

Manage schedules and rotations for your on-call duty policy.

## Example Usage

```terraform
resource "oneuptime_on_call_policy_schedule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `labels` (Set) Relation to Labels Array where this object is categorized in...
- `description` (String) Friendly description that will help you remember..
- `timezone` (String) IANA timezone this schedule's restriction and hand-off wall-clock times are interpreted in. When empty, times are interpreted in the server's local timezone (legacy behavior)...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `current_user_id_on_roster` (String) A unique identifier for an object, represented as a UUID..
- `next_user_id_on_roster` (String) A unique identifier for an object, represented as a UUID..
- `roster_handoff_at` (String) A date time object..
- `roster_next_handoff_at` (String) A date time object..
- `roster_next_start_at` (String) A date time object..
- `roster_start_at` (String) A date time object..
- `shift_config_version` (Number) Incremented whenever the schedule's layers, members, overrides or policy attachments change. Used as the calendar feed SEQUENCE...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_on_call_policy_schedule.example <id>
```
