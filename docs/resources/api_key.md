---
page_title: "oneuptime_api_key Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  Manage API Keys for your project
---

# oneuptime_api_key (Resource)

Manage API Keys for your project

## Example Usage

```terraform
resource "oneuptime_api_key" "example" {
  name        = "Example api key"
  expires_at  = "2030-01-01T00:00:00Z"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `expires_at` (String) Date and Time when this API Key expires.
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.

### Read-Only

- `api_key` (String, Sensitive) The secret key, sent in the ApiKey header of each request. Shown once: in the answer to the create that makes it, or to a reset of it. OneUptime keeps only its hash. Kept in state from the answer to the create: OneUptime never reads it back, so a resource imported into Terraform has none. Replace the resource for a new one.
- `api_key_hint` (String) The key's last four characters, to tell it from the project's other keys. Never enough to use the key.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing api key by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_api_key.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_api_key.example <id>
```
