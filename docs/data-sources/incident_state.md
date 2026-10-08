---
page_title: "oneuptime_incident_state Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident states for your project (Created, Acknowledged for example). Add / edit or remove states.
---

# oneuptime_incident_state (Data Source)

Manage incident states for your project (Created, Acknowledged for example). Add / edit or remove states.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident state may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_state" "example" {
  name = "Example incident state"
}

# Or by id:
data "oneuptime_incident_state" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_acknowledged_state` (Boolean) Is it the acknowledged state of the incident?
- `is_created_state` (Boolean) Is it the created state of the incident?
- `is_resolved_state` (Boolean) Is it the resolved state of the incident?
- `name` (String) Any friendly name of this object.
- `order` (Number) Where this state sits in the project's list of incident states: 1 is the top. Incidents only ever move down the list, and an incident in a state at or below the acknowledged (or resolved) state counts as acknowledged (or resolved), so the created, acknowledged and resolved states have to stay in that order. A new state without a number goes just above the resolved state. Setting a number moves the state to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
