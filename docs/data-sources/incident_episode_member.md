---
page_title: "oneuptime_incident_episode_member Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Link between incidents and episodes
---

# oneuptime_incident_episode_member (Data Source)

Link between incidents and episodes

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident episode member may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_episode_member" "example" {
  incident_episode_id = oneuptime_incident_episode.example.id
}

# Or by id:
data "oneuptime_incident_episode_member" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `added_by` (String) How this incident was added to the episode (rule, manual, or api).
- `added_by_user_id` (String) User ID who manually added this incident to the episode. The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_episode_id` (String) ID of the Incident Episode that this incident belongs to. The ID of a `oneuptime_incident_episode`.
- `incident_id` (String) ID of the Incident that is a member of this episode. The ID of a `oneuptime_incident`.
- `is_owner_notified_of_incident_added` (Boolean) Has the owner been notified that this incident was added to the episode?
- `matched_rule_id` (String) ID of the grouping rule that matched this incident. The ID of a `oneuptime_incident_grouping_rule`.

### Read-Only

- `added_at` (String) When this incident was added to the episode.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
