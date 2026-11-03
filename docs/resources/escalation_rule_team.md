---
page_title: "oneuptime_escalation_rule_team Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add teams to your escalation rules.
---

# oneuptime_escalation_rule_team (Resource)

Add teams to your escalation rules.

~> **Renamed:** this resource was called `oneuptime_team_on_call_duty_escalation_rule` before. The old name still works, but is deprecated. To switch, rename the resource in your configuration and add a `moved` block, so Terraform keeps the existing escalation rule team:

```terraform
moved {
  from = oneuptime_team_on_call_duty_escalation_rule.example
  to   = oneuptime_escalation_rule_team.example
}
```

## Example Usage

```terraform
resource "oneuptime_escalation_rule_team" "example" {
  on_call_duty_policy_id                 = oneuptime_on_call_policy.example.id
  on_call_duty_policy_escalation_rule_id = oneuptime_escalation_rule.example.id
}
```

## Schema

### Required

- `on_call_duty_policy_escalation_rule_id` (String) ID of the escalation rule that pages the team. The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of the on-call policy the escalation rule is in. The ID of a `oneuptime_on_call_policy`.

### Optional

- `team_id` (String) ID of the team the escalation rule pages. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing escalation rule team by its escalation rule and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_escalation_rule_team.example
  id = "<on_call_duty_policy_escalation_rule_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_escalation_rule_team.example <on_call_duty_policy_escalation_rule_id>/<team_id>
```

Its own ID works too.
