---
page_title: "oneuptime_service_ai_note Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  What a service's team tells OneUptime AI about it. Every investigation of an incident or an alert on the service reads its notes.
---

# oneuptime_service_ai_note (Resource)

What a service's team tells OneUptime AI about it. Every investigation of an incident or an alert on the service reads its notes.

## Example Usage

```terraform
resource "oneuptime_service_ai_note" "example" {
  service_id = oneuptime_service.example.id
  note       = "Checkout's database failover takes five minutes; don't page for it."
}
```

## Schema

### Required

- `note` (String) What the team knows about the service that OneUptime AI should read when it investigates it, in plain text, at most 2000 characters.
- `service_id` (String) ID of the service the note is about. The ID of a `oneuptime_service`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who wrote the note. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this note belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing service ai note by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_service_ai_note.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_service_ai_note.example <id>
```
