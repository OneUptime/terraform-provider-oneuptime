---
page_title: "oneuptime_incident_measurement_value Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  The computed value of one incident measurement for one incident, recomputed from the incident's timeline rather than accumulated
---

# oneuptime_incident_measurement_value (Data Source)

The computed value of one incident measurement for one incident, recomputed from the incident's timeline rather than accumulated

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident measurement value may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_measurement_value" "example" {
  incident_id = oneuptime_incident.example.id
}

# Or by id:
data "oneuptime_incident_measurement_value" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `end_incident_state_timeline_id` (String) The incident state timeline entry the end anchor resolved to. Recorded for provenance only - it carries no foreign key, so deleting a timeline entry never blocks or rewrites this row; the next recompute simply produces the right answer.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this measurement value was computed for. The ID of a `oneuptime_incident`.
- `incident_measurement_id` (String) ID of the measurement definition this value was computed from. The ID of a `oneuptime_incident_measurement`.
- `start_incident_state_timeline_id` (String) The incident state timeline entry the start anchor resolved to. Recorded for provenance only - it carries no foreign key, so deleting a timeline entry never blocks or rewrites this row; the next recompute simply produces the right answer.
- `status` (String) The outcome of evaluating this measurement: Recorded, Pending, Not Applicable or Invalid.
- `status_message` (String) Why this measurement has the status it has, in plain words - for example which anchor has not been reached yet, or by how much the end precedes the start.
- `value_in_seconds` (Number) The measured duration in seconds. Only set when the status is Recorded - a measurement that could not be computed is left blank rather than written as zero.

### Read-Only

- `computed_at` (String) When this value was last recomputed.
- `created_at` (String) Date and Time when the object was created.
- `ended_at` (String) When this measurement's end anchor resolved to. Blank while the end anchor has not resolved.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `started_at` (String) When this measurement's start anchor resolved to. Blank while the start anchor has not resolved.
- `updated_at` (String) Date and Time when the object was updated.
