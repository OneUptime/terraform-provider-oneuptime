---
page_title: "oneuptime_on_call_schedule_team_owner Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Add teams as owners to your on-call schedules.
---

# oneuptime_on_call_schedule_team_owner (Resource)

Add teams as owners to your on-call schedules.

## Example Usage

```terraform
resource "oneuptime_on_call_schedule_team_owner" "example" {
  team_id                         = oneuptime_team.example.id
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.example.id
}
```

## Schema

### Required

- `on_call_duty_policy_schedule_id` (String) ID of the on-call schedule the team owns. The ID of a `oneuptime_on_call_policy_schedule`.
- `team_id` (String) ID of the team that owns the on-call schedule. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call schedule team owner by its on call policy schedule and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_schedule_team_owner.example
  id = "<on_call_duty_policy_schedule_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_schedule_team_owner.example <on_call_duty_policy_schedule_id>/<team_id>
```

Its own ID works too.
