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
- `is_created_state` (Bool) Is it the created state of the incident?..
- `is_acknowledged_state` (Bool) Is it the acknowledged state of the incident?..
- `is_resolved_state` (Bool) Is it the resolved state of the incident?..
- `order` (Number) Where this state sits in the project's list of incident states: 1 is the top. Incidents only ever move down the list, and an incident in a state at or below the acknowledged (or resolved) state counts as acknowledged (or resolved), so the created, acknowledged and resolved states have to stay in that order. A new state without a number goes just above the resolved state. Setting a number moves the state to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_state.example <id>
```
