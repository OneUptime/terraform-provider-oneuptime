---
page_title: "oneuptime_security_playbook Data Source - oneuptime"
subcategory: "Other"
description: |-
  Automated responses to security cases: steps that assign, investigate, notify, escalate and contain, run when a matching case opens or by hand.
---

# oneuptime_security_playbook (Data Source)

Automated responses to security cases: steps that assign, investigate, notify, escalate and contain, run when a matching case opens or by hand.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security playbook may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_playbook" "example" {
  name = "Example security playbook"
}

# Or by id:
data "oneuptime_security_playbook" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) ID of the person who created the playbook. The ID of a `oneuptime_user` (see the data source).
- `description` (String) When to use the playbook and what it is for.
- `detection_source` (String) Run only on cases raised by this source: a Detection Rule, Threat Intel or a Connected Product. Empty runs on every source.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether the playbook runs on its own. A disabled playbook can still be run by hand from a case.
- `minimum_priority` (String) Run only on cases at this priority or above. Empty runs on every priority.
- `name` (String) What the playbook does, in a few words.
- `run_count` (Number) How many times the playbook has been run.
- `trigger` (String) When the playbook runs on its own: Case Opened, Detection Added (a new detection joins an open case), or Manual (only when somebody runs it).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `detection_rule_ids` (String) Run only on cases with a detection from one of these detection rules. Empty runs on cases from any rule. A JSON value: write it with `jsonencode()`.
- `last_run_at` (String) When the playbook last started a run.
- `project_id` (String) ID of the project this playbook belongs to. The ID of a `oneuptime_project`.
- `steps` (String) The steps the playbook takes, in order. Each has an action, its settings, and whether it needs approval, may fail without stopping the run, or runs only when threat intel matched. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
