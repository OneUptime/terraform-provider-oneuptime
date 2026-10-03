---
page_title: "oneuptime_scheduled_maintenance_measurement Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  A named duration between two points in a scheduled maintenance event's life, computed automatically for every event
---

# oneuptime_scheduled_maintenance_measurement (Data Source)

A named duration between two points in a scheduled maintenance event's life, computed automatically for every event Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_scheduled_maintenance_measurement" "by_name" {
  name = "example-scheduled_maintenance_measurement"
}

data "oneuptime_scheduled_maintenance_measurement" "by_id" {
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
- `key` (String) Stable, machine readable identifier for this measurement, unique within the project: lowercase letters, numbers and hyphens. Leave it out and it is made from the name - Time to Start becomes time-to-start, with -2, -3 and so on added when another measurement already has it. It is part of the metric name, so it cannot be changed once the measurement is created; to rename a measurement, change the Name instead... Computed.
- `description` (String) Description of what this measurement means to your team.. Computed.
- `metric_name` (String) Name of the metric this measurement writes to, derived from the key as oneuptime.scheduled-maintenance.measurement.<key>.. Computed.
- `start_anchor_type` (String) Where the measurement starts - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role... Computed.
- `end_anchor_type` (String) Where the measurement ends - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role... Computed.
- `start_scheduled_maintenance_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `end_scheduled_maintenance_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `start_scheduled_maintenance_state_role` (String) The role of the state that starts this measurement (Scheduled, Ongoing, Ended or Resolved), when the start anchor is a state role. Resolving by role keeps working when a project renames or replaces the state... Computed.
- `end_scheduled_maintenance_state_role` (String) The role of the state that ends this measurement (Scheduled, Ongoing, Ended or Resolved), when the end anchor is a state role. Resolving by role keeps working when a project renames or replaces the state... Computed.
- `start_state_occurrence` (String) Which entry to use when the start state is entered more than once - the first time it was entered, or the last... Computed.
- `end_state_occurrence` (String) Which entry to use when the end state is entered more than once - the first time it was entered, or the last... Computed.
- `unit` (String) The unit this measurement's charts are in: seconds (the default), minutes, hours or days. Every value is worked out in seconds and stored that way on the scheduled maintenance event; each chart point is written in this unit, so a chart in hours reads 1.5 for an hour and a half. With seconds, charts show seconds, minutes, hours or days as the numbers grow. Changing it rewrites the measurement's chart points in the new unit. A value that is not a time unit charts in seconds... Computed.
- `aggregation_type` (String) How this measurement's chart sums up many scheduled maintenance events by default - Avg (the default), P50, P90, P95, P99, Max or Min. View Chart in the dashboard opens the chart this way. Sum is deliberately absent: adding durations up across scheduled maintenance events produces a number with no meaning... Computed.
- `is_enabled` (Bool) Whether this measurement is computed for scheduled maintenance events.. Computed.
- `show_on_scheduled_maintenance_view` (Bool) Whether this measurement is shown on the scheduled maintenance event page.. Computed.
- `order` (Number) Where this measurement appears in the list of measurements, lowest number first. A new measurement is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them... Computed.
- `is_system_defined` (Bool) Whether this measurement was created by OneUptime rather than by your team.. Computed.
- `backfill_requested_at` (String) A date time object.. Computed.
- `backfill_cursor_created_at` (String) A date time object.. Computed.
- `backfill_completed_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
