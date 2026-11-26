---
page_title: "oneuptime_custom_role_permission Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  The permissions in a custom role, each with the scope and labels a team's own permission takes.
---

# oneuptime_custom_role_permission (Data Source)

The permissions in a custom role, each with the scope and labels a team's own permission takes.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one custom role permission may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_custom_role_permission" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_custom_role_permission" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `custom_role_id` (String) ID of the custom role this permission is part of. The ID of a `oneuptime_custom_role`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `permission` (String) A built-in role, such as IncidentMember, or a single permission, such as ReadProjectIncident. You can find the list of permissions in the Permission Reference. ProjectOwner cannot be part of a custom role.
- `scope` (String) Which resources this permission reaches. One of: All (every resource in the project), Owned (resources the team given the role, or one of its members, owns), Labels (resources carrying one of the labels). Defaults to All.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `labels` (Set of String) With the scope Labels: the labels this permission is limited to. A resource carrying at least one of them is reached. IDs of `oneuptime_label` resources.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
