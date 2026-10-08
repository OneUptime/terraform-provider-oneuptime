---
page_title: "oneuptime_incident_role Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage incident roles for your project (Incident Commander, Responder, etc.). Add, edit, or remove roles.
---

# oneuptime_incident_role (Data Source)

Manage incident roles for your project (Incident Commander, Responder, etc.). Add, edit, or remove roles.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident role may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_role" "example" {
  name = "Example incident role"
}

# Or by id:
data "oneuptime_incident_role" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `can_assign_multiple_users` (Boolean) Can multiple users be assigned to this role? If false, only one user can be assigned.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_deleteable` (Boolean) Can this role be deleted? Primary roles cannot be deleted.
- `is_primary_role` (Boolean) Is this the primary incident role? Primary roles like Incident Commander have special significance.
- `name` (String) Any friendly name of this object.
- `role_icon` (String) Icon for this incident role (e.g., User, Shield, etc.).
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
