---
page_title: "oneuptime_ai_agent Data Source - oneuptime"
subcategory: "Other"
description: |-
  Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management.
---

# oneuptime_ai_agent (Data Source)

Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_ai_agent" "by_name" {
  name = "example-ai_agent"
}

data "oneuptime_ai_agent" "by_id" {
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
- `key` (String) Ai agent key. Computed.
- `description` (String) Ai agent description. Computed.
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `ai_agent_version` (String) Version object. Computed.
- `last_alive` (String) A date time object.. Computed.
- `icon_file_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `connection_status` (String) Connection Status of the AI Agent.. Computed.
- `is_default` (Bool) Is this the default AI Agent for the project? When set, this agent will be used for automated tasks... Computed.
- `labels` (Set) Relation to Labels Array where this object is categorized in... Computed.
