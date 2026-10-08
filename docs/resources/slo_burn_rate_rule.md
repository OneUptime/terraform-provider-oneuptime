---
page_title: "oneuptime_slo_burn_rate_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Configure multi-window burn rate rules that raise alerts and/or declare incidents when a Service Level Objective consumes its error budget too quickly
---

# oneuptime_slo_burn_rate_rule (Resource)

Configure multi-window burn rate rules that raise alerts and/or declare incidents when a Service Level Objective consumes its error budget too quickly

## Example Usage

```terraform
resource "oneuptime_slo_burn_rate_rule" "example" {
  service_level_objective_id = oneuptime_service_level_objective.example.id
  name                       = "Example slo burn rate rule"
  burn_rate_threshold        = 42
  long_window_in_minutes     = 42
  short_window_in_minutes    = 42
}
```

## Schema

### Required

- `burn_rate_threshold` (Number) Alert when the burn rate in both the long and short windows is at or above this threshold (e.g. 14.4).
- `long_window_in_minutes` (Number) Length of the long lookback window in minutes (e.g. 60). The alert fires when both windows exceed the threshold and resolves when the long window drops below it.
- `name` (String) Name of this burn rate rule.
- `service_level_objective_id` (String) ID of the Service Level Objective this burn rate rule belongs to. The ID of a `oneuptime_service_level_objective`.
- `short_window_in_minutes` (Number) Length of the short lookback window in minutes (e.g. 5). Guards against alerting on burn that has already stopped.

### Optional

- `add_slo_owners_as_owners` (Boolean) Also add the owner users and owner teams of the Service Level Objective as owners of the alerts and incidents this burn rate rule creates. Disabled by default. Defaults to `false`.
- `alert_description_template` (String) Description (in Markdown) of the alert raised when this burn rate rule fires. Supports template variables. Leave empty to use the default description.
- `alert_labels` (Set of String) Labels added to alerts raised by this burn rate rule. IDs of `oneuptime_label` resources.
- `alert_owner_teams` (Set of String) Teams added as owners of alerts raised by this burn rate rule. IDs of `oneuptime_team` resources.
- `alert_owner_users` (Set of String) Users added as owners of alerts raised by this burn rate rule. IDs of `oneuptime_user` records.
- `alert_remediation_notes` (String) Remediation notes (in Markdown) attached to the alert raised when this burn rate rule fires. Supports template variables.
- `alert_severity_id` (String) ID of the Alert Severity of the alert created when this burn rate rule fires. The ID of a `oneuptime_alert_severity`.
- `alert_title_template` (String) Title of the alert raised when this burn rate rule fires. Supports template variables such as {{sloName}}. Leave empty to use the default title.
- `auto_resolve_alert` (Boolean) Resolve the alert automatically when the burn rate over the long window drops back below the threshold. Enabled by default. When disabled, the alert stays open until someone resolves it. Defaults to `true`.
- `auto_resolve_incident` (Boolean) Resolve the incident automatically when the burn rate over the long window drops back below the threshold. Enabled by default. When disabled, the incident stays open until someone resolves it. Defaults to `true`.
- `incident_description_template` (String) Description (in Markdown) of the incident declared when this burn rate rule fires. Supports template variables. Leave empty to use the default description.
- `incident_labels` (Set of String) Labels added to incidents declared by this burn rate rule. IDs of `oneuptime_label` resources.
- `incident_on_call_duty_policies` (Set of String) On-call duty policies attached to incidents declared by this burn rate rule. IDs of `oneuptime_on_call_policy` resources.
- `incident_owner_teams` (Set of String) Teams added as owners of incidents declared by this burn rate rule. IDs of `oneuptime_team` resources.
- `incident_owner_users` (Set of String) Users added as owners of incidents declared by this burn rate rule. IDs of `oneuptime_user` records.
- `incident_remediation_notes` (String) Remediation notes (in Markdown) attached to the incident declared when this burn rate rule fires. Supports template variables.
- `incident_severity_id` (String) ID of the Incident Severity of the incident declared when this burn rate rule fires. The ID of a `oneuptime_incident_severity`.
- `incident_title_template` (String) Title of the incident declared when this burn rate rule fires. Supports template variables such as {{sloName}}. Leave empty to use the default title.
- `is_alert_private` (Boolean) Make the alert raised by this burn rate rule private, so only its owners, project admins and project owners can see it. Disabled by default. Defaults to `false`.
- `is_enabled` (Boolean) Whether this burn rate rule is enabled. Defaults to `true`.
- `is_incident_private` (Boolean) Make the incident declared by this burn rate rule private, so only its owners, project admins and project owners can see it. Disabled by default. Defaults to `false`.
- `minimum_sample_count` (Number) For event-based SLIs only: skip this rule when the long window has fewer than this many total events. Prevents noisy alerts on low traffic.
- `on_call_duty_policies` (Set of String) On-call duty policies attached to alerts created by this burn rate rule. Incidents have their own list. IDs of `oneuptime_on_call_policy` resources.
- `refire_suppression_minutes` (Number) Minimum number of minutes after an alert or incident resolves before this rule can declare that same record again. Each output is suppressed independently, from its own resolve. Defaults to the long window length when not set.
- `should_create_alert` (Boolean) Raise an Alert when this burn rate rule fires. Enabled by default. Defaults to `true`.
- `should_create_incident` (Boolean) Declare an Incident when this burn rate rule fires. Disabled by default. Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_alert_created_at` (String) The last time an alert was created by this burn rate rule. Computed by the worker.
- `last_alert_resolved_at` (String) The last time an alert created by this burn rate rule was resolved. Computed by the worker.
- `last_incident_created_at` (String) The last time an incident was declared by this burn rate rule. Computed by the worker.
- `last_incident_resolved_at` (String) The last time an incident declared by this burn rate rule was resolved. Computed by the worker.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing slo burn rate rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_slo_burn_rate_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_slo_burn_rate_rule.example <id>
```
