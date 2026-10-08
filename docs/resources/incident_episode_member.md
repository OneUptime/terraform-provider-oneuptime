---
page_title: "oneuptime_incident_episode_member Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Link between incidents and episodes
---

# oneuptime_incident_episode_member (Resource)

Link between incidents and episodes

## Example Usage

```terraform
resource "oneuptime_incident_episode_member" "example" {
  incident_episode_id = oneuptime_incident_episode.example.id
  incident_id         = oneuptime_incident.example.id
}
```

## Schema

### Required

- `incident_episode_id` (String) ID of the Incident Episode that this incident belongs to. The ID of a `oneuptime_incident_episode`.
- `incident_id` (String) ID of the Incident that is a member of this episode. The ID of a `oneuptime_incident`.

### Optional

- `added_at` (String) When this incident was added to the episode.
- `added_by` (String) How this incident was added to the episode (rule, manual, or api). Defaults to `rule`.
- `is_owner_notified_of_incident_added` (Boolean) Has the owner been notified that this incident was added to the episode?
- `matched_rule_id` (String) ID of the grouping rule that matched this incident. The ID of a `oneuptime_incident_grouping_rule`.

### Read-Only

- `added_by_user_id` (String) User ID who manually added this incident to the episode. The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident episode member by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_episode_member.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_episode_member.example <id>
```
