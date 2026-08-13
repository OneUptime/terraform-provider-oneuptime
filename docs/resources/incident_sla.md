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
  incident_id = "123e4567-e89b-12d3-a456-426614174000"
  incident_sla_rule_id = "123e4567-e89b-12d3-a456-426614174000"
  sla_started_at = "2030-01-01T00:00:00Z"
}
```

## Schema

### Required

- `incident_id` (String) A unique identifier for an object, represented as a UUID..
- `incident_sla_rule_id` (String) A unique identifier for an object, represented as a UUID..
- `sla_started_at` (String) A date time object..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `response_deadline` (String) A date time object..
- `resolution_deadline` (String) A date time object..
- `status` (String) Current SLA status (On Track, At Risk, Breached, Met)..
- `responded_at` (String) A date time object..
- `resolved_at` (String) A date time object..
- `last_internal_note_reminder_sent_at` (String) A date time object..
- `last_public_note_reminder_sent_at` (String) A date time object..
- `breach_notification_sent_at` (String) A date time object..
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
terraform import oneuptime_incident_sla.example <id>
```
