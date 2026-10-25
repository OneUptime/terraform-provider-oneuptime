---
page_title: "oneuptime_security_case_status Data Source - oneuptime"
subcategory: "Other"
description: |-
  The statuses a security case moves through (New, Investigating, Contained, Closed and your own). Add, rename, reorder or archive them.
---

# oneuptime_security_case_status (Data Source)

The statuses a security case moves through (New, Investigating, Contained, Closed and your own). Add, rename, reorder or archive them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one security case status may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_security_case_status" "example" {
  name = "Example security case status"
}

# Or by id:
data "oneuptime_security_case_status" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `archived_by_user_id` (String) ID of the person who archived this status. The ID of a `oneuptime_user` (see the data source).
- `created_by_user_id` (String) ID of the person who added this status. The ID of a `oneuptime_user` (see the data source).
- `description` (String) What a case in this status is waiting for, in a sentence.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_acknowledged_status` (Boolean) Is it the status a reopened case moves to? From it down, a case counts as taken on. A project has one, and it can't be deleted or archived.
- `is_archived` (Boolean) Archived statuses stay on the cases in them, but can't be picked again. The built-in statuses can't be archived.
- `is_closed_status` (Boolean) Is it the status Close Case moves a case to? From it down, a case counts as closed and carries a verdict. A project has one.
- `is_contained_status` (Boolean) Is it the status the time to contain is measured to? From it down, a case counts as contained. A project has one.
- `is_new_status` (Boolean) Is it the status every new case starts in? A project has one, and it can't be deleted or archived.
- `name` (String) What the status is called, like Investigating or Triage.
- `order` (Number) Where this status sits in the project's list of security case statuses: 1 is the top. A case in a status at or below the acknowledged (contained, closed) status counts as acknowledged (contained, closed), so the new, acknowledged, contained and closed statuses have to stay in that order. A new status without a number goes just above the closed status. Setting a number moves the status to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.
- `slug` (String) Friendly globally unique name for the status.

### Read-Only

- `archived_at` (String) When this status was archived. Empty while it is not archived.
- `color` (String) The status's colour, in hex (#32a852 for example).
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of the project this status belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
