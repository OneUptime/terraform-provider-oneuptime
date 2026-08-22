---
page_title: "oneuptime_llm_model_price Resource - oneuptime"
subcategory: "Other"
description: |-
  Custom per-project LLM pricing. When a span carries token counts but no reported cost, ingest prices it against these entries and the built-in list-price catalog — the longest matching model prefix wins, and a project entry beats a built-in one on ties.
---

# oneuptime_llm_model_price (Resource)

Custom per-project LLM pricing. When a span carries token counts but no reported cost, ingest prices it against these entries and the built-in list-price catalog — the longest matching model prefix wins, and a project entry beats a built-in one on ties.

## Example Usage

```terraform
resource "oneuptime_llm_model_price" "example" {
  model_prefix = "Example short text"
  input_price_per_million_tokens_in_usd = 42
  output_price_per_million_tokens_in_usd = 42
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `model_prefix` (String) Model-name prefix this price matches, e.g. gpt-4o or my-custom-finetune. Stored lowercase; the longest matching prefix wins and a project entry beats a built-in catalog entry on ties...
- `input_price_per_million_tokens_in_usd` (Number) Price of one million input (prompt) tokens in USD. Use 0 for free input tokens...
- `output_price_per_million_tokens_in_usd` (Number) Price of one million output (completion) tokens in USD. Use 0 for free output tokens (e.g. embeddings)...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this price entry, e.g. which negotiated rate or self-hosted deployment it reflects...
- `is_enabled` (Bool) Whether this price entry is used when pricing LLM spans...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_llm_model_price.example <id>
```
