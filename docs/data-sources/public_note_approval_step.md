---
page_title: "oneuptime_public_note_approval_step Data Source - oneuptime"
subcategory: "Other"
description: |-
  Who approves public notes before they are published on status pages and sent to subscribers, step by step
---

# oneuptime_public_note_approval_step (Data Source)

Who approves public notes before they are published on status pages and sent to subscribers, step by step

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one public note approval step may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_public_note_approval_step" "example" {
  name = "Example public note approval step"
}

# Or by id:
data "oneuptime_public_note_approval_step" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `event_type` (String) Whose public notes this step approves: Incident (the public notes of incidents and incident episodes) or ScheduledMaintenance (the public notes of scheduled maintenance events).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `name` (String) What the step is called, shown to approvers and authors (for example, Legal). Optional.
- `order` (Number) Where this step sits among the approval steps, lowest number first: a public note goes through them in this order. A new step is added to the end. Setting a number another step already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.

### Read-Only

- `approver_teams` (Set of String) Teams whose members may approve a public note for this step. Someone added to a team can decide on the notes already waiting. IDs of `oneuptime_team` resources.
- `approver_users` (Set of String) People who may approve a public note for this step. Any one of them, or of the members of the teams below, is enough. IDs of `oneuptime_user` records.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
