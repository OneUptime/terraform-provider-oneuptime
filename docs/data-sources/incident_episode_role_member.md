---
page_title: "oneuptime_incident_episode_role_member Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Assign users with specific roles to incident episodes. These assignments propagate to all incidents in the episode.
---

# oneuptime_incident_episode_role_member (Data Source)

Assign users with specific roles to incident episodes. These assignments propagate to all incidents in the episode.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident episode role member may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_episode_role_member" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_incident_episode_role_member" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_episode_id` (String) ID of your OneUptime Incident Episode in which this object belongs. The ID of a `oneuptime_incident_episode`.
- `incident_role_id` (String) ID of the Incident Role assigned to this user. The ID of a `oneuptime_incident_role`.
- `notes` (String) Assignment context or notes.
- `user_id` (String) ID of your OneUptime User assigned to this episode. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
