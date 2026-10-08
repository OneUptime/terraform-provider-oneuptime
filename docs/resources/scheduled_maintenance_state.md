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
  name        = "Example scheduled maintenance state"
  color       = "#ff0000"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_ended_state` (Boolean) Is this state a ended state? Defaults to `false`.
- `is_ongoing_state` (Boolean) Is this state a ongoing state? Defaults to `false`.
- `is_resolved_state` (Boolean) Is this state a resolved state? Defaults to `false`.
- `is_scheduled_state` (Boolean) Is this state a scheduled state? Defaults to `false`.
- `order` (Number) Where this state sits in the project's list of scheduled maintenance states: 1 is the top. Events only ever move down the list, so the scheduled, ongoing, ended and completed states have to stay in that order. A new state without a number goes just above the completed state. Setting a number moves the state to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing scheduled maintenance state by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_scheduled_maintenance_state.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_scheduled_maintenance_state.example <id>
```
