---
page_title: "oneuptime_incident_member Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Assign users with specific roles to incidents.
---

# oneuptime_incident_member (Resource)

Assign users with specific roles to incidents.

## Example Usage

```terraform
resource "oneuptime_incident_member" "example" {
  incident_id      = oneuptime_incident.example.id
  incident_role_id = oneuptime_incident_role.example.id
  user_email       = "alice@example.com"
}
```

## Schema

### Required

- `incident_id` (String) ID of your OneUptime Incident in which this object belongs. The ID of a `oneuptime_incident`.
- `incident_role_id` (String) ID of the Incident Role assigned to this user. The ID of a `oneuptime_incident_role`.

### Optional

- `notes` (String) Assignment context or notes.
- `user_email` (String) The email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) ID of your OneUptime User in which this object belongs. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident member by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_member.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_member.example <id>
```
