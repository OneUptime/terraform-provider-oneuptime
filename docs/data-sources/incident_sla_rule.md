---
page_title: "oneuptime_incident_sla_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure SLA rules to define response and resolution time targets for incidents
---

# oneuptime_incident_sla_rule (Data Source)

Configure SLA rules to define response and resolution time targets for incidents Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_sla_rule" "by_name" {
  name = "example-incident_sla_rule"
}

data "oneuptime_incident_sla_rule" "by_id" {
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
- `description` (String) Description of this SLA rule.. Computed.
- `order` (Number) Order/priority of this rule. Rules are evaluated in order (lowest first). First matching rule wins... Computed.
- `is_enabled` (Bool) Whether this SLA rule is enabled.. Computed.
- `response_time_in_minutes` (Number) Target response time in minutes. This is the maximum time allowed before the incident must be acknowledged... Computed.
- `resolution_time_in_minutes` (Number) Target resolution time in minutes. This is the maximum time allowed before the incident must be resolved... Computed.
- `at_risk_threshold_in_percentage` (Number) Percentage of the deadline at which the SLA status changes to At Risk. For example, 80 means the status becomes At Risk when 80% of the time has elapsed... Computed.
- `internal_note_reminder_interval_in_minutes` (Number) How often (in minutes) to automatically post internal notes to unresolved incidents. Internal notes are only visible to your team. For example, set to 30 to remind your team every 30 minutes to provide an update. Leave empty to disable... Computed.
- `public_note_reminder_interval_in_minutes` (Number) How often (in minutes) to automatically post public notes to unresolved incidents. Public notes are visible to external stakeholders on your status page. For example, set to 60 to post a status update every hour. Leave empty to disable... Computed.
- `internal_note_reminder_template` (String) The content of the automatic internal note posted to your team. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used... Computed.
- `public_note_reminder_template` (String) The content of the automatic public note shown on your status page. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used... Computed.
- `monitors` (Set) Only apply this SLA rule to incidents affecting these monitors. Leave empty to match incidents from any monitor... Computed.
- `incident_severities` (Set) Only apply this SLA rule to incidents with these severities. Leave empty to match incidents of any severity... Computed.
- `incident_labels` (Set) Only apply this SLA rule to incidents that have at least one of these labels. Leave empty to match incidents regardless of labels... Computed.
- `monitor_labels` (Set) Only apply this SLA rule to incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels... Computed.
- `incident_title_pattern` (String) Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'... Computed.
- `incident_description_pattern` (String) Regular expression pattern to match incident descriptions. Leave empty to match any description... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
