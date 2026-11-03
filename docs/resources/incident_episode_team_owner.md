---
page_title: "oneuptime_incident_episode_team_owner Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Add teams as owners to your incident episodes.
---

# oneuptime_incident_episode_team_owner (Resource)

Add teams as owners to your incident episodes.

## Example Usage

```terraform
resource "oneuptime_incident_episode_team_owner" "example" {
  team_id             = oneuptime_team.example.id
  incident_episode_id = oneuptime_incident_episode.example.id
}
```

## Schema

### Required

- `incident_episode_id` (String) ID of the incident episode the team owns. The ID of a `oneuptime_incident_episode`.
- `team_id` (String) ID of the team that owns the incident episode. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident episode team owner by its incident episode and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_episode_team_owner.example
  id = "<incident_episode_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_episode_team_owner.example <incident_episode_id>/<team_id>
```

Its own ID works too.
