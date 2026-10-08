---
page_title: "oneuptime_incident_measurement Resource - oneuptime"
subcategory: "Incidents"
description: |-
  A named duration between two points in an incident's life, computed automatically for every incident
---

# oneuptime_incident_measurement (Resource)

A named duration between two points in an incident's life, computed automatically for every incident

## Example Usage

```terraform
resource "oneuptime_incident_measurement" "example" {
  name              = "Example incident measurement"
  start_anchor_type = "Example short text"
  end_anchor_type   = "Example short text"
  description       = "Managed by Terraform"
}
```

## Schema

### Required

- `end_anchor_type` (String) Where this measurement ends. One of: Impact Started At, Declared At, Created At, Timeline Start, State Entered, State Role Entered, Postmortem Posted At.
- `name` (String) Human readable name of this measurement, such as Time to Acknowledge. This is what charts call it.
- `start_anchor_type` (String) Where this measurement starts. One of: Impact Started At, Declared At, Created At, Timeline Start, State Entered, State Role Entered, Postmortem Posted At.

### Optional

- `aggregation_type` (String) How this measurement's chart sums up many incidents by default - Avg (the default), P50, P90, P95, P99, Max or Min. View Chart in the dashboard opens the chart this way. Sum is deliberately absent: adding durations up across incidents produces a number with no meaning. Defaults to `Avg`.
- `description` (String) Description of what this measurement means to your team.
- `end_incident_state_id` (String) ID of the incident state this measurement ends at. Required only when the End Anchor Type is State Entered. Cleared if that state is deleted, at which point the measurement reports Not Applicable rather than a wrong number. The ID of a `oneuptime_incident_state`.
- `end_incident_state_role` (String) The role of the state this measurement ends at - Created, Acknowledged or Resolved. Used when the End Anchor Type is State Role Entered.
- `end_state_occurrence` (String) Which entry to use when the end state is entered more than once - First or Last. First matches the built-in incident metrics; Last follows a reopened incident to its final pass through that state. Defaults to `First`.
- `is_enabled` (Boolean) Whether this measurement is computed for new and updated incidents. Defaults to `true`.
- `key` (String) Stable, machine readable identifier for this measurement, unique within the project: lowercase letters, numbers and hyphens. Leave it out and it is made from the name - Time to Detect becomes time-to-detect, with -2, -3 and so on added when another measurement already has it. It cannot be changed once the measurement is created, because it is used to build the metric name that every recorded point is written under; to rename a measurement, change the Name instead.
- `order` (Number) Where this measurement appears in the list of measurements, lowest number first. A new measurement is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `show_on_incident_view` (Boolean) Whether this measurement is shown on the incident page alongside the incident's other timings. Defaults to `true`.
- `start_incident_state_id` (String) ID of the incident state this measurement starts at. Required only when the Start Anchor Type is State Entered. Cleared if that state is deleted, at which point the measurement reports Not Applicable rather than a wrong number. The ID of a `oneuptime_incident_state`.
- `start_incident_state_role` (String) The role of the state this measurement starts at - Created, Acknowledged or Resolved. Used when the Start Anchor Type is State Role Entered. Resolving by role keeps the measurement working when a project renames or replaces the state that plays that part.
- `start_state_occurrence` (String) Which entry to use when the start state is entered more than once - First or Last. First matches the built-in incident metrics; Last follows a reopened incident to its final pass through that state. Defaults to `First`.
- `unit` (String) The unit this measurement's charts are in: seconds (the default), minutes, hours or days. Every value is worked out in seconds and stored that way on the incident; each chart point is written in this unit, so a chart in hours reads 1.5 for an hour and a half. With seconds, charts show seconds, minutes, hours or days as the numbers grow. Changing it rewrites the measurement's chart points in the new unit. A value that is not a time unit charts in seconds. Defaults to `seconds`.

### Read-Only

- `backfill_completed_at` (String) When the backfill of this measurement over existing incidents finished.
- `backfill_cursor_created_at` (String) How far the backfill has walked this project, so a restart resumes instead of starting over.
- `backfill_requested_at` (String) When a backfill of this measurement over existing incidents was requested.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_system_defined` (Boolean) Whether this measurement was seeded by OneUptime rather than created by your team.
- `metric_name` (String) The metric name every recorded point of this measurement is written under. Derived from the key as oneuptime.incident.measurement.<key> and maintained for you.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident measurement by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_measurement.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_measurement.example <id>
```
