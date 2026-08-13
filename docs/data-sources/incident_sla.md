---
page_title: "oneuptime_incident_sla Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Track SLA status and deadlines for incidents
---

# oneuptime_incident_sla (Data Source)

Track SLA status and deadlines for incidents Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_sla" "by_name" {
  name = "example-incident_sla"
}

data "oneuptime_incident_sla" "by_id" {
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
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_sla_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `response_deadline` (String) A date time object.. Computed.
- `resolution_deadline` (String) A date time object.. Computed.
- `status` (String) Current SLA status (On Track, At Risk, Breached, Met).. Computed.
- `responded_at` (String) A date time object.. Computed.
- `resolved_at` (String) A date time object.. Computed.
- `last_internal_note_reminder_sent_at` (String) A date time object.. Computed.
- `last_public_note_reminder_sent_at` (String) A date time object.. Computed.
- `breach_notification_sent_at` (String) A date time object.. Computed.
- `sla_started_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
