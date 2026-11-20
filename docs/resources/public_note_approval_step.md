---
page_title: "oneuptime_public_note_approval_step Resource - oneuptime"
subcategory: "Other"
description: |-
  Who approves public notes before they are published on status pages and sent to subscribers, step by step
---

# oneuptime_public_note_approval_step (Resource)

Who approves public notes before they are published on status pages and sent to subscribers, step by step

## Example Usage

```terraform
resource "oneuptime_public_note_approval_step" "example" {
  event_type = "Incident"
  name       = "Example public note approval step"
}
```

<a id="order"></a>
## Order

`order` is this public note approval step's place among the public note approval steps of its project with the same event type, lowest first. Setting it places this public note approval step and moves no other, so a configuration that numbers every public note approval step of a project and event type is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.public_note_approval_steps`:

```terraform
order = index(local.public_note_approval_steps, each.key) + 1
```

- Give each public note approval step a number of its own: two public note approval steps at one number are shown the older first.
- A public note approval step created without `order` goes after the last one of its project with the same event type.
- Deleting a public note approval step leaves a gap in the numbers, which changes nothing: public note approval steps go by their numbers, not by counting.
- Moving a public note approval step in the OneUptime Dashboard numbers the public note approval steps of the project with that event type 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `event_type` (String) Whose public notes this step approves: Incident (the public notes of incidents and incident episodes) or ScheduledMaintenance (the public notes of maintenance events). Allowed values: `Incident`, `ScheduledMaintenance`.

### Optional

- `approver_teams` (Set of String) Teams whose members may approve a public note for this step. Someone added to a team can decide on the notes already waiting. IDs of `oneuptime_team` resources.
- `approver_user_emails` (Set of String) The email addresses of the people in `approver_users`, instead of their IDs: each one of someone in this project, or invited to it. Give one of the two.
- `approver_users` (Set of String) People who may approve a public note for this step. Any one of them, or of the members of the teams below, is enough. IDs of `oneuptime_user` records. Or name them by email address, in `approver_user_emails`.
- `name` (String) What the step is called, shown to approvers and authors (for example, Legal). Optional.
- `order` (Number) Where this step sits among the approval steps, lowest number first: a public note goes through them in this order. A new step is added to the end. Setting a number places the step there and moves no other: two steps at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing public note approval step by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_public_note_approval_step.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_public_note_approval_step.example <id>
```
