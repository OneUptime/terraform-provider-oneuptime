---
page_title: "oneuptime_incident_episode_user_owner Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Add users as owners to your incident episodes.
---

# oneuptime_incident_episode_user_owner (Resource)

Add users as owners to your incident episodes.

## Example Usage

```terraform
resource "oneuptime_incident_episode_user_owner" "example" {
  user_id             = data.oneuptime_user.example.id
  incident_episode_id = oneuptime_incident_episode.example.id
}
```

## Schema

### Required

- `incident_episode_id` (String) ID of your OneUptime Incident Episode in which this object belongs. The ID of a `oneuptime_incident_episode`.
- `user_id` (String) ID of the user who is the owner. The ID of a `oneuptime_user` (see the data source).

### Optional

- `is_owner_notified` (Boolean) Are owners notified of this resource ownership? Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident episode user owner by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_episode_user_owner.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_episode_user_owner.example <id>
```
