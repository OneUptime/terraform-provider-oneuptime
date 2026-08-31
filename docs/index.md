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
      source = "oneuptime/oneuptime"
      version = "12.0.29"
    }
  }
}

provider "oneuptime" {
  oneuptime_url = "oneuptime.com"  # Optional, defaults to oneuptime.com (internally becomes oneuptime.com/api)
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

Every resource has a matching data source of the same name for looking up existing items by `id` or `name`.

## OpenTofu

This provider is published to the OpenTofu Registry as well, and its end-to-end suite runs against both engines on every change. Configuration is identical — see the [OpenTofu guide](./guides/opentofu).

## Reusable modules

Hand-written modules ship in the repository under `modules/`, for setups that would otherwise be copy-pasted per service:

- [`monitoring-and-incident-response`](https://github.com/OneUptime/terraform-provider-oneuptime/tree/master/modules/monitoring-and-incident-response) — HTTP monitors, an on-call rotation paged when they fail, and a status page listing them.
