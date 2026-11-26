---
page_title: "oneuptime_check_setting Resource - oneuptime"
subcategory: "Data Security"
description: |-
  How a project tunes a built-in check: turn it off, or give its failures another severity.
---

# oneuptime_check_setting (Resource)

How a project tunes a built-in check: turn it off, or give its failures another severity.

## Example Usage

```terraform
resource "oneuptime_check_setting" "example" {
  check_id = "Example short text"
}
```

## Schema

### Required

- `check_id` (String) The built-in check this setting changes.

### Optional

- `is_enabled` (Boolean) Whether the check runs. Defaults to `true`.
- `note` (String) Why the check was changed.
- `severity_override` (String) The severity every failure of the check gets, instead of the one its store's sensitivity sets. Allowed values: `Critical`, `High`, `Medium`, `Low`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this check setting. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this check setting belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing check setting by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_check_setting.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_check_setting.example <id>
```
