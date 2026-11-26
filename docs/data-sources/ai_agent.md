---
page_title: "oneuptime_ai_agent Data Source - oneuptime"
subcategory: "AI"
description: |-
  Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management.
---

# oneuptime_ai_agent (Data Source)

Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai agent may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_ai_agent" "example" {
  name = "Example ai agent"
}

# Or by id:
data "oneuptime_ai_agent" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `connection_status` (String) Connection Status of the AI Agent.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `description` (String) Friendly description that will help you remember.
- `icon_file_id` (String) AI Agent Icon File ID. The ID of a `oneuptime_file`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_default` (Boolean) Is this the default AI Agent for the project? When set, this agent will be used for automated tasks.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `ai_agent_version` (String) A version label for the AI agent. Leave it out and it is 1.0.0.
- `created_at` (String) Date and Time when the object was created.
- `key` (String, Sensitive) The secret the AI agent connects with, its AI_AGENT_KEY. Leave it out and OneUptime generates one.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.
- `last_alive` (String) When the AI agent last checked in. Set by OneUptime.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project the AI agent belongs to. A global AI agent belongs to none.
- `updated_at` (String) Date and Time when the object was updated.
