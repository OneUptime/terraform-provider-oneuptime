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
  name              = "Example alert measurement"
  start_anchor_type = "Impact Started At"
  end_anchor_type   = "Impact Started At"
  description       = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this alert measurement's place among the alert measurements of its project, lowest first. Setting it places this alert measurement and moves no other, so a configuration that numbers every alert measurement of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.alert_measurements`:

```terraform
order = index(local.alert_measurements, each.key) + 1
```

- Give each alert measurement a number of its own: two alert measurements at one number are shown the older first.
- An alert measurement created without `order` goes after the last one of its project.
- Deleting an alert measurement leaves a gap in the numbers, which changes nothing: alert measurements go by their numbers, not by counting.
- Moving an alert measurement in the OneUptime Dashboard numbers the project's alert measurements 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `end_anchor_type` (String) Where this measurement ends. One of: Impact Started At, Created At, Timeline Start, State Entered, State Role Entered, First Paged At. Allowed values: `Impact Started At`, `Created At`, `Timeline Start`, `State Entered`, `State Role Entered`, `First Paged At`.
- `name` (String) Human readable name of this measurement, such as Time to Acknowledge. This is what charts call it.
- `start_anchor_type` (String) Where this measurement starts. One of: Impact Started At, Created At, Timeline Start, State Entered, State Role Entered, First Paged At. Allowed values: `Impact Started At`, `Created At`, `Timeline Start`, `State Entered`, `State Role Entered`, `First Paged At`.

### Optional

- `aggregation_type` (String) How this measurement's chart sums up many alerts by default - Avg (the default), P50, P90, P95, P99, Max or Min. View Chart in the dashboard opens the chart this way. Sum is deliberately absent: adding durations up across alerts produces a number with no meaning. Allowed values: `Avg`, `Max`, `Min`, `P50`, `P90`, `P95`, `P99`. Defaults to `Avg`.
- `description` (String) Description of what this measurement means to your team.
- `end_alert_state_id` (String) ID of the alert state this measurement ends at. Required only when the End Anchor Type is State Entered. Cleared if that state is deleted, at which point the measurement reports Not Applicable rather than a wrong number. The ID of a `oneuptime_alert_state`.
- `end_alert_state_role` (String) The role of the state this measurement ends at - Created, Acknowledged or Resolved. Used when the End Anchor Type is State Role Entered. Allowed values: `Created`, `Acknowledged`, `Resolved`.
- `end_state_occurrence` (String) Which entry to use when the end state is entered more than once - First or Last. First matches the built-in alert metrics; Last follows a reopened alert to its final pass through that state. Allowed values: `First`, `Last`. Defaults to `First`.
- `is_enabled` (Boolean) Whether this measurement is computed for new and updated alerts. Defaults to `true`.
- `key` (String) Stable, machine readable identifier for this measurement, unique within the project: lowercase letters, numbers and hyphens. Leave it out and it is made from the name - Time to Acknowledge becomes time-to-acknowledge, with -2, -3 and so on added when another measurement already has it. It cannot be changed once the measurement is created, because it is used to build the metric name that every recorded point is written under; to rename a measurement, change the Name instead.
- `order` (Number) Where this measurement appears in the list of measurements, lowest number first. A new measurement is added to the end of the list. Setting a number places it there and moves no other one: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).
- `show_on_alert_view` (Boolean) Whether this measurement is shown on the alert page alongside the alert's other timings. Defaults to `true`.
- `start_alert_state_id` (String) ID of the alert state this measurement starts at. Required only when the Start Anchor Type is State Entered. Cleared if that state is deleted, at which point the measurement reports Not Applicable rather than a wrong number. The ID of a `oneuptime_alert_state`.
- `start_alert_state_role` (String) The role of the state this measurement starts at - Created, Acknowledged or Resolved. Used when the Start Anchor Type is State Role Entered. Resolving by role keeps the measurement working when a project renames or replaces the state that plays that part. Allowed values: `Created`, `Acknowledged`, `Resolved`.
- `start_state_occurrence` (String) Which entry to use when the start state is entered more than once - First or Last. First matches the built-in alert metrics; Last follows a reopened alert to its final pass through that state. Allowed values: `First`, `Last`. Defaults to `First`.
- `unit` (String) The unit this measurement's charts are in: seconds (the default), minutes, hours or days. Every value is worked out in seconds and stored that way on the alert; each chart point is written in this unit, so a chart in hours reads 1.5 for an hour and a half. With seconds, charts show seconds, minutes, hours or days as the numbers grow. Changing it rewrites the measurement's chart points in the new unit. A value that is not a time unit charts in seconds. Defaults to `seconds`.

### Read-Only

- `backfill_completed_at` (String) When the backfill of this measurement over existing alerts finished.
- `backfill_cursor_created_at` (String) How far the backfill has walked this project, so a restart resumes instead of starting over.
- `backfill_cursor_id` (String) The ID of the last alert the backfill worked out, which orders the alerts created at the same moment.
- `backfill_requested_at` (String) When a backfill of this measurement over existing alerts was requested.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_system_defined` (Boolean) Whether this measurement was seeded by OneUptime rather than created by your team.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `metric_name` (String) The metric name every recorded point of this measurement is written under. Derived from the key as oneuptime.alert.measurement.<key> and maintained for you.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert measurement by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_measurement.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_measurement.example <id>
```
