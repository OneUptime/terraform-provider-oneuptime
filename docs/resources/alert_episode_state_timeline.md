---
page_title: "oneuptime_alert_episode_state_timeline Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Change state of the alert episodes (Created to Acknowledged for example)
---

# oneuptime_alert_episode_state_timeline (Resource)

Change state of the alert episodes (Created to Acknowledged for example)

## Example Usage

```terraform
resource "oneuptime_alert_episode_state_timeline" "example" {
  alert_episode_id = oneuptime_alert_episode.example.id
  alert_state_id   = oneuptime_alert_state.example.id
}
```

## Schema

### Required

- `alert_episode_id` (String) Relation to Alert Episode ID in which this resource belongs. The ID of a `oneuptime_alert_episode`.
- `alert_state_id` (String) Alert State ID Relation. Which alert state does this episode change to? The ID of a `oneuptime_alert_state`.

### Optional

- `ends_at` (String) When did this status change end? OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `root_cause` (String) What is the root cause of this status change?
- `starts_at` (String) When did this status change?

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of state change?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `state_change_log` (String) A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert episode state timeline by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_episode_state_timeline.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_episode_state_timeline.example <id>
```
