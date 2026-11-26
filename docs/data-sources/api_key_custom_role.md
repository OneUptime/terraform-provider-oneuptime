---
page_title: "oneuptime_api_key_custom_role Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  The custom roles an API key is given. The key holds every permission of the role except those with the Owned scope: a key owns nothing.
---

# oneuptime_api_key_custom_role (Data Source)

The custom roles an API key is given. The key holds every permission of the role except those with the Owned scope: a key owns nothing.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one api key custom role may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_api_key_custom_role" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_api_key_custom_role" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `api_key_id` (String) ID of the API key given the role. The ID of a `oneuptime_api_key`.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `custom_role_id` (String) ID of the custom role given to the API key. The ID of a `oneuptime_custom_role`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
