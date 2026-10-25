---
page_title: "oneuptime_security_case_status Resource - oneuptime"
subcategory: "Other"
description: |-
  The statuses a security case moves through (New, Investigating, Contained, Closed and your own). Add, rename, reorder or archive them.
---

# oneuptime_security_case_status (Resource)

The statuses a security case moves through (New, Investigating, Contained, Closed and your own). Add, rename, reorder or archive them.

## Example Usage

```terraform
resource "oneuptime_security_case_status" "example" {
  name        = "Example security case status"
  color       = "#ff0000"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `color` (String) The status's colour, in hex (#32a852 for example).
- `name` (String) What the status is called, like Investigating or Triage.

### Optional

- `description` (String) What a case in this status is waiting for, in a sentence.
- `is_acknowledged_status` (Boolean) Is it the status a reopened case moves to? From it down, a case counts as taken on. A project has one, and it can't be deleted or archived.
- `is_archived` (Boolean) Archived statuses stay on the cases in them, but can't be picked again. The built-in statuses can't be archived. Defaults to `false`.
- `is_closed_status` (Boolean) Is it the status Close Case moves a case to? From it down, a case counts as closed and carries a verdict. A project has one.
- `is_contained_status` (Boolean) Is it the status the time to contain is measured to? From it down, a case counts as contained. A project has one.
- `is_new_status` (Boolean) Is it the status every new case starts in? A project has one, and it can't be deleted or archived.
- `order` (Number) Where this status sits in the project's list of security case statuses: 1 is the top. A case in a status at or below the acknowledged (contained, closed) status counts as acknowledged (contained, closed), so the new, acknowledged, contained and closed statuses have to stay in that order. A new status without a number goes just above the closed status. Setting a number moves the status to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.

### Read-Only

- `archived_at` (String) When this status was archived. Empty while it is not archived.
- `archived_by_user_id` (String) ID of the person who archived this status. The ID of a `oneuptime_user` (see the data source).
- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the person who added this status. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of the project this status belongs to. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for the status.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing security case status by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_security_case_status.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_security_case_status.example <id>
```
