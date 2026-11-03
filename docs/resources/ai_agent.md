---
page_title: "oneuptime_ai_agent Resource - oneuptime"
subcategory: "AI"
description: |-
  Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management.
---

# oneuptime_ai_agent (Resource)

Manages custom AI agents. Deploy AI agents anywhere and connect them to your project for automated incident management.

## Example Usage

```terraform
resource "oneuptime_ai_agent" "example" {
  name        = "Example ai agent"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Any friendly name of this object.

### Optional

- `ai_agent_version` (String) A version label for the AI agent. Leave it out and it is 1.0.0.
- `description` (String) Friendly description that will help you remember.
- `icon_file_id` (String) AI Agent Icon File ID. The ID of a `oneuptime_file`.
- `is_default` (Boolean) Is this the default AI Agent for the project? When set, this agent will be used for automated tasks. Defaults to `false`.
- `key` (String, Sensitive) The secret the AI agent connects with, its AI_AGENT_KEY. Leave it out and OneUptime generates one.
- `labels` (Set of String) Labels to group and filter by. IDs of `oneuptime_label` resources.
- `last_alive` (String) When the AI agent last checked in. Set by OneUptime.

### Read-Only

- `connection_status` (String) Connection Status of the AI Agent.
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project the AI agent belongs to. A global AI agent belongs to none.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing ai agent by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_ai_agent.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_ai_agent.example <id>
```
