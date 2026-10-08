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

### Read-Only

- `alert_id` (String) ID of the alert whose investigation box this shared conversation belongs to. The ID of a `oneuptime_alert`.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this conversation.
- `id` (String) Unique identifier for the resource.
- `incident_id` (String) ID of the incident whose investigation box this shared conversation belongs to. The ID of a `oneuptime_incident`.
- `last_message_at` (String) When the last message in this conversation was sent.
- `llm_provider_id` (String) The LLM provider selected for this conversation. If empty, the project default (or global) provider is used.
- `page_context` (String) The dashboard page (entity) this conversation is about. Set from the first message that carried a page context. A JSON value: write it with `jsonencode()`.
- `permission_mode` (String) How the agent is allowed to run mutating tools: AskForApproval, AutoRun or ReadOnly.
- `project_id` (String) ID of the project this conversation belongs to. The ID of a `oneuptime_project`.
- `title` (String) Title of the conversation. Generated from the first message.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing ai conversation by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_ai_conversation.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_ai_conversation.example <id>
```
