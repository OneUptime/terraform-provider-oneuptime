---
page_title: "oneuptime_llm_provider Data Source - oneuptime"
subcategory: "Other"
description: |-
  Manage LLM Provider configurations. Connect to OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAI-compatible servers (e.g. vLLM, LocalAI), or other LLM providers to enable AI features.
---

# oneuptime_llm_provider (Data Source)

Manage LLM Provider configurations. Connect to OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAI-compatible servers (e.g. vLLM, LocalAI), or other LLM providers to enable AI features.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one llm provider may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_llm_provider" "example" {
  name = "Example llm provider"
}

# Or by id:
data "oneuptime_llm_provider" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `api_key` (String) The API key for the LLM provider. Required for OpenAI, Azure OpenAI, Anthropic, Groq, and Mistral.
- `base_url` (String) The base URL for the LLM API. Required for Azure OpenAI and Ollama, optional for others. The API key and the Additional Parameters are sent to it, so only project owners and admins can change it. Everyone who may read the project's settings can read it, so never put a key, a token or a password in it: use the API Key.
- `cost_per_million_tokens_in_usd_cents` (Number) Cost per million tokens in USD cents. Used for billing when using global LLM providers.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `description` (String) Description of this LLM configuration.
- `has_additional_params` (Boolean) Whether Additional Parameters are saved. Worked out from the parameters on every read.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_default` (Boolean) Is this the default LLM provider for the project? When set, the global LLM provider will not be used.
- `llm_type` (String) The type of LLM provider (OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAICompatible, etc.).
- `model_name` (String) The name of the model to use (e.g., gpt-4, claude-3-opus, llama2).
- `name` (String) A friendly name for this LLM configuration.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `additional_params` (String) Optional JSON object with extra parameters sent directly to the provider API. These are merged last and override any defaults. Read only by project owners and admins, like the API key. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this LLM belongs to. If null, it is a global LLM.
- `updated_at` (String) Date and Time when the object was updated.
