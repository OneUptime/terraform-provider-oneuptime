---
page_title: "oneuptime_domain Resource - oneuptime"
subcategory: "Organization"
description: |-
  Manage Custom Domains for your project
---

# oneuptime_domain (Resource)

Manage Custom Domains for your project

## Example Usage

```terraform
resource "oneuptime_domain" "example" {
  domain = "example.com"
}
```

## Schema

### Required

- `domain` (String) Domain - acmeinc.com for example.

### Optional

- `is_verified` (Boolean) Is this domain verified? OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent. Defaults to `false`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `domain_verification_text` (String) Verification text that you need to add to your domains TXT record to veify the domain.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing domain by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_domain.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_domain.example <id>
```
