---
page_title: "oneuptime_slo_error_budget_policy_tier Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Tiers of a Service Level Objective's error budget policy: what the team does as the error budget runs down, and the alerts and incidents raised when the SLO enters a tier
---

# oneuptime_slo_error_budget_policy_tier (Data Source)

Tiers of a Service Level Objective's error budget policy: what the team does as the error budget runs down, and the alerts and incidents raised when the SLO enters a tier

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one slo error budget policy tier may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_slo_error_budget_policy_tier" "example" {
  name = "Example slo error budget policy tier"
}

# Or by id:
data "oneuptime_slo_error_budget_policy_tier" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `add_slo_owners_as_owners` (Boolean) Also add the owner users and owner teams of the Service Level Objective as owners of the alerts and incidents this tier raises. Disabled by default.
- `alert_description_template` (String) Description (in Markdown) of the alert raised when the SLO enters this tier. Supports template variables. Leave empty to use the default description.
- `alert_remediation_notes` (String) Remediation notes (in Markdown) attached to the alert raised when the SLO enters this tier. Supports template variables.
- `alert_severity_id` (String) ID of the Alert Severity of the alert raised when the SLO enters this tier. The ID of a `oneuptime_alert_severity`.
- `alert_title_template` (String) Title of the alert raised when the SLO enters this tier. Supports template variables such as {{sloName}} and {{policyTier}}. Leave empty to use the default title.
- `auto_resolve_alert` (Boolean) Resolve the alert automatically when the error budget recovers out of this tier. Enabled by default. When disabled, the alert stays open until someone resolves it. Entering a more severe tier always resolves it.
- `auto_resolve_incident` (Boolean) Resolve the incident automatically when the error budget recovers out of this tier. Enabled by default. When disabled, the incident stays open until someone resolves it. Entering a more severe tier always resolves it.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `freezes_deployments` (Boolean) While the SLO is in this tier, its deploymentsAllowed field is false, so CI/CD pipelines that check it hold back releases. Disabled by default.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_description_template` (String) Description (in Markdown) of the incident declared when the SLO enters this tier. Supports template variables. Leave empty to use the default description.
- `incident_remediation_notes` (String) Remediation notes (in Markdown) attached to the incident declared when the SLO enters this tier. Supports template variables.
- `incident_severity_id` (String) ID of the Incident Severity of the incident declared when the SLO enters this tier. The ID of a `oneuptime_incident_severity`.
- `incident_title_template` (String) Title of the incident declared when the SLO enters this tier. Supports template variables such as {{sloName}} and {{policyTier}}. Leave empty to use the default title.
- `is_alert_private` (Boolean) Make the alert raised by this tier private, so only its owners, project admins and project owners can see it. Disabled by default.
- `is_enabled` (Boolean) Whether this tier is part of the policy. A disabled tier never applies, and anything it raised is resolved.
- `is_incident_private` (Boolean) Make the incident declared by this tier private, so only its owners, project admins and project owners can see it. Disabled by default.
- `name` (String) Name of this tier, for example "Reliability first". Shown on the SLO's Overview, in its feed and as {{policyTier}} in alert and incident templates.
- `service_level_objective_id` (String) ID of the Service Level Objective whose error budget policy this tier belongs to. The ID of a `oneuptime_service_level_objective`.
- `should_create_alert` (Boolean) Raise an Alert when the SLO enters this tier. Disabled by default.
- `should_create_incident` (Boolean) Declare an Incident when the SLO enters this tier. Disabled by default.
- `threshold_percentage` (Number) The tier applies while the share of the error budget left is at or below this percentage, from 0 to 100. 0 is the tier of an exhausted budget. The SLO is in the applying tier with the lowest threshold. Unique per SLO.
- `what_we_do` (String) What the team agreed to do while the SLO is in this tier, in Markdown - for example "Only reliability work and critical fixes ship." Shown on the SLO's Overview, in the feed item posted when the SLO enters the tier, and as {{policyTierWhatWeDo}} in alert and incident templates.

### Read-Only

- `alert_labels` (Set of String) Labels added to alerts raised by this tier. IDs of `oneuptime_label` resources.
- `alert_owner_teams` (Set of String) Teams added as owners of alerts raised by this tier. IDs of `oneuptime_team` resources.
- `alert_owner_users` (Set of String) Users added as owners of alerts raised by this tier. IDs of `oneuptime_user` records.
- `created_at` (String) Date and Time when the object was created.
- `incident_labels` (Set of String) Labels added to incidents declared by this tier. IDs of `oneuptime_label` resources.
- `incident_on_call_duty_policies` (Set of String) On-call policies attached to incidents declared by this tier. IDs of `oneuptime_on_call_policy` resources.
- `incident_owner_teams` (Set of String) Teams added as owners of incidents declared by this tier. IDs of `oneuptime_team` resources.
- `incident_owner_users` (Set of String) Users added as owners of incidents declared by this tier. IDs of `oneuptime_user` records.
- `last_alert_created_at` (String) The last time an alert was raised by this tier. Computed by the worker.
- `last_alert_resolved_at` (String) The last time an alert raised by this tier was resolved. Computed by the worker.
- `last_incident_created_at` (String) The last time an incident was declared by this tier. Computed by the worker.
- `last_incident_resolved_at` (String) The last time an incident declared by this tier was resolved. Computed by the worker.
- `on_call_duty_policies` (Set of String) On-call policies attached to alerts raised by this tier. Incidents have their own list. IDs of `oneuptime_on_call_policy` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
