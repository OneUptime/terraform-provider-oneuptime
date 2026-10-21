---
page_title: "oneuptime_ai_connection Data Source - oneuptime"
subcategory: "Other"
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

- `conversation_readers` (String) For a company AI tool: 'admins' lets only project owners and admins read its conversations; 'everyone' lets everyone who reads the project's AI conversations read them.
- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this connection is read on its schedule.
- `last_error` (String) Why the last read failed, with credentials removed. Cleared by the next read that works.
- `name` (String) Friendly name for this connection.
- `provider_value` (String) What this connection reads, such as 'openai', 'anthropic' or 'claude-enterprise'. Fixed once created.
- `service_id` (String) ID of the telemetry service everything this connection brings in belongs to. The ID of a `oneuptime_service`.
- `store_conversation_text` (Boolean) For a company AI tool: keep what people asked and what the AI answered. When off, only who used it, when and how much is kept.

### Read-Only

- `config` (String) The provider's settings that are not secret, such as a tenant, a subscription or a Region. Keys are defined by the provider. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `last_pull_result` (String) What the last read found: how many days, records and conversations, and how long it took. A JSON value: write it with `jsonencode()`.
- `last_pulled_at` (String) When the connection was last read. Empty means never.
- `last_successful_pull_at` (String) When a read last finished without an error.
- `next_pull_at` (String) When the connection is read next. Empty means as soon as possible.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
