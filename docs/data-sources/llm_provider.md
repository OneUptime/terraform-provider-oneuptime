---
page_title: "oneuptime_llm_provider Data Source - oneuptime"
subcategory: "Other"
description: |-
  Manage LLM Provider configurations. Connect to OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAI-compatible servers (e.g. vLLM, LocalAI), or other LLM providers to enable AI features.
---

# oneuptime_llm_provider (Data Source)

Manage LLM Provider configurations. Connect to OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAI-compatible servers (e.g. vLLM, LocalAI), or other LLM providers to enable AI features. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_llm_provider" "by_name" {
  name = "example-llm_provider"
}

data "oneuptime_llm_provider" "by_id" {
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
- `description` (String) Description of this LLM configuration... Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `llm_type` (String) The type of LLM provider (OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAICompatible, etc.).. Computed.
- `api_key` (String) The API key for the LLM provider. Required for OpenAI, Azure OpenAI, Anthropic, Groq, and Mistral... Computed.
- `model_name` (String) The name of the model to use (e.g., gpt-4, claude-3-opus, llama2)... Computed.
- `base_url` (String) The base URL for the LLM API. Required for Azure OpenAI and Ollama, optional for others... Computed.
- `additional_params` (String) Optional JSON object with extra parameters sent directly to the provider API. These are merged last and override any defaults... Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_default` (Bool) Is this the default LLM provider for the project? When set, the global LLM provider will not be used... Computed.
- `cost_per_million_tokens_in_usd_cents` (Number) Cost per million tokens in USD cents. Used for billing when using global LLM providers... Computed.
