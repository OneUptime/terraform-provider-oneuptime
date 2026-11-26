---
page_title: "oneuptime_postmortem_action_item Resource - oneuptime"
subcategory: "Incidents"
description: |-
  What a postmortem says the team will do so an incident does not happen again: who owns it, when it is due, and whether it is done.
---

# oneuptime_postmortem_action_item (Resource)

What a postmortem says the team will do so an incident does not happen again: who owns it, when it is due, and whether it is done.

## Example Usage

```terraform
resource "oneuptime_postmortem_action_item" "example" {
  incident_id = oneuptime_incident.example.id
  title       = "Alert on connection pool usage above 80%"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `incident_id` (String) ID of the incident whose postmortem this action item is from. The ID of a `oneuptime_incident`.
- `title` (String) What the team will do, in a line.

### Optional

- `description` (String) More about the action item, in Markdown: why it matters, and how anyone will know it is done.
- `due_date` (String) The day the action item is due: the start of that day in your time zone. An open action item is overdue once the day is over.
- `monitor_id` (String) ID of the monitor created for this action item. Empty until one is. The ID of a `oneuptime_monitor`.
- `owner_user_email` (String) The email address of the person in `owner_user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `owner_user_id` (String) ID of the person who owns the action item. Empty while nobody does. The ID of a `oneuptime_user` (see the data source). Or name them by email address, in `owner_user_email`.
- `status` (String) Open, Done or WontDo. Done or WontDo records who closed it and when; Open again clears both. Allowed values: `Open`, `Done`, `WontDo`. Defaults to `Open`.

### Read-Only

- `closed_at` (String) When the action item was marked Done or Won't Do. Empty while it is open.
- `closed_by_user_id` (String) ID of the person who marked the action item Done or Won't Do. The ID of a `oneuptime_user` (see the data source).
- `code_task_id` (String) The code task OneUptime AI works on for this action item: a pull request, or a test, it writes. Empty until someone asks for one.
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the person who added the action item. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of the project this action item belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing postmortem action item by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_postmortem_action_item.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_postmortem_action_item.example <id>
```
