---
page_title: "oneuptime_service_ai_note Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  What a service's team tells OneUptime AI about it. Every investigation of an incident or an alert on the service reads its notes.
---

# oneuptime_service_ai_note (Data Source)

What a service's team tells OneUptime AI about it. Every investigation of an incident or an alert on the service reads its notes.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one service ai note may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_service_ai_note" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_service_ai_note" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who wrote the note. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `note` (String) What the team knows about the service that OneUptime AI should read when it investigates it, in plain text, at most 2000 characters.
- `service_id` (String) ID of the service the note is about. The ID of a `oneuptime_service`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this note belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
