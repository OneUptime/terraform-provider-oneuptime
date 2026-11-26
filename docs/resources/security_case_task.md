---
page_title: "oneuptime_security_case_task Resource - oneuptime"
subcategory: "Security"
description: |-
  What has to be done on a security case, who is doing it, and whether it is done.
---

# oneuptime_security_case_task (Resource)

What has to be done on a security case, who is doing it, and whether it is done.

## Example Usage

```terraform
resource "oneuptime_security_case_task" "example" {
  security_case_id = oneuptime_security_case.example.id
  title            = "Reset the user's password and end their sessions"
}
```

## Schema

### Required

- `security_case_id` (String) ID of the case this task is on. The ID of a `oneuptime_security_case`.
- `title` (String) What has to be done, in a line.

### Optional

- `assigned_to_user_email` (String) The email address of the person in `assigned_to_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `assigned_to_user_id` (String) ID of the person doing the task. Empty while nobody is. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `assigned_to_user_email`.
- `is_done` (Boolean) Whether the task is done. Setting it records who did it and when; clearing it opens the task again. Defaults to `false`.
- `response_step_key` (String) The step of the case's response plan this task was made from, when it was made from one. Set when the task is added.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who added the task. The ID of a `oneuptime_user` (see the data source).
- `done_at` (String) When the task was marked done. Empty while it is open.
- `done_by_user_id` (String) ID of the person who marked the task done. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this task belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case task by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case_task.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case_task.example <id>
```
