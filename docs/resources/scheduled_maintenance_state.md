---
page_title: "oneuptime_scheduled_maintenance_state Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage different scheduled maintenance state to your project (Scheduled, Ongoing, Completed for example)
---

# oneuptime_scheduled_maintenance_state (Resource)

Manage different scheduled maintenance state to your project (Scheduled, Ongoing, Completed for example)

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_state" "example" {
  name = "Example short text"
  color = jsonencode({
    "_type": "Color",
    "value": "#ff0000"
  })
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object..
- `color` (String) Color object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Friendly description that will help you remember..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_scheduled_state` (Bool) Is this state a scheduled state?..
- `is_ongoing_state` (Bool) Is this state a ongoing state?..
- `is_ended_state` (Bool) Is this state a ended state?..
- `is_resolved_state` (Bool) Is this state a resolved state?..
- `order` (Number) Where this state sits in the project's list of scheduled maintenance states: 1 is the top. Events only ever move down the list, so the scheduled, ongoing, ended and completed states have to stay in that order. A new state without a number goes just above the completed state. Setting a number moves the state to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_scheduled_maintenance_state.example <id>
```
