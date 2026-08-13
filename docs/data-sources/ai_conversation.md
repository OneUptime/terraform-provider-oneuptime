---
page_title: "oneuptime_ai_conversation Data Source - oneuptime"
subcategory: "Other"
description: |-
  A conversation with the OneUptime AI about observability data (logs, traces, metrics, exceptions, incidents, monitors and alerts).
---

# oneuptime_ai_conversation (Data Source)

A conversation with the OneUptime AI about observability data (logs, traces, metrics, exceptions, incidents, monitors and alerts). Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_ai_conversation" "by_name" {
  name = "example-ai_conversation"
}

data "oneuptime_ai_conversation" "by_id" {
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
- `title` (String) Title of the conversation. Generated from the first message... Computed.
- `last_message_at` (String) A date time object.. Computed.
- `llm_provider_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `permission_mode` (String) How the agent is allowed to run mutating tools: AskForApproval, AutoRun or ReadOnly... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
