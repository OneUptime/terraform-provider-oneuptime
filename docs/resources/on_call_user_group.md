---
page_title: "oneuptime_on_call_user_group Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  A Slack user group or a Microsoft Teams tag whose members OneUptime keeps equal to whoever is on call right now in the on-call schedules it follows, so anyone can mention it to reach them.
---

# oneuptime_on_call_user_group (Resource)

A Slack user group or a Microsoft Teams tag whose members OneUptime keeps equal to whoever is on call right now in the on-call schedules it follows, so anyone can mention it to reach them.

## Example Usage

```terraform
resource "oneuptime_on_call_user_group" "example" {
  workspace_type                = "Example short text"
  workspace_group_id            = "This is an example of longer text content that might be stored in this field."
  on_call_duty_policy_schedules = [oneuptime_on_call_policy_schedule.example.id]
  name                          = "Example on call user group"
}
```

## Schema

### Required

- `on_call_duty_policy_schedules` (Set of String) The on-call schedules the group follows: it holds whoever is on call in any of them right now. IDs of `oneuptime_on_call_policy_schedule` resources.
- `workspace_group_id` (String) The group's id in the workspace: a Slack user group id, such as S0604QSJC, or a Microsoft Teams tag id.
- `workspace_type` (String) Where the group is: Slack (a user group) or MicrosoftTeams (a tag in a team).

### Optional

- `handle` (String) Slack only: the handle people type to mention the group, without the @, such as payments-oncall.
- `name` (String) The group's name in the workspace, such as Payments on-call. Read from Slack or Microsoft Teams when the group is linked.
- `workspace_team_id` (String) Microsoft Teams only: the team the tag is in (its Microsoft 365 group id).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `next_sync_at` (String) When OneUptime looks at the group again, besides every hand-off: to undo changes made to it by hand, and to try again after the workspace refused a change.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `sync_state` (String) What the last sync came to: who OneUptime put in the group, who is on call but could not be added and why, whether the group is turned off while nobody is on call, and what stopped the sync, if anything. A JSON value: write it with `jsonencode()`.
- `synced_at` (String) When the group last matched who was on call.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing on call user group by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_on_call_user_group.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_on_call_user_group.example <id>
```
