---
page_title: "oneuptime_scheduled_maintenance_measurement Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  A named duration between two points in a scheduled maintenance event's life, computed automatically for every event
---

# oneuptime_scheduled_maintenance_measurement (Resource)

A named duration between two points in a scheduled maintenance event's life, computed automatically for every event

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_measurement" "example" {
  name = "Example short text"
  key = "Example short text"
  start_anchor_type = "Example short text"
  end_anchor_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this measurement, shown on charts and on the scheduled maintenance event page..
- `key` (String) Stable machine-readable key for this measurement. It is part of the metric name, so it cannot be changed once the measurement is created...
- `start_anchor_type` (String) Where the measurement starts - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role...
- `end_anchor_type` (String) Where the measurement ends - the moment the event was created, either end of the planned window, the start of its timeline, a specific state, or a state role...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of what this measurement means to your team..
- `start_scheduled_maintenance_state_id` (String) A unique identifier for an object, represented as a UUID..
- `end_scheduled_maintenance_state_id` (String) A unique identifier for an object, represented as a UUID..
- `start_scheduled_maintenance_state_role` (String) The role of the state that starts this measurement (Scheduled, Ongoing, Ended or Resolved), when the start anchor is a state role. Resolving by role keeps working when a project renames or replaces the state...
- `end_scheduled_maintenance_state_role` (String) The role of the state that ends this measurement (Scheduled, Ongoing, Ended or Resolved), when the end anchor is a state role. Resolving by role keeps working when a project renames or replaces the state...
- `start_state_occurrence` (String) Which entry to use when the start state is entered more than once - the first time it was entered, or the last...
- `end_state_occurrence` (String) Which entry to use when the end state is entered more than once - the first time it was entered, or the last...
- `unit` (String) The unit this measurement is displayed in. Values are always stored in seconds...
- `aggregation_type` (String) The aggregation this measurement's chart defaults to. Summing durations across events produces a number with no meaning, so Sum is not offered...
- `is_enabled` (Bool) Whether this measurement is computed for scheduled maintenance events..
- `show_on_scheduled_maintenance_view` (Bool) Whether this measurement is shown on the scheduled maintenance event page..
- `order` (Number) Order in which this measurement is displayed. Lowest first...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `metric_name` (String) Name of the metric this measurement writes to, derived from the key as oneuptime.scheduled-maintenance.measurement.<key>..
- `is_system_defined` (Bool) Whether this measurement was created by OneUptime rather than by your team..
- `backfill_requested_at` (String) A date time object..
- `backfill_cursor_created_at` (String) A date time object..
- `backfill_completed_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_scheduled_maintenance_measurement.example <id>
```
