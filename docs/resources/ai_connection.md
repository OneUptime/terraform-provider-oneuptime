---
page_title: "oneuptime_ai_connection Resource - oneuptime"
subcategory: "Other"
description: |-
  Connections of AI / LLM observability to AI provider accounts (usage and cost) and to company AI tools (conversations), read on a schedule without any code.
---

# oneuptime_ai_connection (Resource)

Connections of AI / LLM observability to AI provider accounts (usage and cost) and to company AI tools (conversations), read on a schedule without any code.

## Example Usage

```terraform
resource "oneuptime_ai_connection" "example" {
  name           = "Example ai connection"
  provider_value = "Example short text"
  secrets        = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
}
```

## Schema

### Required

- `name` (String) Friendly name for this connection.
- `provider_value` (String) What this connection reads, such as 'openai', 'anthropic' or 'claude-enterprise'. Fixed once created.
- `secrets` (String, Sensitive) The provider's secrets (an admin key, a client secret, a service account key) as a JSON object. Encrypted at rest and never returned by the API.

### Optional

- `config` (String) The provider's settings that are not secret, such as a tenant, a subscription or a Region. Keys are defined by the provider. A JSON value: write it with `jsonencode()`.
- `conversation_readers` (String) For a company AI tool: 'admins' lets only project owners and admins read its conversations; 'everyone' lets everyone who reads the project's AI conversations read them. Defaults to `admins`.
- `is_enabled` (Boolean) Whether this connection is read on its schedule. Defaults to `true`.
- `store_conversation_text` (Boolean) For a company AI tool: keep what people asked and what the AI answered. When off, only who used it, when and how much is kept. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this connection. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_error` (String) Why the last read failed, with credentials removed. Cleared by the next read that works.
- `last_pull_result` (String) What the last read found: how many days, records and conversations, and how long it took. A JSON value: write it with `jsonencode()`.
- `last_pulled_at` (String) When the connection was last read. Empty means never.
- `last_successful_pull_at` (String) When a read last finished without an error.
- `next_pull_at` (String) When the connection is read next. Empty means as soon as possible.
- `project_id` (String) ID of the project this connection belongs to. The ID of a `oneuptime_project`.
- `service_id` (String) ID of the telemetry service everything this connection brings in belongs to. The ID of a `oneuptime_service`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing ai connection by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_ai_connection.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_ai_connection.example <id>
```
