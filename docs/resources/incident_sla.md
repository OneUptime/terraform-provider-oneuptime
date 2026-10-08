---
page_title: "oneuptime_incident_sla Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Track SLA status and deadlines for incidents
---

# oneuptime_incident_sla (Resource)

Track SLA status and deadlines for incidents

## Example Usage

```terraform
resource "oneuptime_incident_sla" "example" {
  incident_id          = oneuptime_incident.example.id
  incident_sla_rule_id = oneuptime_incident_sla_rule.example.id
  sla_started_at       = "2030-01-01T00:00:00Z"
}
```

## Schema

### Required

- `incident_id` (String) ID of the incident this SLA record is tracking. The ID of a `oneuptime_incident`.
- `incident_sla_rule_id` (String) ID of the SLA rule that was applied to this incident. The ID of a `oneuptime_incident_sla_rule`.
- `sla_started_at` (String) The time when SLA tracking started (usually the incident declaredAt time).

### Optional

- `breach_notification_sent_at` (String) The time when breach notification was sent to incident owners.
- `last_internal_note_reminder_sent_at` (String) The last time an internal note reminder was sent.
- `last_public_note_reminder_sent_at` (String) The last time a public note reminder was sent.
- `resolution_deadline` (String) The deadline by which the incident must be resolved to meet the SLA.
- `resolved_at` (String) The actual time when the incident was resolved.
- `responded_at` (String) The actual time when the incident was acknowledged.
- `response_deadline` (String) The deadline by which the incident must be acknowledged to meet the SLA.
- `status` (String) Current SLA status (On Track, At Risk, Breached, Met). Defaults to `On Track`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident sla by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_sla.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_sla.example <id>
```
