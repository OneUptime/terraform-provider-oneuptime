---
page_title: "oneuptime_ai_conversation_message Data Source - oneuptime"
subcategory: "Other"
description: |-
  A message in an AI conversation. Assistant messages carry citations, tool events and cost.
---

# oneuptime_ai_conversation_message (Data Source)

A message in an AI conversation. Assistant messages carry citations, tool events and cost. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_ai_conversation_message" "by_name" {
  name = "example-ai_conversation_message"
}

data "oneuptime_ai_conversation_message" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `conversation_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `role` (String) Who authored this message: User or Assistant... Computed.
- `content_in_markdown` (String) Message content in markdown... Computed.
- `status` (String) Current status of this message... Computed.
- `ai_run_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `citations` (String) Server-minted citations for this assistant message. Each citation records the tool, the exact validated query arguments and the row count... Computed.
- `widgets` (String) Inline widgets (charts, tables, trace waterfalls, resource cards) built from this assistant message's tool results and rendered inline in the chat... Computed.
- `tool_actions` (String) Mutating actions the agent proposed or performed in this turn, with their approval status (pending, approved, denied, executed)... Computed.
- `error_message` (String) Error message if this message failed to generate... Computed.
- `user_feedback` (String) Thumbs feedback the user left on this assistant message: Up or Down... Computed.
