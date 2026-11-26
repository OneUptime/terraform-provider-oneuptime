---
page_title: "oneuptime_docker_host_team_owner Data Source - oneuptime"
subcategory: "Infrastructure"
description: |-
  Add teams as owners to your Docker hosts.
---

# oneuptime_docker_host_team_owner (Data Source)

Add teams as owners to your Docker hosts.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one docker host team owner may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_docker_host_team_owner" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_docker_host_team_owner" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `docker_host_id` (String) ID of the Docker host the team owns. The ID of a `oneuptime_docker_host`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `team_id` (String) ID of the team that owns the Docker host. The ID of a `oneuptime_team`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
