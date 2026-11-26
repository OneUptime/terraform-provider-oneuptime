---
page_title: "oneuptime_incident_sla_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure SLA rules to define response and resolution time targets for incidents
---

# oneuptime_incident_sla_rule (Data Source)

Configure SLA rules to define response and resolution time targets for incidents

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident sla rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_sla_rule" "example" {
  name = "Example incident sla rule"
}

# Or by id:
data "oneuptime_incident_sla_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `at_risk_threshold_in_percentage` (Number) Percentage of the deadline at which the SLA status changes to At Risk. For example, 80 means the status becomes At Risk when 80% of the time has elapsed.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this SLA rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_description_pattern` (String) Regular expression pattern to match incident descriptions. Leave empty to match any description.
- `incident_title_pattern` (String) Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.
- `internal_note_reminder_interval_in_minutes` (Number) How often (in minutes) to automatically post internal notes to unresolved incidents. Only the people who can see an incident can read its internal notes. For example, set to 30 to remind your team every 30 minutes to provide an update. Leave empty to disable.
- `internal_note_reminder_template` (String) The content of the automatic internal note posted on the incident. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used.
- `is_enabled` (Boolean) Whether this SLA rule is enabled.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Name of this SLA rule.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them.
- `public_note_reminder_interval_in_minutes` (Number) How often (in minutes) to automatically post public notes to unresolved incidents. Public notes are visible to external stakeholders on your status page. For example, set to 60 to post a status update every hour. Leave empty to disable.
- `public_note_reminder_template` (String) The content of the automatic public note shown on your status page. Use variables like {{incidentTitle}}, {{elapsedTime}}, {{slaStatus}}, {{timeToResolutionDeadline}} to include dynamic incident data. If left empty, a default template will be used.
- `resolution_time_in_minutes` (Number) Target resolution time in minutes. This is the maximum time allowed before the incident must be resolved.
- `response_time_in_minutes` (Number) Target response time in minutes. This is the maximum time allowed before the incident must be acknowledged.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: read it with `jsondecode()`.
- `emptied_limits` (String) The lists this rule was limited to whose every record was deleted since: each list, the deleted records' names, how many and when. Deleting the last record a list held would make the rule match every record instead, so OneUptime switched the rule off and wrote this. It is cleared when the rule is switched on again; until a list named here holds a record again, switching the rule on must say what the list holds in the same update. Written by OneUptime; any value sent for it is refused. A JSON value: read it with `jsondecode()`.
- `incident_labels` (Set of String) Only apply this SLA rule to incidents that have at least one of these labels. Leave empty to match incidents regardless of labels. IDs of `oneuptime_label` resources.
- `incident_severities` (Set of String) Only apply this SLA rule to incidents with these severities. Leave empty to match incidents of any severity. IDs of `oneuptime_incident_severity` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `monitor_labels` (Set of String) Only apply this SLA rule to incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only apply this SLA rule to incidents affecting these monitors. Leave empty to match incidents from any monitor. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
