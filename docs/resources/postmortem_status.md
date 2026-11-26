---
page_title: "oneuptime_postmortem_status Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage the statuses an incident's postmortem moves through (Not Started, Draft, Completed for example). Add, edit, reorder or archive them.
---

# oneuptime_postmortem_status (Resource)

Manage the statuses an incident's postmortem moves through (Not Started, Draft, Completed for example). Add, edit, reorder or archive them.

## Example Usage

```terraform
resource "oneuptime_postmortem_status" "example" {
  name        = "Example postmortem status"
  color       = "#ff0000"
  description = "Managed by Terraform"
}
```

<a id="order"></a>
## Order

`order` is this postmortem status's place among the postmortem statuses of its project, lowest first. Setting it places this postmortem status and moves no other, so a configuration that numbers every postmortem status of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.postmortem_statuses`:

```terraform
order = index(local.postmortem_statuses, each.key) + 1
```

- Give each postmortem status a number of its own: two postmortem statuses at one number are shown the older first.
- A postmortem status created without `order` goes just above the completed status: it takes the completed status's number, and the completed status and the postmortem statuses right below it step down one. Give it `order` when the configuration numbers the others, or the next plan shows theirs changed.
- The not started status, draft status and completed status stay in that order, from the top: a write that would put two of them at one number, or out of order, is refused. An apply that moves them past each other's numbers can meet that refusal on one write; apply again and it finishes.
- Deleting a postmortem status leaves a gap in the numbers, which changes nothing: postmortem statuses go by their numbers, not by counting.
- Moving a postmortem status in the OneUptime Dashboard numbers the project's postmortem statuses 1..n again; the next plan shows that as changes to `order`.


## Schema

### Required

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_archived` (Boolean) Archived statuses stay on the incidents whose postmortems are in them, but can't be picked again. The not started, draft and completed statuses can't be archived. Defaults to `false`.
- `is_completed_status` (Boolean) Is it the status a postmortem moves to once it is published on a status page? A postmortem in this status, or in one below it, counts as done. A project has one, and it can't be deleted or archived.
- `is_draft_status` (Boolean) Is it the status a postmortem moves to once its note is written? A project has one, and it can't be deleted or archived.
- `is_not_started_status` (Boolean) Is it the status a new incident's postmortem starts in? A project has one, and it can't be deleted or archived.
- `order` (Number) Where this status sits in the project's list of postmortem statuses: 1 is the top. A postmortem in the completed status, or in a status below it, counts as done; one above it is still owed. So the not started, draft and completed statuses have to stay in that order. A new status without a number goes just above the completed status. Setting a number places the status there and moves no other: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `archived_at` (String) When this postmortem status was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing postmortem status by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_postmortem_status.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_postmortem_status.example <id>
```
