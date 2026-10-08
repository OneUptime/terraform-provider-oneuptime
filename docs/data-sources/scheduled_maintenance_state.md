---
page_title: "oneuptime_scheduled_maintenance_state Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage different scheduled maintenance state to your project (Scheduled, Ongoing, Completed for example)
---

# oneuptime_scheduled_maintenance_state (Data Source)

Manage different scheduled maintenance state to your project (Scheduled, Ongoing, Completed for example)

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled maintenance state may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_state" "example" {
  name = "Example scheduled maintenance state"
}

# Or by id:
data "oneuptime_scheduled_maintenance_state" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_ended_state` (Boolean) Is this state a ended state?
- `is_ongoing_state` (Boolean) Is this state a ongoing state?
- `is_resolved_state` (Boolean) Is this state a resolved state?
- `is_scheduled_state` (Boolean) Is this state a scheduled state?
- `name` (String) Any friendly name of this object.
- `order` (Number) Where this state sits in the project's list of scheduled maintenance states: 1 is the top. Events only ever move down the list, so the scheduled, ongoing, ended and completed states have to stay in that order. A new state without a number goes just above the completed state. Setting a number moves the state to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
