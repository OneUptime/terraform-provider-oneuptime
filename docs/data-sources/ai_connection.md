---
page_title: "oneuptime_ai_connection Data Source - oneuptime"
subcategory: "AI"
description: |-
  Connections of AI / LLM observability to AI provider accounts (usage and cost) and to company AI tools (conversations), read on a schedule without any code.
---

# oneuptime_ai_connection (Data Source)

Connections of AI / LLM observability to AI provider accounts (usage and cost) and to company AI tools (conversations), read on a schedule without any code.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai connection may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_ai_connection" "example" {
  name = "Example ai connection"
}

# Or by id:
data "oneuptime_ai_connection" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `conversation_readers` (String) For a company AI tool: 'admins' lets only Project Owners and Project Admins read its conversations; 'everyone' lets everyone who reads the project's AI conversations read them.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this connection is read on its schedule.
- `last_error` (String) Why the last read failed, with credentials removed. Cleared by the next read that works.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Friendly name for this connection.
- `provider_value` (String) What this connection reads, such as 'openai', 'anthropic' or 'claude-enterprise'. Fixed once created.
- `service_id` (String) ID of the telemetry service everything this connection brings in belongs to. The ID of a `oneuptime_service`.
- `store_conversation_text` (Boolean) For a company AI tool: keep what people asked and what the AI answered. When off, only who used it, when and how much is kept.

### Read-Only

- `config` (String) The provider's settings that are not secret, such as a tenant, a subscription or a Region. Keys are defined by the provider. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_pull_result` (String) What the last read found: how many days, records and conversations, and how long it took. A JSON value: read it with `jsondecode()`.
- `last_pulled_at` (String) When the connection was last read. Empty means never.
- `last_successful_pull_at` (String) When a read last finished without an error.
- `next_pull_at` (String) When the connection is read next. Empty means as soon as possible.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
