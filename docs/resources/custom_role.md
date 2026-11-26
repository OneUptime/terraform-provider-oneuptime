---
page_title: "oneuptime_custom_role Resource - oneuptime"
subcategory: "Teams & Access"
description: |-
  A set of permissions you name once, such as Incident Responder, and give to teams and API keys. Changing a role changes every team and key that has it.
---

# oneuptime_custom_role (Resource)

A set of permissions you name once, such as Incident Responder, and give to teams and API keys. Changing a role changes every team and key that has it.

## Example Usage

```terraform
resource "oneuptime_custom_role" "example" {
  name        = "Example custom role"
  description = "Managed by Terraform"
}
```

<a id="deleting-it"></a>
## Deleting it

OneUptime deletes a custom role that workflows run as only when the delete says what they run as instead, and refuses one that does not, naming the workflows. In Terraform, `workflows_run_as_on_delete` says it:

```terraform
resource "oneuptime_custom_role" "example" {
  name = "Example custom role"

  workflows_run_as_on_delete = "Default"
}
```

or

```terraform
resource "oneuptime_custom_role" "example" {
  name = "Example custom role"

  workflows_run_as_on_delete                = "CustomRole"
  workflows_run_as_custom_role_id_on_delete = oneuptime_custom_role.other.id
}
```

- Set it, and run `terraform apply`, before you destroy the custom role or remove its block. Terraform deletes a resource with what its state holds, so a value set in the same apply that deletes or replaces the custom role is not sent yet.
- Nothing is sent until the delete: setting or changing it never changes the custom role or anything that depends on it.
- Left unset, the delete says nothing more, as before. When OneUptime needs it, the destroy is refused, the custom role stays in state, and the error says what to set.
- OneUptime keeps none of it, so an import leaves it unset: the first plan after an import sets it in state, and changes nothing in OneUptime.
- When this custom role and the one in `workflows_run_as_custom_role_id_on_delete` are destroyed together, Terraform deletes this one first, since it refers to the other, and what depended on it moves to the other custom role. Give that custom role its own `workflows_run_as_on_delete` too, or its delete is refused.

## Schema

### Required

- `name` (String) What the role is called, such as Incident Responder. Unique in the project, and never the name of a built-in role.

### Optional

- `description` (String) What the role is for, in a sentence or two.
- `workflows_run_as_custom_role_id_on_delete` (String) The ID of the other custom role the workflows run as, when `workflows_run_as_on_delete` is "CustomRole": one you could give a team, and not the role being deleted. The ID of a `oneuptime_custom_role`. Sent only when Terraform deletes the custom role, and changes nothing before: set it, and apply, before you destroy the custom role. See [Deleting it](#deleting-it).
- `workflows_run_as_on_delete` (String) What the workflows that run as this custom role run as once it is deleted: "Default", the default permissions, or "CustomRole", another custom role you could give a team, named in `workflows_run_as_custom_role_id_on_delete`. Never Project Owner: deleting a role never raises what a workflow can do. Sent only when Terraform deletes the custom role, and changes nothing before: set it, and apply, before you destroy the custom role. See [Deleting it](#deleting-it). Allowed values: `Default`, `CustomRole`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `permission_count` (Number) How many permissions the role holds. OneUptime keeps it as permissions are added to the role and taken out of it.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing custom role by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_custom_role.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_custom_role.example <id>
```
