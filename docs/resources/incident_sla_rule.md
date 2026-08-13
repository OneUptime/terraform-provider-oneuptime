---
page_title: "oneuptime_incident_sla_rule Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Configure SLA rules to define response and resolution time targets for incidents
---

# oneuptime_incident_sla_rule (Resource)

Configure SLA rules to define response and resolution time targets for incidents

## Example Usage

```terraform
resource "oneuptime_incident_sla_rule" "example" {
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this SLA rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this SLA rule..
- `order` (Number) Order/priority of this rule. Rules are evaluated in order (lowest first). First matching rule wins...
- `is_enabled` (Bool) Whether this SLA rule is enabled..
- `response_time_in_minutes` (Number) Target response time in minutes. This is the maximum time allowed before the incident must be acknowledged...
- `resolution_time_in_minutes` (Number) Target resolution time in minutes. This is the maximum time allowed before the incident must be resolved...
- `at_risk_threshold_in_percentage` (Number) Percentage of the deadline at which the SLA status changes to At Risk. For example, 80 means the status becomes At Risk when 80% of the time has elapsed...
- `internal_note_reminder_interval_in_minutes` (Number) How often (in minutes) to automatically post internal notes to unresolved incidents. Internal notes are only visible to your team. For example, set to 30 to remind your team every 30 minutes to provide an update. Leave empty to disable...
- `public_note_reminder_interval_in_minutes` (Number) How often (in minutes) to automatically post public notes to unresolved incidents. Public notes are visible to external stakeholders on your status page. For example, set to 60 to post a status update every hour. Leave empty to disable...
- `internal_note_reminder_template` (String) The content of the automatic internal note posted to your team. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used...
- `public_note_reminder_template` (String) The content of the automatic public note shown on your status page. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used...
- `monitors` (Set) Only apply this SLA rule to incidents affecting these monitors. Leave empty to match incidents from any monitor...
- `incident_severities` (Set) Only apply this SLA rule to incidents with these severities. Leave empty to match incidents of any severity...
- `incident_labels` (Set) Only apply this SLA rule to incidents that have at least one of these labels. Leave empty to match incidents regardless of labels...
- `monitor_labels` (Set) Only apply this SLA rule to incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels...
- `incident_title_pattern` (String) Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'...
- `incident_description_pattern` (String) Regular expression pattern to match incident descriptions. Leave empty to match any description...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_sla_rule.example <id>
```
