---
page_title: "oneuptime_team_member Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Add users to your teams.
---

# oneuptime_team_member (Resource)

Add users to your teams.

## Example Usage

```terraform
resource "oneuptime_team_member" "example" {
  team_id    = oneuptime_team.example.id
  user_email = "alice@example.com"
}
```

## Invite people by email address

Name the person by `user_email`. Someone who is in the project already joins the team at once; someone who is not is invited, as the dashboard's **Invite** does: OneUptime makes them an account if they have none, and emails them. `has_accepted_invitation` is `false` until they accept, and their accepting is no change to the plan.

To put someone this configuration invites on call, or make them an owner, refer to their membership: Terraform then invites them first.

```terraform
resource "oneuptime_monitor_user_owner" "alice" {
  monitor_id = oneuptime_monitor.example.id
  user_id    = oneuptime_team_member.example.user_id
}
```

## Remove someone

Removing a member from the configuration removes them from the team. Out of their last team of the project, they are out of the project: OneUptime takes them off every on-call schedule, escalation rule and override, and off what they own, except the incidents, alerts and maintenance events that are over. Remove the blocks that name them in the same change. Their OneUptime account stays theirs. A team that must keep a member refuses to lose its last one who has accepted.

## Email addresses

OneUptime keeps addresses in lower case, and finds a person by theirs in any case. The address stays in state as you wrote it, and a person who changes theirs in OneUptime stays the same person to Terraform: nothing changes until you write the new address, and then only the address does.

## Schema

### Required

- `team_id` (String) ID of the team the person is in. The ID of a `oneuptime_team`.

### Optional

- `user_email` (String) The email address of the person in `user_id`, instead of their ID. Someone who is not in the project yet is invited to it, as the dashboard's Invite does: OneUptime emails them, and they join when they accept. Give one of the two.
- `user_id` (String) ID of the person in the team. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `has_accepted_invitation` (Boolean) Whether the person has joined: false while their invitation waits, true once they accept it. Someone in the project already joins at once.
- `id` (String) Unique identifier for the resource.
- `invitation_accepted_at` (String) When the person accepted their invitation.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing team member by its team and the person's email address (or user ID), with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_team_member.example
  id = "<team_id>/alice@example.com"
}
```

or on the command line:

```shell
terraform import oneuptime_team_member.example <team_id>/alice@example.com
```

Imported by an address, the person's `user_email` stays as you wrote it. Its own ID works too.
