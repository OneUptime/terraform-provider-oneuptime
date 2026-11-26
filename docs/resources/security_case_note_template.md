---
page_title: "oneuptime_security_case_note_template Resource - oneuptime"
subcategory: "Security"
description: |-
  Ready-made text for the notes a team writes on its security cases, picked from a case's note composer.
---

# oneuptime_security_case_note_template (Resource)

Ready-made text for the notes a team writes on its security cases, picked from a case's note composer.

## Example Usage

```terraform
resource "oneuptime_security_case_note_template" "example" {
  template_name = "Containment checklist"
  note          = "# Heading\n\nThis is **markdown** content"
}
```

## Schema

### Required

- `note` (String) The text the template puts in a case's note, in Markdown, to edit before it is posted.
- `template_name` (String) What people pick the template by in a case's note composer.

### Optional

- `template_description` (String) What the template is for, and when to use it.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who added the template. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this note template belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case note template by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case_note_template.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case_note_template.example <id>
```
