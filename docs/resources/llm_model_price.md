---
page_title: "oneuptime_llm_model_price Resource - oneuptime"
subcategory: "AI"
description: |-
  Custom per-project LLM pricing. When a span carries token counts but no reported cost, ingest prices it against these entries and the built-in list-price catalog — the longest matching model prefix wins, and a project entry beats a built-in one on ties.
---

# oneuptime_llm_model_price (Resource)

Custom per-project LLM pricing. When a span carries token counts but no reported cost, ingest prices it against these entries and the built-in list-price catalog — the longest matching model prefix wins, and a project entry beats a built-in one on ties.

## Example Usage

```terraform
resource "oneuptime_llm_model_price" "example" {
  model_prefix                           = "Example short text"
  input_price_per_million_tokens_in_usd  = 42
  output_price_per_million_tokens_in_usd = 42
  description                            = "Managed by Terraform"
}
```

## Schema

### Required

- `input_price_per_million_tokens_in_usd` (Number) Price of one million input (prompt) tokens in USD. Use 0 for free input tokens.
- `model_prefix` (String) Model-name prefix this price matches, e.g. gpt-4o or my-custom-finetune. Stored lowercase; the longest matching prefix wins and a project entry beats a built-in catalog entry on ties.
- `output_price_per_million_tokens_in_usd` (Number) Price of one million output (completion) tokens in USD. Use 0 for free output tokens (e.g. embeddings).

### Optional

- `description` (String) Description of this price entry, e.g. which negotiated rate or self-hosted deployment it reflects.
- `is_enabled` (Boolean) Whether this price entry is used when pricing LLM spans. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this price entry. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this LLM model price belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing llm model price by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_llm_model_price.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_llm_model_price.example <id>
```
