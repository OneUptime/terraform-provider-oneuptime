---
page_title: "oneuptime_runbook_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created.
---

# oneuptime_runbook_rule (Data Source)

Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one runbook rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_runbook_rule" "example" {
  name = "Example runbook rule"
}

# Or by id:
data "oneuptime_runbook_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this runbook rule.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `monitor_description_pattern` (String) Case-insensitive regex matched against the descriptions of the monitors of the incident, alert or scheduled maintenance event. Leave empty to match any monitor description.
- `monitor_name_pattern` (String) Case-insensitive regex matched against the names of the monitors of the incident, alert or scheduled maintenance event. Leave empty to match any monitor name.
- `name` (String) Name of this runbook rule.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title.
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident, Alert, or ScheduledMaintenance.

### Read-Only

- `alert_severities` (Set of String) Only match alerts with one of these severities. Alert rules only. Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.
- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `incident_severities` (Set of String) Only match incidents with one of these severities. Incident rules only. Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.
- `labels` (Set of String) Only match incidents, alerts or scheduled maintenance events that carry at least one of these labels. Leave empty to match regardless of their labels. IDs of `oneuptime_label` resources.
- `monitor_labels` (Set of String) Only match when a monitor of the incident, alert or scheduled maintenance event carries at least one of these labels. Leave empty to match regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only match incidents and scheduled maintenance events that affect, and alerts raised by, at least one of these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `runbooks` (Set of String) Runbooks to start when this rule matches. Each runbook produces its own execution. IDs of `oneuptime_runbook` resources.
- `updated_at` (String) Date and Time when the object was updated.
