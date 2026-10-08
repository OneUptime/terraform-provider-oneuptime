---
page_title: "oneuptime_on_call_policy_schedule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Manage schedules and rotations for your on-call duty policy.
---

# oneuptime_on_call_policy_schedule (Data Source)

Manage schedules and rotations for your on-call duty policy.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call policy schedule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_policy_schedule" "example" {
  name = "Example on call policy schedule"
}

# Or by id:
data "oneuptime_on_call_policy_schedule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `current_user_id_on_roster` (String) User ID who is currently on roster. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) Any friendly name of this object.
- `next_user_id_on_roster` (String) Next ID who is currently on roster. The ID of a `oneuptime_user` (see the data source).
- `shift_config_version` (Number) Incremented whenever the schedule's layers, members, overrides or policy attachments change. Used as the calendar feed SEQUENCE.
- `slug` (String) Friendly globally unique name for your object.
- `timezone` (String) IANA timezone this schedule's restriction and hand-off wall-clock times are interpreted in. When empty, times are interpreted in the server's local timezone (legacy behavior).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) Relation to Labels Array where this object is categorized in. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `roster_handoff_at` (String) When does the roster handoff occur for this schedule for the current user?
- `roster_next_handoff_at` (String) When does the next roster handoff occur for this schedule for the next user?
- `roster_next_start_at` (String) When does the next event start for this schedule for the next user?
- `roster_start_at` (String) When does the current event start for this schedule for the current user?
- `updated_at` (String) Date and Time when the object was updated.
