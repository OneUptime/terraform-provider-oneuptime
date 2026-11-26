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
- `added_by` (String) How this incident was added to the episode (rule, manual, or api). Allowed values: `rule`, `manual`, `api`. Defaults to `rule`.

### Read-Only

- `added_by_user_id` (String) User ID who manually added this incident to the episode. The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified_of_incident_added` (Boolean) Has the owner been notified that this incident was added to the episode?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `matched_rule_id` (String) ID of the grouping rule that matched this incident. The ID of a `oneuptime_incident_grouping_rule`.
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
