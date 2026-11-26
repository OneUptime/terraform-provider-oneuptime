---
page_title: "oneuptime_postmortem_status Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage the statuses an incident's postmortem moves through (Not Started, Draft, Completed for example). Add, edit, reorder or archive them.
---

# oneuptime_postmortem_status (Data Source)

Manage the statuses an incident's postmortem moves through (Not Started, Draft, Completed for example). Add, edit, reorder or archive them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one postmortem status may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_postmortem_status" "example" {
  name = "Example postmortem status"
}

# Or by id:
data "oneuptime_postmortem_status" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Friendly description that will help you remember.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_archived` (Boolean) Archived statuses stay on the incidents whose postmortems are in them, but can't be picked again. The not started, draft and completed statuses can't be archived.
- `is_completed_status` (Boolean) Is it the status a postmortem moves to once it is published on a status page? A postmortem in this status, or in one below it, counts as done. A project has one, and it can't be deleted or archived.
- `is_draft_status` (Boolean) Is it the status a postmortem moves to once its note is written? A project has one, and it can't be deleted or archived.
- `is_not_started_status` (Boolean) Is it the status a new incident's postmortem starts in? A project has one, and it can't be deleted or archived.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Any friendly name of this object.
- `order` (Number) Where this status sits in the project's list of postmortem statuses: 1 is the top. A postmortem in the completed status, or in a status below it, counts as done; one above it is still owed. So the not started, draft and completed statuses have to stay in that order. A new status without a number goes just above the completed status. Setting a number places the status there and moves no other: two at one number are read the older first. In the dashboard, drag the rows to reorder them.
- `slug` (String) Friendly globally unique name for your object.

### Read-Only

- `archived_at` (String) When this postmortem status was archived. Empty while it is not archived.
- `color` (String) Color of this resource in Hex (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
