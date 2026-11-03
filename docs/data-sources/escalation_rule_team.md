---
page_title: "oneuptime_escalation_rule_team Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add teams to your escalation rules.
---

# oneuptime_escalation_rule_team (Data Source)

Add teams to your escalation rules.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one escalation rule team may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_team_on_call_duty_escalation_rule` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_escalation_rule_team" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
}

# Or by id:
data "oneuptime_escalation_rule_team" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `on_call_duty_policy_escalation_rule_id` (String) ID of the escalation rule that pages the team. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of the on-call policy the escalation rule is in. The ID of a `oneuptime_on_call_policy`.
- `team_id` (String) ID of the team the escalation rule pages. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
