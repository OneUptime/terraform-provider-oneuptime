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
  name        = "Example incident role"
  color       = "#ff0000"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `name` (String) Any friendly name of this object.

### Optional

- `can_assign_multiple_users` (Boolean) Can multiple users be assigned to this role? If false, only one user can be assigned.
- `description` (String) Friendly description that will help you remember.
- `role_icon` (String) Icon for this incident role (e.g., User, Shield, etc.).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_deleteable` (Boolean) Can this role be deleted? Primary roles cannot be deleted.
- `is_primary_role` (Boolean) Is this the primary incident role? Primary roles like Incident Commander have special significance.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident role by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_role.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_role.example <id>
```
