---
page_title: "oneuptime_incident_sla Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Track SLA status and deadlines for incidents
---

# oneuptime_incident_sla (Data Source)

Track SLA status and deadlines for incidents

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident sla may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_sla" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_incident_sla" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this SLA record is tracking. The ID of a `oneuptime_incident`.
- `incident_sla_rule_id` (String) ID of the SLA rule that was applied to this incident. The ID of a `oneuptime_incident_sla_rule`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `status` (String) Current SLA status (On Track, At Risk, Breached, Met).

### Read-Only

- `breach_notification_sent_at` (String) The time when breach notification was sent to incident owners.
- `created_at` (String) Date and Time when the object was created.
- `last_internal_note_reminder_sent_at` (String) The last time an internal note reminder was sent.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_public_note_reminder_sent_at` (String) The last time a public note reminder was sent.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `resolution_deadline` (String) The deadline by which the incident must be resolved to meet the SLA.
- `resolved_at` (String) The actual time when the incident was resolved.
- `responded_at` (String) The actual time when the incident was acknowledged.
- `response_deadline` (String) The deadline by which the incident must be acknowledged to meet the SLA.
- `sla_started_at` (String) The time when SLA tracking started (usually the incident declaredAt time).
- `updated_at` (String) Date and Time when the object was updated.
