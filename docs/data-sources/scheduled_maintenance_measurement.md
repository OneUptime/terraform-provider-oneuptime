---
page_title: "oneuptime_scheduled_maintenance_measurement Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  A named duration between two points in a scheduled maintenance event's life, computed automatically for every event
---

# oneuptime_scheduled_maintenance_measurement (Data Source)

A named duration between two points in a scheduled maintenance event's life, computed automatically for every event

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled maintenance measurement may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_measurement" "example" {
  name = "Example scheduled maintenance measurement"
}

# Or by id:
data "oneuptime_scheduled_maintenance_measurement" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `aggregation_type` (String) How this measurement's chart sums up many scheduled maintenance events by default - Avg (the default), P50, P90, P95, P99, Max or Min. View Chart in the dashboard opens the chart this way. Sum is deliberately absent: adding durations up across scheduled maintenance events produces a number with no meaning.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of what this measurement means to your team.
- `end_anchor_type` (String) Where the measurement ends - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role.
- `end_scheduled_maintenance_state_id` (String) ID of the state whose entry ends this measurement, when the end anchor is a specific state. The ID of a `oneuptime_scheduled_maintenance_state`.
- `end_scheduled_maintenance_state_role` (String) The role of the state that ends this measurement (Scheduled, Ongoing, Ended or Resolved), when the end anchor is a state role. Resolving by role keeps working when a project renames or replaces the state.
- `end_state_occurrence` (String) Which entry to use when the end state is entered more than once - the first time it was entered, or the last.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this measurement is computed for scheduled maintenance events.
- `is_system_defined` (Boolean) Whether this measurement was created by OneUptime rather than by your team.
- `key` (String) Stable, machine readable identifier for this measurement, unique within the project: lowercase letters, numbers and hyphens. Leave it out and it is made from the name - Time to Start becomes time-to-start, with -2, -3 and so on added when another measurement already has it. It is part of the metric name, so it cannot be changed once the measurement is created; to rename a measurement, change the Name instead.
- `metric_name` (String) Name of the metric this measurement writes to, derived from the key as oneuptime.scheduled-maintenance.measurement.<key>.
- `name` (String) Human readable name of this measurement, such as Start Delay. This is what charts call it.
- `order` (Number) Where this measurement appears in the list of measurements, lowest number first. A new measurement is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `show_on_scheduled_maintenance_view` (Boolean) Whether this measurement is shown on the scheduled maintenance event page.
- `start_anchor_type` (String) Where the measurement starts - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role.
- `start_scheduled_maintenance_state_id` (String) ID of the state whose entry starts this measurement, when the start anchor is a specific state. The ID of a `oneuptime_scheduled_maintenance_state`.
- `start_scheduled_maintenance_state_role` (String) The role of the state that starts this measurement (Scheduled, Ongoing, Ended or Resolved), when the start anchor is a state role. Resolving by role keeps working when a project renames or replaces the state.
- `start_state_occurrence` (String) Which entry to use when the start state is entered more than once - the first time it was entered, or the last.
- `unit` (String) The unit this measurement's charts are in: seconds (the default), minutes, hours or days. Every value is worked out in seconds and stored that way on the scheduled maintenance event; each chart point is written in this unit, so a chart in hours reads 1.5 for an hour and a half. With seconds, charts show seconds, minutes, hours or days as the numbers grow. Changing it rewrites the measurement's chart points in the new unit. A value that is not a time unit charts in seconds.

### Read-Only

- `backfill_completed_at` (String) When the backfill of this measurement over existing scheduled maintenance events finished.
- `backfill_cursor_created_at` (String) How far the backfill has walked this project, so a restart resumes instead of starting over.
- `backfill_requested_at` (String) When a backfill of this measurement over existing scheduled maintenance events was requested.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
