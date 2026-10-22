---
page_title: "oneuptime_postmortem_status Resource - oneuptime"
subcategory: "Other"
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
- `order` (Number) Where this status sits in the project's list of postmortem statuses: 1 is the top. A postmortem in the completed status, or in a status below it, counts as done; one above it is still owed. So the not started, draft and completed statuses have to stay in that order. A new status without a number goes just above the completed status. Setting a number moves the status to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.

### Read-Only

- `archived_at` (String) When this postmortem status was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
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
