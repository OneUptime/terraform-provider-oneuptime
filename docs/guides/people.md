---
page_title: "Managing people and teams"
subcategory: "Guides"
description: |-
  Teams and their members by email address: invitations, on-call, owners and offboarding as code.
---

# Managing people and teams

Every resource that names a person takes their email address as well as their user ID: `user_email` beside `user_id`, `owner_user_emails` beside `owner_users`, `assigned_to_user_email` beside `assigned_to_user_id`, and the same everywhere else. Give one of the two. The provider finds the person among the project's members and the people invited to it, and sends their ID. It searches only the API key's project: nobody outside it is ever found.

## Onboarding

A team member named by an address nobody in the project has is invited by it, as the dashboard's **Invite** does: OneUptime makes them an account if they have none, and emails them. Someone in the project already joins the team at once.

```terraform
resource "oneuptime_team" "sre" {
  name = "SRE"
}

resource "oneuptime_team_member" "alice" {
  team_id    = oneuptime_team.sre.id
  user_email = "alice@example.com"
}

resource "oneuptime_on_call_policy_schedule" "weekly" {
  name     = "SRE weekly"
  timezone = "UTC"
}

resource "oneuptime_on_call_schedule_layer" "weekly" {
  on_call_duty_policy_schedule_id = oneuptime_on_call_policy_schedule.weekly.id
  name                            = "Weekly"
  starts_at                       = "2030-01-07T09:00:00Z"
  hand_off_time                   = "2030-01-14T09:00:00Z"
}

# Alice takes turns on call from her first day: invited people can.
resource "oneuptime_on_call_schedule_layer_user" "alice" {
  on_call_duty_policy_schedule_id       = oneuptime_on_call_policy_schedule.weekly.id
  on_call_duty_policy_schedule_layer_id = oneuptime_on_call_schedule_layer.weekly.id
  user_id                               = oneuptime_team_member.alice.user_id
}

resource "oneuptime_monitor" "api" {
  name         = "Public API"
  monitor_type = "Manual"
}

# Bob is in the project already: name him by address anywhere.
resource "oneuptime_monitor_user_owner" "api" {
  monitor_id = oneuptime_monitor.api.id
  user_email = "bob@example.com"
}
```

Refer to a membership this configuration makes (`user_id = oneuptime_team_member.alice.user_id`) to name the person it invites: Terraform then invites them before it names them anywhere else. An address nobody in the project has is a warning in the plan, and an error in the apply.

The API key needs a permission that may read the project's team members (Project Member, Project Admin or Read Teams) to find people, and one that may invite them (Project Admin, or Invite New Members) to add them.

## Finding someone

`data "oneuptime_user"` looks someone up by `email`, `name` or `id`, and gives their ID, name, whether they have joined, and the teams they are in.

```terraform
data "oneuptime_user" "alice" {
  email = "alice@example.com"
}
```

## Offboarding

Remove the person's memberships from the configuration. Out of their last team, they are out of the project, and OneUptime takes them off every on-call schedule, escalation rule and override, and off what they own, except the incidents, alerts and maintenance events that are over. Remove the blocks that name them in the same change. Their OneUptime account stays theirs; only the project is left.

## Importing

A membership no screen shows the ID of is imported by its team and the person's address:

```shell
terraform import oneuptime_team_member.alice <team id>/alice@example.com
```

## When someone changes their address

The address stays in state as you wrote it, and the person stays the same person to Terraform: nothing changes until you write the new address, and then only the address does. OneUptime keeps addresses in lower case and compares them in any case.

## With an identity provider

If your identity provider manages your people over SCIM, let it: while SCIM Push Groups manages a project's teams, OneUptime refuses team member changes from anywhere else, Terraform included. Keep teams, on-call and owners in Terraform, named by address.
