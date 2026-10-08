---
page_title: "oneuptime_incident_episode_user_owner Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Add users as owners to your incident episodes.
---

# oneuptime_incident_episode_user_owner (Data Source)

Add users as owners to your incident episodes.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident episode user owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_episode_user_owner" "example" {
  user_id = data.oneuptime_user.example.id
}

# Or by id:
data "oneuptime_incident_episode_user_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_episode_id` (String) ID of your OneUptime Incident Episode in which this object belongs. The ID of a `oneuptime_incident_episode`.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `user_id` (String) ID of the user who is the owner. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
