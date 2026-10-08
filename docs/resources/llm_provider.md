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
  name        = "Example llm provider"
  llm_type    = "Example short text"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `llm_type` (String) The type of LLM provider (OpenAI, Azure OpenAI, Anthropic, Groq, Mistral, Ollama, OpenAICompatible, etc.).
- `name` (String) A friendly name for this LLM configuration.

### Optional

- `additional_params` (String) Optional JSON object with extra parameters sent directly to the provider API. These are merged last and override any defaults. A JSON value: write it with `jsonencode()`.
- `api_key` (String) The API key for the LLM provider. Required for OpenAI, Azure OpenAI, Anthropic, Groq, and Mistral.
- `base_url` (String) The base URL for the LLM API. Required for Azure OpenAI and Ollama, optional for others.
- `description` (String) Description of this LLM configuration.
- `is_default` (Boolean) Is this the default LLM provider for the project? When set, the global LLM provider will not be used. Defaults to `false`.
- `model_name` (String) The name of the model to use (e.g., gpt-4, claude-3-opus, llama2).

### Read-Only

- `cost_per_million_tokens_in_usd_cents` (Number) Cost per million tokens in USD cents. Used for billing when using global LLM providers.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this LLM belongs to. If null, it is a global LLM.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing llm provider by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_llm_provider.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_llm_provider.example <id>
```
