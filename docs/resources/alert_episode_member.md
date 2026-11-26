---
page_title: "oneuptime_alert_episode_member Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Link between alerts and episodes
---

# oneuptime_alert_episode_member (Resource)

Link between alerts and episodes

## Example Usage

```terraform
resource "oneuptime_alert_episode_member" "example" {
  alert_episode_id = oneuptime_alert_episode.example.id
  alert_id         = oneuptime_alert.example.id
}
```

## Schema

### Required

- `alert_episode_id` (String) ID of the Alert Episode that this alert belongs to. The ID of a `oneuptime_alert_episode`.
- `alert_id` (String) ID of the Alert that is a member of this episode. The ID of a `oneuptime_alert`.

### Optional

- `added_at` (String) When this alert was added to the episode.
- `added_by` (String) How this alert was added to the episode (rule, manual, or api). Allowed values: `rule`, `manual`, `api`. Defaults to `rule`.

### Read-Only

- `added_by_user_id` (String) User ID who manually added this alert to the episode. The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified_of_alert_added` (Boolean) Has the owner been notified that this alert was added to the episode?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `matched_rule_id` (String) ID of the grouping rule that matched this alert. The ID of a `oneuptime_alert_grouping_rule`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert episode member by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_episode_member.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_episode_member.example <id>
```
