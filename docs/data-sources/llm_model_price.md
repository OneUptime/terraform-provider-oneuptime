---
page_title: "oneuptime_llm_model_price Data Source - oneuptime"
subcategory: "Other"
description: |-
  Custom per-project LLM pricing. When a span carries token counts but no reported cost, ingest prices it against these entries and the built-in list-price catalog — the longest matching model prefix wins, and a project entry beats a built-in one on ties.
---

# oneuptime_llm_model_price (Data Source)

Custom per-project LLM pricing. When a span carries token counts but no reported cost, ingest prices it against these entries and the built-in list-price catalog — the longest matching model prefix wins, and a project entry beats a built-in one on ties. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_llm_model_price" "by_name" {
  name = "example-llm_model_price"
}

data "oneuptime_llm_model_price" "by_id" {
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
- `model_prefix` (String) Model-name prefix this price matches, e.g. gpt-4o or my-custom-finetune. Stored lowercase; the longest matching prefix wins and a project entry beats a built-in catalog entry on ties... Computed.
- `description` (String) Description of this price entry, e.g. which negotiated rate or self-hosted deployment it reflects... Computed.
- `is_enabled` (Bool) Whether this price entry is used when pricing LLM spans... Computed.
- `input_price_per_million_tokens_in_usd` (Number) Price of one million input (prompt) tokens in USD. Use 0 for free input tokens... Computed.
- `output_price_per_million_tokens_in_usd` (Number) Price of one million output (completion) tokens in USD. Use 0 for free output tokens (e.g. embeddings)... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
