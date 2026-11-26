---
page_title: "oneuptime_scheduled_maintenance_measurement_value Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  The computed value of one measurement for one maintenance event. Written by OneUptime, never edited by hand.
---

# oneuptime_scheduled_maintenance_measurement_value (Data Source)

The computed value of one measurement for one maintenance event. Written by OneUptime, never edited by hand.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled maintenance measurement value may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_measurement_value" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_scheduled_maintenance_measurement_value" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `end_scheduled_maintenance_state_timeline_id` (String) The state timeline entry the end anchor resolved to, when it resolved to one. Not a foreign key - the value survives the timeline entry being removed.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `scheduled_maintenance_id` (String) ID of the Maintenance Event this measurement value was computed for. The ID of a `oneuptime_scheduled_maintenance_event`.
- `scheduled_maintenance_measurement_id` (String) ID of the Scheduled Maintenance Measurement definition this value was computed from. The ID of a `oneuptime_scheduled_maintenance_measurement`.
- `start_scheduled_maintenance_state_timeline_id` (String) The state timeline entry the start anchor resolved to, when it resolved to one. Not a foreign key - the value survives the timeline entry being removed.
- `status` (String) Outcome of evaluating this measurement - Recorded, Pending, Not Applicable or Invalid.
- `status_message` (String) Why this measurement has the status it has - which anchor is still open, or why it can never resolve.
- `value_in_seconds` (Number) The measured duration in seconds. Empty unless the status is Recorded.

### Read-Only

- `computed_at` (String) When this measurement value was last computed.
- `created_at` (String) Date and Time when the object was created.
- `ended_at` (String) When the end anchor of this measurement resolved.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `started_at` (String) When the start anchor of this measurement resolved.
- `updated_at` (String) Date and Time when the object was updated.
