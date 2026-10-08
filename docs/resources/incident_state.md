---
page_title: "oneuptime_incident_state Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident states for your project (Created, Acknowledged for example). Add / edit or remove states.
---

# oneuptime_incident_state (Resource)

Manage incident states for your project (Created, Acknowledged for example). Add / edit or remove states.

## Example Usage

```terraform
resource "oneuptime_incident_state" "example" {
  name        = "Example incident state"
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
- `is_acknowledged_state` (Boolean) Is it the acknowledged state of the incident?
- `is_created_state` (Boolean) Is it the created state of the incident?
- `is_resolved_state` (Boolean) Is it the resolved state of the incident?
- `order` (Number) Where this state sits in the project's list of incident states: 1 is the top. Incidents only ever move down the list, and an incident in a state at or below the acknowledged (or resolved) state counts as acknowledged (or resolved), so the created, acknowledged and resolved states have to stay in that order. A new state without a number goes just above the resolved state. Setting a number moves the state to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident state by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_state.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_state.example <id>
```
