---
page_title: "oneuptime_slo_burn_rate_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Configure multi-window burn rate rules that raise alerts and/or declare incidents when a Service Level Objective consumes its error budget too quickly
---

# oneuptime_slo_burn_rate_rule (Data Source)

Configure multi-window burn rate rules that raise alerts and/or declare incidents when a Service Level Objective consumes its error budget too quickly Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_slo_burn_rate_rule" "by_name" {
  name = "example-slo_burn_rate_rule"
}

data "oneuptime_slo_burn_rate_rule" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `service_level_objective_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_enabled` (Bool) Whether this burn rate rule is enabled.. Computed.
- `burn_rate_threshold` (Number) Alert when the burn rate in both the long and short windows is at or above this threshold (e.g. 14.4)... Computed.
- `long_window_in_minutes` (Number) Length of the long lookback window in minutes (e.g. 60). The alert fires when both windows exceed the threshold and resolves when the long window drops below it... Computed.
- `short_window_in_minutes` (Number) Length of the short lookback window in minutes (e.g. 5). Guards against alerting on burn that has already stopped... Computed.
- `minimum_sample_count` (Number) For event-based SLIs only: skip this rule when the long window has fewer than this many total events. Prevents noisy alerts on low traffic... Computed.
- `refire_suppression_minutes` (Number) Minimum number of minutes after an alert or incident resolves before this rule can declare that same record again. Each output is suppressed independently, from its own resolve. Defaults to the long window length when not set... Computed.
- `should_create_alert` (Bool) Raise an Alert when this burn rate rule fires. Enabled by default... Computed.
- `should_create_incident` (Bool) Declare an Incident when this burn rate rule fires. Disabled by default... Computed.
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policies` (Set) On-call duty policies attached to alerts created by this burn rate rule. Incidents have their own list... Computed.
- `alert_title_template` (String) Title of the alert raised when this burn rate rule fires. Supports template variables such as {{sloName}}. Leave empty to use the default title... Computed.
- `alert_description_template` (String) Description (in Markdown) of the alert raised when this burn rate rule fires. Supports template variables. Leave empty to use the default description... Computed.
- `alert_remediation_notes` (String) Remediation notes (in Markdown) attached to the alert raised when this burn rate rule fires. Supports template variables... Computed.
- `is_alert_private` (Bool) Make the alert raised by this burn rate rule private, so only its owners, project admins and project owners can see it. Disabled by default... Computed.
- `auto_resolve_alert` (Bool) Resolve the alert automatically when the burn rate over the long window drops back below the threshold. Enabled by default. When disabled, the alert stays open until someone resolves it... Computed.
- `alert_labels` (Set) Labels added to alerts raised by this burn rate rule... Computed.
- `alert_owner_teams` (Set) Teams added as owners of alerts raised by this burn rate rule... Computed.
- `alert_owner_users` (Set) Users added as owners of alerts raised by this burn rate rule... Computed.
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_on_call_duty_policies` (Set) On-call duty policies attached to incidents declared by this burn rate rule... Computed.
- `incident_title_template` (String) Title of the incident declared when this burn rate rule fires. Supports template variables such as {{sloName}}. Leave empty to use the default title... Computed.
- `incident_description_template` (String) Description (in Markdown) of the incident declared when this burn rate rule fires. Supports template variables. Leave empty to use the default description... Computed.
- `incident_remediation_notes` (String) Remediation notes (in Markdown) attached to the incident declared when this burn rate rule fires. Supports template variables... Computed.
- `is_incident_private` (Bool) Make the incident declared by this burn rate rule private, so only its owners, project admins and project owners can see it. Disabled by default... Computed.
- `auto_resolve_incident` (Bool) Resolve the incident automatically when the burn rate over the long window drops back below the threshold. Enabled by default. When disabled, the incident stays open until someone resolves it... Computed.
- `incident_labels` (Set) Labels added to incidents declared by this burn rate rule... Computed.
- `incident_owner_teams` (Set) Teams added as owners of incidents declared by this burn rate rule... Computed.
- `incident_owner_users` (Set) Users added as owners of incidents declared by this burn rate rule... Computed.
- `add_slo_owners_as_owners` (Bool) Also add the owner users and owner teams of the Service Level Objective as owners of the alerts and incidents this burn rate rule creates. Disabled by default... Computed.
- `last_alert_created_at` (String) A date time object.. Computed.
- `last_alert_resolved_at` (String) A date time object.. Computed.
- `last_incident_created_at` (String) A date time object.. Computed.
- `last_incident_resolved_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
