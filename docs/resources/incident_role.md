---
page_title: "oneuptime_incident_role Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident roles for your project (Incident Commander, Responder, etc.). Add, edit, or remove roles.
---

# oneuptime_incident_role (Resource)

Manage incident roles for your project (Incident Commander, Responder, etc.). Add, edit, or remove roles.

## Example Usage

```terraform
resource "oneuptime_incident_role" "example" {
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
- `role_icon` (String) Icon for this incident role (e.g., User, Shield, etc.)..
- `can_assign_multiple_users` (Bool) Can multiple users be assigned to this role? If false, only one user can be assigned...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_primary_role` (Bool) Is this the primary incident role? Primary roles like Incident Commander have special significance...
- `is_deleteable` (Bool) Can this role be deleted? Primary roles cannot be deleted...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_role.example <id>
```
