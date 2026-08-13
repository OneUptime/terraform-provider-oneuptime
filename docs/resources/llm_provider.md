---
page_title: "oneuptime_llm_provider Resource - oneuptime"
subcategory: "Other"
description: |-
  Manage LLM Provider configurations. Connect to OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAI-compatible servers (e.g. vLLM, LocalAI), or other LLM providers to enable AI features.
---

# oneuptime_llm_provider (Resource)

Manage LLM Provider configurations. Connect to OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAI-compatible servers (e.g. vLLM, LocalAI), or other LLM providers to enable AI features.

## Example Usage

```terraform
resource "oneuptime_llm_provider" "example" {
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  llm_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name object.
- `llm_type` (String) The type of LLM provider (OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAICompatible, etc.)..

### Optional

- `description` (String) Description of this LLM configuration...
- `api_key` (String) The API key for the LLM provider. Required for OpenAI, Azure OpenAI, Anthropic, Groq, and Mistral...
- `model_name` (String) The name of the model to use (e.g., gpt-4, claude-3-opus, llama2)...
- `base_url` (String) The base URL for the LLM API. Required for Azure OpenAI and Ollama, optional for others...
- `additional_params` (String) Optional JSON object with extra parameters sent directly to the provider API. These are merged last and override any defaults...
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_default` (Bool) Is this the default LLM provider for the project? When set, the global LLM provider will not be used...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `cost_per_million_tokens_in_usd_cents` (Number) Cost per million tokens in USD cents. Used for billing when using global LLM providers...

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_llm_provider.example <id>
```
