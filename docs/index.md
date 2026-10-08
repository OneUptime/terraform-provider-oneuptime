---
page_title: "oneuptime Provider"
subcategory: ""
description: |-
  Terraform provider for Oneuptime.
---

# Oneuptime Provider

OpenAPI specification for OneUptime. This document describes the API endpoints, request and response formats, and other details necessary for developers to interact with the OneUptime API.

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

- `api_key` (String, Sensitive) Project-scoped API key for authentication. Falls back to the `ONEUPTIME_API_KEY` environment variable; the provider fails at configure time when neither is set.
- `oneuptime_url` (String) The oneuptime URL (without /api path). Defaults to 'oneuptime.com' if not specified. The provider automatically appends '/api' to the URL. Can also be set via the `ONEUPTIME_URL` environment variable.

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

These resources were renamed so that names like IoT and vCenter read as one word. Each old name still works, as a deprecated alias, so no configuration breaks; plans that use one say so.

| Old name | New name |
|----------|----------|
| `oneuptime_io_t_device_credential` | [`oneuptime_iot_device_credential`](./resources/iot_device_credential) |
| `oneuptime_io_t_fleet` | [`oneuptime_iot_fleet`](./resources/iot_fleet) |
| `oneuptime_io_t_fleet_label_rule` | [`oneuptime_iot_fleet_label_rule`](./resources/iot_fleet_label_rule) |
| `oneuptime_io_t_fleet_owner_rule` | [`oneuptime_iot_fleet_owner_rule`](./resources/iot_fleet_owner_rule) |
| `oneuptime_io_t_fleet_team_owner` | [`oneuptime_iot_fleet_team_owner`](./resources/iot_fleet_team_owner) |
| `oneuptime_io_t_fleet_user_owner` | [`oneuptime_iot_fleet_user_owner`](./resources/iot_fleet_user_owner) |
| `oneuptime_v_center` | [`oneuptime_vcenter`](./resources/vcenter) |
| `oneuptime_v_center_feed` | [`oneuptime_vcenter_feed`](./resources/vcenter_feed) |
| `oneuptime_v_center_label_rule` | [`oneuptime_vcenter_label_rule`](./resources/vcenter_label_rule) |
| `oneuptime_v_center_owner_rule` | [`oneuptime_vcenter_owner_rule`](./resources/vcenter_owner_rule) |
| `oneuptime_v_center_team_owner` | [`oneuptime_vcenter_team_owner`](./resources/vcenter_team_owner) |
| `oneuptime_v_center_user_owner` | [`oneuptime_vcenter_user_owner`](./resources/vcenter_user_owner) |

To switch, rename the resource in your configuration and add a `moved` block, and Terraform keeps the existing resource instead of replacing it:

```terraform
moved {
  from = oneuptime_v_center.example
  to   = oneuptime_vcenter.example
}
```

## OpenTofu

This provider is published to the OpenTofu Registry as well, and its end-to-end suite runs against both engines on every change. Configuration is identical — see the [OpenTofu guide](./guides/opentofu).

## Reusable modules

Hand-written modules ship in the repository under `modules/`, for setups that would otherwise be copy-pasted per service:

- [`monitoring-and-incident-response`](https://github.com/OneUptime/terraform-provider-oneuptime/tree/master/modules/monitoring-and-incident-response) — HTTP monitors, an on-call rotation paged when they fail, and a status page listing them.
