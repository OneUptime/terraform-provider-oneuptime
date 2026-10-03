---
page_title: "oneuptime_alert_measurement Resource - oneuptime"
subcategory: "Alerts"
description: |-
  A named duration between two points in an alert's life, computed automatically for every alert
---

# oneuptime_alert_measurement (Resource)

A named duration between two points in an alert's life, computed automatically for every alert

## Example Usage

```terraform
resource "oneuptime_alert_measurement" "example" {
  name = "Example short text"
  start_anchor_type = "Example short text"
  end_anchor_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Human readable name of this measurement, such as Time to Acknowledge. This is what charts call it...
- `start_anchor_type` (String) Where this measurement starts. One of: Impact Started At, Created At, Timeline Start, State Entered, State Role Entered...
- `end_anchor_type` (String) Where this measurement ends. One of: Impact Started At, Created At, Timeline Start, State Entered, State Role Entered...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `key` (String) Stable, machine readable identifier for this measurement, unique within the project: lowercase letters, numbers and hyphens. Leave it out and it is made from the name - Time to Acknowledge becomes time-to-acknowledge, with -2, -3 and so on added when another measurement already has it. It cannot be changed once the measurement is created, because it is used to build the metric name that every recorded point is written under; to rename a measurement, change the Name instead...
- `description` (String) Description of what this measurement means to your team..
- `start_alert_state_id` (String) A unique identifier for an object, represented as a UUID..
- `end_alert_state_id` (String) A unique identifier for an object, represented as a UUID..
- `start_alert_state_role` (String) The role of the state this measurement starts at - Created, Acknowledged or Resolved. Used when the Start Anchor Type is State Role Entered. Resolving by role keeps the measurement working when a project renames or replaces the state that plays that part...
- `end_alert_state_role` (String) The role of the state this measurement ends at - Created, Acknowledged or Resolved. Used when the End Anchor Type is State Role Entered...
- `start_state_occurrence` (String) Which entry to use when the start state is entered more than once - First or Last. First matches the built-in alert metrics; Last follows a reopened alert to its final pass through that state...
- `end_state_occurrence` (String) Which entry to use when the end state is entered more than once - First or Last. First matches the built-in alert metrics; Last follows a reopened alert to its final pass through that state...
- `unit` (String) The unit this measurement's charts are in: seconds (the default), minutes, hours or days. Every value is worked out in seconds and stored that way on the alert; each chart point is written in this unit, so a chart in hours reads 1.5 for an hour and a half. With seconds, charts show seconds, minutes, hours or days as the numbers grow. Changing it rewrites the measurement's chart points in the new unit. A value that is not a time unit charts in seconds...
- `aggregation_type` (String) How this measurement's chart sums up many alerts by default - Avg (the default), P50, P90, P95, P99, Max or Min. View Chart in the dashboard opens the chart this way. Sum is deliberately absent: adding durations up across alerts produces a number with no meaning...
- `is_enabled` (Bool) Whether this measurement is computed for new and updated alerts..
- `show_on_alert_view` (Bool) Whether this measurement is shown on the alert page alongside the alert's other timings..
- `order` (Number) Where this measurement appears in the list of measurements, lowest number first. A new measurement is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `metric_name` (String) The metric name every recorded point of this measurement is written under. Derived from the key as oneuptime.alert.measurement.<key> and maintained for you...
- `is_system_defined` (Bool) Whether this measurement was seeded by OneUptime rather than created by your team..
- `backfill_requested_at` (String) A date time object..
- `backfill_cursor_created_at` (String) A date time object..
- `backfill_completed_at` (String) A date time object..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_alert_measurement.example <id>
```
