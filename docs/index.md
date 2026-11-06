---
page_title: "oneuptime Provider"
subcategory: ""
description: |-
  Manage a OneUptime project as code, with Terraform or OpenTofu: monitors, status pages, incidents, on-call policies, teams and everything else the OneUptime API manages.
---

# OneUptime Provider

Manage a OneUptime project as code, with Terraform or OpenTofu: monitors, status pages, incidents, on-call policies, teams and everything else the OneUptime API manages.

## Example Usage

```terraform
terraform {
  required_providers {
    oneuptime = {
      source  = "oneuptime/oneuptime"
      version = "~> 14.0"
    }
  }
}

provider "oneuptime" {
  # Both can come from the environment instead:
  # ONEUPTIME_URL and ONEUPTIME_API_KEY.
  oneuptime_url = "https://oneuptime.com" # or your self-hosted instance
  api_key       = var.oneuptime_api_key
}
```

## Schema

### Optional

- `api_key` (String, Sensitive) An API key of the project to manage, from its **Settings** > **API Keys**. Falls back to the `ONEUPTIME_API_KEY` environment variable; the provider stops before any plan when neither is set.
- `oneuptime_url` (String) The address of your OneUptime: https://oneuptime.com, or your own installation's. Leave out /api: the provider adds it. Falls back to the `ONEUPTIME_URL` environment variable, then to https://oneuptime.com.

## Start Here

The provider covers the full OneUptime API surface. Most configurations begin with these resources:

- [`oneuptime_monitor`](./resources/monitor) — Uptime and health checks for your services
- [`oneuptime_monitor_status`](./resources/monitor_status) — The states a monitor can be in
- [`oneuptime_label`](./resources/label) — Organize resources across the project
- [`oneuptime_status_page`](./resources/status_page) — Public status pages for your users
- [`oneuptime_status_page_domain`](./resources/status_page_domain) — Serve a status page on your own domain
- [`oneuptime_incident_severity`](./resources/incident_severity) — Severity levels for incidents
- [`oneuptime_on_call_duty_policy`](./resources/on_call_duty_policy) — On-call rotations and escalation
- [`oneuptime_team`](./resources/team) — Teams that own monitors and get paged
- [`oneuptime_scheduled_maintenance`](./resources/scheduled_maintenance) — Planned maintenance windows

Every resource has a matching data source of the same name, for referring to something that already exists instead of creating it. Look it up by `id`, or by any of its other arguments - each one you set must match, and exactly one item may match them all:

```terraform
data "oneuptime_monitor_status" "offline" {
  name = "Offline"
}

data "oneuptime_incident_severity" "critical" {
  name = "Critical Incident"
}
```

## IDs inside resources

Server-side ids that a resource carries inside it - the ids of a monitor's steps, criteria and incident templates - are the server's to manage. Never write them: the provider sends none, the server gives them, and it keeps them across every apply, so incidents keep pointing at the criteria that raised them.

## Renamed resources

These resources were renamed to read the way the product names them: IoT and vCenter are one word, auto remediation rules are AI rules, an escalation rule's users, teams and schedules are `escalation_rule_user`, `escalation_rule_team` and `escalation_rule_schedule`, and every team owner is `<thing>_team_owner`. Each old name still works, as a deprecated alias, so no configuration breaks; plans that use one say so.

| Old name | New name |
|----------|----------|
| `oneuptime_ai_agent_owner_team` | [`oneuptime_ai_agent_team_owner`](./resources/ai_agent_team_owner) |
| `oneuptime_auto_remediation_rule` | [`oneuptime_ai_rule`](./resources/ai_rule) |
| `oneuptime_schedules_on_call_duty_escalation_rule` | [`oneuptime_escalation_rule_schedule`](./resources/escalation_rule_schedule) |
| `oneuptime_team_on_call_duty_escalation_rule` | [`oneuptime_escalation_rule_team`](./resources/escalation_rule_team) |
| `oneuptime_users_on_call_duty_escalation_rule` | [`oneuptime_escalation_rule_user`](./resources/escalation_rule_user) |
| `oneuptime_io_t_device_credential` | [`oneuptime_iot_device_credential`](./resources/iot_device_credential) |
| `oneuptime_io_t_fleet` | [`oneuptime_iot_fleet`](./resources/iot_fleet) |
| `oneuptime_io_t_fleet_label_rule` | [`oneuptime_iot_fleet_label_rule`](./resources/iot_fleet_label_rule) |
| `oneuptime_io_t_fleet_owner_rule` | [`oneuptime_iot_fleet_owner_rule`](./resources/iot_fleet_owner_rule) |
| `oneuptime_io_t_fleet_team_owner` | [`oneuptime_iot_fleet_team_owner`](./resources/iot_fleet_team_owner) |
| `oneuptime_io_t_fleet_user_owner` | [`oneuptime_iot_fleet_user_owner`](./resources/iot_fleet_user_owner) |
| `oneuptime_open_stack_cloud` | [`oneuptime_openstack_cloud`](./resources/openstack_cloud) |
| `oneuptime_open_stack_cloud_feed` | [`oneuptime_openstack_cloud_feed`](./resources/openstack_cloud_feed) |
| `oneuptime_open_stack_cloud_label_rule` | [`oneuptime_openstack_cloud_label_rule`](./resources/openstack_cloud_label_rule) |
| `oneuptime_open_stack_cloud_owner_rule` | [`oneuptime_openstack_cloud_owner_rule`](./resources/openstack_cloud_owner_rule) |
| `oneuptime_open_stack_cloud_team_owner` | [`oneuptime_openstack_cloud_team_owner`](./resources/openstack_cloud_team_owner) |
| `oneuptime_open_stack_cloud_user_owner` | [`oneuptime_openstack_cloud_user_owner`](./resources/openstack_cloud_user_owner) |
| `oneuptime_probe_owner_team` | [`oneuptime_probe_team_owner`](./resources/probe_team_owner) |
| `oneuptime_v_center` | [`oneuptime_vcenter`](./resources/vcenter) |
| `oneuptime_v_center_feed` | [`oneuptime_vcenter_feed`](./resources/vcenter_feed) |
| `oneuptime_v_center_label_rule` | [`oneuptime_vcenter_label_rule`](./resources/vcenter_label_rule) |
| `oneuptime_v_center_owner_rule` | [`oneuptime_vcenter_owner_rule`](./resources/vcenter_owner_rule) |
| `oneuptime_v_center_team_owner` | [`oneuptime_vcenter_team_owner`](./resources/vcenter_team_owner) |
| `oneuptime_v_center_user_owner` | [`oneuptime_vcenter_user_owner`](./resources/vcenter_user_owner) |

To switch, rename the resource in your configuration and add a `moved` block, and Terraform keeps the existing resource instead of replacing it:

```terraform
moved {
  from = oneuptime_auto_remediation_rule.example
  to   = oneuptime_ai_rule.example
}
```

One old name is another resource's now, so no alias can keep it. A configuration that uses it for what it managed before stops planning until you rename the resource and add a `moved` block:

| Old name, now another resource | What it managed is now |
|----------|----------|
| [`oneuptime_subscriber_notification_template`](./resources/subscriber_notification_template) | [`oneuptime_status_page_subscriber_notification_template`](./resources/status_page_subscriber_notification_template) |

```terraform
moved {
  from = oneuptime_subscriber_notification_template.example
  to   = oneuptime_status_page_subscriber_notification_template.example
}
```

Where the engine cannot move a resource to another type (Terraform before 1.8), remove it from the state with `terraform state rm` and import it under its new name instead.

## People by email address

Every resource that names a person takes their email address as well as their user ID - `user_email` beside `user_id`, `owner_user_emails` beside `owner_users` - and a team member named by an address that is new to the project invites them. See [Managing people and teams](./guides/people) for onboarding and offboarding as code.

```terraform
resource "oneuptime_team_member" "alice" {
  team_id    = oneuptime_team.sre.id
  user_email = "alice@example.com"
}
```

## OpenTofu

This provider is published to the OpenTofu Registry as well, and its end-to-end suite runs against both engines on every change. Configuration is identical — see the [OpenTofu guide](./guides/opentofu).

## Reusable modules

Hand-written modules ship in the repository under `modules/`, for setups that would otherwise be copy-pasted per service:

- [`monitoring-and-incident-response`](https://github.com/OneUptime/terraform-provider-oneuptime/tree/master/modules/monitoring-and-incident-response) — HTTP monitors, an on-call rotation paged when they fail, and a status page listing them.
