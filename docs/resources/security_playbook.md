---
page_title: "oneuptime_security_playbook Resource - oneuptime"
subcategory: "Other"
description: |-
  Automated responses to security cases: steps that assign, investigate, notify, escalate and contain, run when a matching case opens or by hand.
---

# oneuptime_security_playbook (Resource)

Automated responses to security cases: steps that assign, investigate, notify, escalate and contain, run when a matching case opens or by hand.

## Example Usage

```terraform
resource "oneuptime_security_playbook" "example" {
  name        = "Example security playbook"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) What the playbook does, in a few words.

### Optional

- `description` (String) When to use the playbook and what it is for.
- `detection_rule_ids` (String) Run only on cases with a detection from one of these detection rules. Empty runs on cases from any rule. A JSON value: write it with `jsonencode()`.
- `detection_source` (String) Run only on cases raised by this source: a Detection Rule, Threat Intel or a Connected Product. Empty runs on every source.
- `is_enabled` (Boolean) Whether the playbook runs on its own. A disabled playbook can still be run by hand from a case. Defaults to `true`.
- `minimum_priority` (String) Run only on cases at this priority or above. Empty runs on every priority.
- `steps` (String) The steps the playbook takes, in order. Each has an action, its settings, and whether it needs approval, may fail without stopping the run, or runs only when threat intel matched. A JSON value: write it with `jsonencode()`.
- `trigger` (String) When the playbook runs on its own: Case Opened, Detection Added (a new detection joins an open case), or Manual (only when somebody runs it). Defaults to `Case Opened`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who created the playbook. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_run_at` (String) When the playbook last started a run.
- `project_id` (String) ID of the project this playbook belongs to. The ID of a `oneuptime_project`.
- `run_count` (Number) How many times the playbook has been run.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security playbook by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_playbook.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_playbook.example <id>
```
