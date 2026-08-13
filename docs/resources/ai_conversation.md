---
page_title: "oneuptime_ai_conversation Resource - oneuptime"
subcategory: "Other"
description: |-
  A conversation with the OneUptime AI about observability data (logs, traces, metrics, exceptions, incidents, monitors and alerts).
---

# oneuptime_ai_conversation (Resource)

A conversation with the OneUptime AI about observability data (logs, traces, metrics, exceptions, incidents, monitors and alerts).

## Example Usage

```terraform
resource "oneuptime_ai_conversation" "example" {

}
```

## Schema

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `title` (String) Title of the conversation. Generated from the first message...
- `last_message_at` (String) A date time object..
- `llm_provider_id` (String) A unique identifier for an object, represented as a UUID..
- `permission_mode` (String) How the agent is allowed to run mutating tools: AskForApproval, AutoRun or ReadOnly...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_ai_conversation.example <id>
```
