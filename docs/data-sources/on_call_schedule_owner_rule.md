---
page_title: "oneuptime_on_call_schedule_owner_rule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Configure rules for automatically assigning owner users and teams when matching on-call schedules are created
---

# oneuptime_on_call_schedule_owner_rule (Data Source)

Configure rules for automatically assigning owner users and teams when matching on-call schedules are created

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call schedule owner rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_schedule_owner_rule" "example" {
  name = "Example on call schedule owner rule"
}

# Or by id:
data "oneuptime_on_call_schedule_owner_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this on-call schedule owner rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this on-call schedule owner rule.
- `notify_owners` (Boolean) Send notifications to owner users and teams when they are added by this rule.
- `on_call_duty_policy_schedule_description_pattern` (String) Regex (case-insensitive) matched against the on-call schedule description. Leave empty to match any description.
- `on_call_duty_policy_schedule_name_pattern` (String) Regex (case-insensitive) matched against the on-call schedule name. Leave empty to match any name.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `on_call_duty_policy_schedule_labels` (Set of String) Only trigger for on-call schedules that have at least one of these labels. Leave empty to match regardless of labels. IDs of `oneuptime_label` resources.
- `owner_teams` (Set of String) Teams to add as owners on the on-call schedule when this rule matches. IDs of `oneuptime_team` resources.
- `owner_users` (Set of String) Users to add as owners on the on-call schedule when this rule matches. IDs of `oneuptime_user` records.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
