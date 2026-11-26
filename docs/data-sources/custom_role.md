---
page_title: "oneuptime_custom_role Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  A set of permissions you name once, such as Incident Responder, and give to teams and API keys. Changing a role changes every team and key that has it.
---

# oneuptime_custom_role (Data Source)

A set of permissions you name once, such as Incident Responder, and give to teams and API keys. Changing a role changes every team and key that has it.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one custom role may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_custom_role" "example" {
  name = "Example custom role"
}

# Or by id:
data "oneuptime_custom_role" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) What the role is for, in a sentence or two.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) What the role is called, such as Incident Responder. Unique in the project, and never the name of a built-in role.
- `permission_count` (Number) How many permissions the role holds. OneUptime keeps it as permissions are added to the role and taken out of it.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
