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
  key = "Example short text"
  start_anchor_type = "Example short text"
  end_anchor_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Human readable name of this measurement. This is what appears on charts and on the alert page...
- `key` (String) Stable, machine readable identifier for this measurement, unique within the project. It is immutable once created because it is used to build the metric name that every recorded point is written under - changing it would orphan all the history. Pick it carefully; to rename a measurement, change the Name instead...
- `start_anchor_type` (String) Where this measurement starts. One of: Impact Started At, Created At, Timeline Start, State Entered, State Role Entered...
- `end_anchor_type` (String) Where this measurement ends. One of: Impact Started At, Created At, Timeline Start, State Entered, State Role Entered...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of what this measurement means to your team..
- `start_alert_state_id` (String) A unique identifier for an object, represented as a UUID..
- `end_alert_state_id` (String) A unique identifier for an object, represented as a UUID..
- `start_alert_state_role` (String) The role of the state this measurement starts at - Created, Acknowledged or Resolved. Used when the Start Anchor Type is State Role Entered. Resolving by role keeps the measurement working when a project renames or replaces the state that plays that part...
- `end_alert_state_role` (String) The role of the state this measurement ends at - Created, Acknowledged or Resolved. Used when the End Anchor Type is State Role Entered...
- `start_state_occurrence` (String) Which entry to use when the start state is entered more than once - First or Last. First matches the built-in alert metrics; Last follows a reopened alert to its final pass through that state...
- `end_state_occurrence` (String) Which entry to use when the end state is entered more than once - First or Last. First matches the built-in alert metrics; Last follows a reopened alert to its final pass through that state...
- `unit` (String) The unit this measurement's values are displayed in. Values are always stored in seconds; this only changes how they are rendered...
- `aggregation_type` (String) The aggregation this measurement's charts default to - Avg, Max, Min, P50, P90, P95 or P99. Sum is deliberately absent because summing durations across alerts produces a number with no meaning...
- `is_enabled` (Bool) Whether this measurement is computed for new and updated alerts..
- `show_on_alert_view` (Bool) Whether this measurement is shown on the alert page alongside the alert's other timings..
- `order` (Number) Order in which this measurement is displayed. Lowest first...
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
