---
page_title: "oneuptime_incident_measurement Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  A named duration between two points in an incident's life, computed automatically for every incident
---

# oneuptime_incident_measurement (Data Source)

A named duration between two points in an incident's life, computed automatically for every incident Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_measurement" "by_name" {
  name = "example-incident_measurement"
}

data "oneuptime_incident_measurement" "by_id" {
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
- `key` (String) Stable, machine readable identifier for this measurement, unique within the project. It is immutable once created because it is used to build the metric name that every recorded point is written under - changing it would orphan all the history. Pick it carefully; to rename a measurement, change the Name instead... Computed.
- `description` (String) Description of what this measurement means to your team.. Computed.
- `metric_name` (String) The metric name every recorded point of this measurement is written under. Derived from the key as oneuptime.incident.measurement.<key> and maintained for you... Computed.
- `start_anchor_type` (String) Where this measurement starts. One of: Impact Started At, Declared At, Created At, Timeline Start, State Entered, State Role Entered, Postmortem Posted At... Computed.
- `end_anchor_type` (String) Where this measurement ends. One of: Impact Started At, Declared At, Created At, Timeline Start, State Entered, State Role Entered, Postmortem Posted At... Computed.
- `start_incident_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `end_incident_state_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `start_incident_state_role` (String) The role of the state this measurement starts at - Created, Acknowledged or Resolved. Used when the Start Anchor Type is State Role Entered. Resolving by role keeps the measurement working when a project renames or replaces the state that plays that part... Computed.
- `end_incident_state_role` (String) The role of the state this measurement ends at - Created, Acknowledged or Resolved. Used when the End Anchor Type is State Role Entered... Computed.
- `start_state_occurrence` (String) Which entry to use when the start state is entered more than once - First or Last. First matches the built-in incident metrics; Last follows a reopened incident to its final pass through that state... Computed.
- `end_state_occurrence` (String) Which entry to use when the end state is entered more than once - First or Last. First matches the built-in incident metrics; Last follows a reopened incident to its final pass through that state... Computed.
- `unit` (String) The unit this measurement's values are displayed in. Values are always stored in seconds; this only changes how they are rendered... Computed.
- `aggregation_type` (String) The aggregation this measurement's charts default to - Avg, Max, Min, P50, P90, P95 or P99. Sum is deliberately absent because summing durations across incidents produces a number with no meaning... Computed.
- `is_enabled` (Bool) Whether this measurement is computed for new and updated incidents.. Computed.
- `show_on_incident_view` (Bool) Whether this measurement is shown on the incident page alongside the incident's other timings.. Computed.
- `order` (Number) Order in which this measurement is displayed. Lowest first... Computed.
- `is_system_defined` (Bool) Whether this measurement was seeded by OneUptime rather than created by your team.. Computed.
- `backfill_requested_at` (String) A date time object.. Computed.
- `backfill_cursor_created_at` (String) A date time object.. Computed.
- `backfill_completed_at` (String) A date time object.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
