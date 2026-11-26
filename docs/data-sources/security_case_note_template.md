---
page_title: "oneuptime_security_case_note_template Data Source - oneuptime"
subcategory: "Security"
description: |-
  Ready-made text for the notes a team writes on its security cases, picked from a case's note composer.
---

# oneuptime_security_case_note_template (Data Source)

Ready-made text for the notes a team writes on its security cases, picked from a case's note composer.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case note template may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_note_template" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_security_case_note_template" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who added the template. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `note` (String) The text the template puts in a case's note, in Markdown, to edit before it is posted.
- `template_description` (String) What the template is for, and when to use it.
- `template_name` (String) What people pick the template by in a case's note composer.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of the project this note template belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
