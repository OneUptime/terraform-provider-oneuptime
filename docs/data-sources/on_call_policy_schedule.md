---
page_title: "oneuptime_on_call_policy_schedule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage schedules and rotations for your on-call duty policy.
---

# oneuptime_on_call_policy_schedule (Data Source)

Manage schedules and rotations for your on-call duty policy. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_on_call_policy_schedule" "by_name" {
  name = "example-on_call_policy_schedule"
}

data "oneuptime_on_call_policy_schedule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
- `description` (String) Friendly description that will help you remember.. Computed.
- `timezone` (String) IANA timezone this schedule's restriction and hand-off wall-clock times are interpreted in. When empty, times are interpreted in the server's local timezone (legacy behavior)... Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `current_user_id_on_roster` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `next_user_id_on_roster` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `roster_handoff_at` (String) A date time object.. Computed.
- `roster_next_handoff_at` (String) A date time object.. Computed.
- `roster_next_start_at` (String) A date time object.. Computed.
- `roster_start_at` (String) A date time object.. Computed.
