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

### Read-Only

- `connection_status` (String) Connection Status of the AI Agent.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User).
- `id` (String) Unique identifier for the resource.
- `last_alive` (String) When the AI agent last checked in. Set by OneUptime.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
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
