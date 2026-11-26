---
page_title: "oneuptime_ai_conversation Data Source - oneuptime"
subcategory: "AI"
description: |-
  A conversation with the OneUptime AI about observability data (logs, traces, metrics, exceptions, incidents, monitors and alerts).
---

# oneuptime_ai_conversation (Data Source)

A conversation with the OneUptime AI about observability data (logs, traces, metrics, exceptions, incidents, monitors and alerts).

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai conversation may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_ai_conversation" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_ai_conversation" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of the alert whose investigation box this shared conversation belongs to. The ID of a `oneuptime_alert`.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this conversation.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident whose investigation box this shared conversation belongs to. The ID of a `oneuptime_incident`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `llm_provider_id` (String) The LLM provider selected for this conversation. If empty, the project default (or global) provider is used.
- `permission_mode` (String) How the agent is allowed to run mutating tools: AskForApproval, AutoRun or ReadOnly.
- `title` (String) Title of the conversation. Generated from the first message.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_message_at` (String) When the last message in this conversation was sent.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `page_context` (String) The dashboard page (entity) this conversation is about. Set from the first message that carried a page context. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of the project this conversation belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
