---
page_title: "oneuptime_runbook_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created.
---

# oneuptime_runbook_rule (Data Source)

Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_runbook_rule" "by_name" {
  name = "example-runbook_rule"
}

data "oneuptime_runbook_rule" "by_id" {
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
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource... Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this runbook rule... Computed.
- `is_enabled` (Bool) Whether this rule is enabled... Computed.
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident, Alert, or ScheduledMaintenance... Computed.
- `monitors` (Set) Only match incidents and scheduled maintenance events that affect, and alerts raised by, at least one of these monitors. Leave empty to match any monitor... Computed.
- `incident_severities` (Set) Only match incidents with one of these severities. Incident rules only. Leave empty to match any severity... Computed.
- `alert_severities` (Set) Only match alerts with one of these severities. Alert rules only. Leave empty to match any severity... Computed.
- `labels` (Set) Only match incidents, alerts or scheduled maintenance events that carry at least one of these labels. Leave empty to match regardless of their labels... Computed.
- `monitor_labels` (Set) Only match when a monitor of the incident, alert or scheduled maintenance event carries at least one of these labels. Leave empty to match regardless of monitor labels... Computed.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title... Computed.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description... Computed.
- `monitor_name_pattern` (String) Case-insensitive regex matched against the names of the monitors of the incident, alert or scheduled maintenance event. Leave empty to match any monitor name... Computed.
- `monitor_description_pattern` (String) Case-insensitive regex matched against the descriptions of the monitors of the incident, alert or scheduled maintenance event. Leave empty to match any monitor description... Computed.
- `runbooks` (Set) Runbooks to start when this rule matches. Each runbook produces its own execution... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
