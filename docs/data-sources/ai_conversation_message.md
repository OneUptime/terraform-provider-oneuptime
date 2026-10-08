---
page_title: "oneuptime_ai_conversation_message Data Source - oneuptime"
subcategory: "Other"
description: |-
  A message in an AI conversation. Assistant messages carry citations, tool events and cost.
---

# oneuptime_ai_conversation_message (Data Source)

A message in an AI conversation. Assistant messages carry citations, tool events and cost.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai conversation message may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_ai_conversation_message" "example" {
  conversation_id = oneuptime_ai_conversation.example.id
}

# Or by id:
data "oneuptime_ai_conversation_message" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_run_id` (String) ID of the AI run that produced this assistant message.
- `content_in_markdown` (String) Message content in markdown.
- `conversation_id` (String) ID of the conversation this message belongs to. The ID of a `oneuptime_ai_conversation`.
- `error_message` (String) Error message if this message failed to generate.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `role` (String) Who authored this message: User or Assistant.
- `status` (String) Current status of this message.
- `user_feedback` (String) Thumbs feedback the user left on this assistant message: Up or Down.
- `user_id` (String) ID of the user who owns the conversation. The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `citations` (String) Server-minted citations for this assistant message. Each citation records the tool, the exact validated query arguments and the row count. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this message belongs to. The ID of a `oneuptime_project`.
- `tool_actions` (String) Mutating actions the agent proposed or performed in this turn, with their approval status (pending, approved, denied, executed). A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
- `widgets` (String) Inline widgets (charts, tables, trace waterfalls, resource cards) built from this assistant message's tool results and rendered inline in the chat. A JSON value: write it with `jsonencode()`.
