---
page_title: "oneuptime_queue_team_owner Resource - oneuptime"
subcategory: "Infrastructure"
description: |-
  Add teams as owners to your queues.
---

# oneuptime_queue_team_owner (Resource)

Add teams as owners to your queues.

## Example Usage

```terraform
resource "oneuptime_queue_team_owner" "example" {
  team_id          = oneuptime_team.example.id
  message_queue_id = oneuptime_queue.example.id
}
```

## Schema

### Required

- `message_queue_id` (String) ID of the queue the team owns. The ID of a `oneuptime_queue`.
- `team_id` (String) ID of the team that owns the queue. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing queue team owner by its queue and its `team_id`, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_queue_team_owner.example
  id = "<message_queue_id>/<team_id>"
}
```

or on the command line:

```shell
terraform import oneuptime_queue_team_owner.example <message_queue_id>/<team_id>
```

Its own ID works too.
