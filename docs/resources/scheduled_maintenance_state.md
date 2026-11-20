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

<a id="order"></a>
## Order

`order` is this scheduled maintenance state's place among the scheduled maintenance states of its project, lowest first. Setting it places this scheduled maintenance state and moves no other, so a configuration that numbers every scheduled maintenance state of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.scheduled_maintenance_states`:

```terraform
order = index(local.scheduled_maintenance_states, each.key) + 1
```

- Give each scheduled maintenance state a number of its own: two scheduled maintenance states at one number are shown the older first.
- A scheduled maintenance state created without `order` goes just above the completed state: it takes the completed state's number, and the completed state and the scheduled maintenance states right below it step down one. Give it `order` when the configuration numbers the others, or the next plan shows theirs changed.
- The scheduled state, ongoing state, ended state and completed state stay in that order, from the top: a write that would put two of them at one number, or out of order, is refused. An apply that moves them past each other's numbers can meet that refusal on one write; apply again and it finishes.
- Deleting a scheduled maintenance state leaves a gap in the numbers, which changes nothing: scheduled maintenance states go by their numbers, not by counting.
- Moving a scheduled maintenance state in the OneUptime Dashboard numbers the project's scheduled maintenance states 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_archived` (Boolean) Archived states stay on the maintenance events that are in them, but can't be picked again. The scheduled, ongoing, ended and completed states can't be archived. Defaults to `false`.
- `is_ended_state` (Boolean) Is this state a ended state? Defaults to `false`.
- `is_ongoing_state` (Boolean) Is this state a ongoing state? Defaults to `false`.
- `is_resolved_state` (Boolean) Is this state a resolved state? Defaults to `false`.
- `is_scheduled_state` (Boolean) Is this state a scheduled state? Defaults to `false`.
- `order` (Number) Where this state sits in the project's list of scheduled maintenance states: 1 is the top. Events only ever move down the list, so the scheduled, ongoing, ended and completed states have to stay in that order. A new state without a number goes just above the completed state. Setting a number places the state there and moves no other: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `archived_at` (String) When this scheduled maintenance state was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
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
