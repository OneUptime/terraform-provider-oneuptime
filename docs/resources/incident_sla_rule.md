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
  name        = "Example incident sla rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this SLA rule.

### Optional

- `at_risk_threshold_in_percentage` (Number) Percentage of the deadline at which the SLA status changes to At Risk. For example, 80 means the status becomes At Risk when 80% of the time has elapsed. Defaults to `80`.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this SLA rule.
- `incident_description_pattern` (String) Regular expression pattern to match incident descriptions. Leave empty to match any description.
- `incident_labels` (Set of String) Only apply this SLA rule to incidents that have at least one of these labels. Leave empty to match incidents regardless of labels. IDs of `oneuptime_label` resources.
- `incident_severities` (Set of String) Only apply this SLA rule to incidents with these severities. Leave empty to match incidents of any severity. IDs of `oneuptime_incident_severity` resources.
- `incident_title_pattern` (String) Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.
- `internal_note_reminder_interval_in_minutes` (Number) How often (in minutes) to automatically post internal notes to unresolved incidents. Internal notes are only visible to your team. For example, set to 30 to remind your team every 30 minutes to provide an update. Leave empty to disable.
- `internal_note_reminder_template` (String) The content of the automatic internal note posted to your team. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used.
- `is_enabled` (Boolean) Whether this SLA rule is enabled. Defaults to `true`.
- `monitor_labels` (Set of String) Only apply this SLA rule to incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only apply this SLA rule to incidents affecting these monitors. Leave empty to match incidents from any monitor. IDs of `oneuptime_monitor` resources.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `public_note_reminder_interval_in_minutes` (Number) How often (in minutes) to automatically post public notes to unresolved incidents. Public notes are visible to external stakeholders on your status page. For example, set to 60 to post a status update every hour. Leave empty to disable.
- `public_note_reminder_template` (String) The content of the automatic public note shown on your status page. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used.
- `resolution_time_in_minutes` (Number) Target resolution time in minutes. This is the maximum time allowed before the incident must be resolved.
- `response_time_in_minutes` (Number) Target response time in minutes. This is the maximum time allowed before the incident must be acknowledged.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident sla rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_sla_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_sla_rule.example <id>
```
