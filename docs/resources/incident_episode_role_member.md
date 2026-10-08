---
page_title: "oneuptime_incident_episode_role_member Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Assign users with specific roles to incident episodes. These assignments propagate to all incidents in the episode.
---

# oneuptime_incident_episode_role_member (Resource)

Assign users with specific roles to incident episodes. These assignments propagate to all incidents in the episode.

## Example Usage

```terraform
resource "oneuptime_incident_episode_role_member" "example" {
  user_id             = data.oneuptime_user.example.id
  incident_episode_id = oneuptime_incident_episode.example.id
  incident_role_id    = oneuptime_incident_role.example.id
}
```

## Schema

### Required

- `incident_episode_id` (String) ID of your OneUptime Incident Episode in which this object belongs. The ID of a `oneuptime_incident_episode`.
- `incident_role_id` (String) ID of the Incident Role assigned to this user. The ID of a `oneuptime_incident_role`.
- `user_id` (String) ID of your OneUptime User assigned to this episode. The ID of a `oneuptime_user` (see the data source).

### Optional

- `notes` (String) Assignment context or notes.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident episode role member by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_episode_role_member.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_episode_role_member.example <id>
```
