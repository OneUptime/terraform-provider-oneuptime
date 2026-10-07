---
page_title: "oneuptime_runbook_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created.
---

# oneuptime_runbook_rule (Resource)

Auto-attach runbooks to incidents, alerts, or scheduled maintenance events when they are created.

## Example Usage

```terraform
resource "oneuptime_runbook_rule" "example" {
  name = "Example short text"
  trigger_entity_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this runbook rule...
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident, Alert, or ScheduledMaintenance...

### Optional

- `criteria` (String) Versioned conditions that determine whether this rule matches a resource...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this runbook rule...
- `is_enabled` (Bool) Whether this rule is enabled...
- `monitors` (Set) Only match incidents and scheduled maintenance events that affect, and alerts raised by, at least one of these monitors. Leave empty to match any monitor...
- `incident_severities` (Set) Only match incidents with one of these severities. Incident rules only. Leave empty to match any severity...
- `alert_severities` (Set) Only match alerts with one of these severities. Alert rules only. Leave empty to match any severity...
- `labels` (Set) Only match incidents, alerts or scheduled maintenance events that carry at least one of these labels. Leave empty to match regardless of their labels...
- `monitor_labels` (Set) Only match when a monitor of the incident, alert or scheduled maintenance event carries at least one of these labels. Leave empty to match regardless of monitor labels...
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title...
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description...
- `monitor_name_pattern` (String) Case-insensitive regex matched against the names of the monitors of the incident, alert or scheduled maintenance event. Leave empty to match any monitor name...
- `monitor_description_pattern` (String) Case-insensitive regex matched against the descriptions of the monitors of the incident, alert or scheduled maintenance event. Leave empty to match any monitor description...
- `runbooks` (Set) Runbooks to start when this rule matches. Each runbook produces its own execution...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_runbook_rule.example <id>
```
