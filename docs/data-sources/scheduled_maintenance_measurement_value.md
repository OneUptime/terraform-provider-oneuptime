---
page_title: "oneuptime_scheduled_maintenance_measurement_value Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  The computed value of one measurement for one scheduled maintenance event. Written by OneUptime, never edited by hand.
---

# oneuptime_scheduled_maintenance_measurement_value (Data Source)

The computed value of one measurement for one scheduled maintenance event. Written by OneUptime, never edited by hand. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_scheduled_maintenance_measurement_value" "by_name" {
  name = "example-scheduled_maintenance_measurement_value"
}

data "oneuptime_scheduled_maintenance_measurement_value" "by_id" {
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
- `scheduled_maintenance_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `scheduled_maintenance_measurement_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `started_at` (String) A date time object.. Computed.
- `ended_at` (String) A date time object.. Computed.
- `value_in_seconds` (Number) The measured duration in seconds. Empty unless the status is Recorded... Computed.
- `status` (String) Outcome of evaluating this measurement - Recorded, Pending, Not Applicable or Invalid.. Computed.
- `status_message` (String) Why this measurement has the status it has - which anchor is still open, or why it can never resolve.. Computed.
- `start_scheduled_maintenance_state_timeline_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `end_scheduled_maintenance_state_timeline_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `computed_at` (String) A date time object.. Computed.
