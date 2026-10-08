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
  incident_id = oneuptime_incident.example.id
}

# Or by id:
data "oneuptime_incident_sla" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this SLA record is tracking. The ID of a `oneuptime_incident`.
- `incident_sla_rule_id` (String) ID of the SLA rule that was applied to this incident. The ID of a `oneuptime_incident_sla_rule`.
- `status` (String) Current SLA status (On Track, At Risk, Breached, Met).

### Read-Only

- `breach_notification_sent_at` (String) The time when breach notification was sent to incident owners.
- `created_at` (String) Date and Time when the object was created.
- `last_internal_note_reminder_sent_at` (String) The last time an internal note reminder was sent.
- `last_public_note_reminder_sent_at` (String) The last time a public note reminder was sent.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `resolution_deadline` (String) The deadline by which the incident must be resolved to meet the SLA.
- `resolved_at` (String) The actual time when the incident was resolved.
- `responded_at` (String) The actual time when the incident was acknowledged.
- `response_deadline` (String) The deadline by which the incident must be acknowledged to meet the SLA.
- `sla_started_at` (String) The time when SLA tracking started (usually the incident declaredAt time).
- `updated_at` (String) Date and Time when the object was updated.
