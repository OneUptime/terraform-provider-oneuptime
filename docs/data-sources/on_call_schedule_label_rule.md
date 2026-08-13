---
page_title: "oneuptime_on_call_schedule_label_rule Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Configure rules for automatically attaching labels to on-call schedules when matching schedules are created
---

# oneuptime_on_call_schedule_label_rule (Data Source)

Configure rules for automatically attaching labels to on-call schedules when matching schedules are created Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_on_call_schedule_label_rule" "by_name" {
  name = "example-on_call_schedule_label_rule"
}

data "oneuptime_on_call_schedule_label_rule" "by_id" {
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
- `description` (String) Description of this on-call schedule label rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `on_call_duty_policy_schedule_labels` (Set) Only trigger for on-call schedules that already have at least one of these labels. Leave empty to match regardless of labels... Computed.
- `on_call_duty_policy_schedule_name_pattern` (String) Regex (case-insensitive) matched against the on-call schedule name. Leave empty to match any name... Computed.
- `on_call_duty_policy_schedule_description_pattern` (String) Regex (case-insensitive) matched against the on-call schedule description. Leave empty to match any description... Computed.
- `labels_to_add` (Set) Labels to attach to the on-call schedule when this rule matches. Already-attached labels are not duplicated... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
