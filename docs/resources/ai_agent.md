---
page_title: "oneuptime_ai_agent Resource - oneuptime"
subcategory: "Other"
description: |-
  Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management.
---

# oneuptime_ai_agent (Resource)

Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management.

## Example Usage

```terraform
resource "oneuptime_ai_agent" "example" {
  key = "Example short text"
  name = jsonencode({
    "_type": "Name",
    "value": "John Doe"
  })
  ai_agent_version = jsonencode({
    "_type": "Version",
    "value": "1.0.0"
  })
  description = "This is a description of the item"
}
```

## Schema

### Required

- `key` (String) Ai agent key.
- `name` (String) Name object.
- `ai_agent_version` (String) Version object.

### Optional

- `description` (String) Ai agent description.
- `last_alive` (String) A date time object..
- `icon_file_id` (String) A unique identifier for an object, represented as a UUID..
- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_default` (Bool) Is this the default AI Agent for the project? When set, this agent will be used for automated tasks...
- `labels` (Set) Relation to Labels Array where this object is categorized in...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `slug` (String) Friendly globally unique name for your object..
- `connection_status` (String) Connection Status of the AI Agent..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_ai_agent.example <id>
```
