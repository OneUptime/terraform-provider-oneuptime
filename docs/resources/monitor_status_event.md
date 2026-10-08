---
page_title: "oneuptime_monitor_status_event Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Change state of the monitor (Operational to Offline for example)
---

# oneuptime_monitor_status_event (Resource)

Change state of the monitor (Operational to Offline for example)

## Example Usage

```terraform
resource "oneuptime_monitor_status_event" "example" {
  monitor_id        = oneuptime_monitor.example.id
  monitor_status_id = oneuptime_monitor_status.example.id
}
```

## Schema

### Required

- `monitor_id` (String) Relation to Monitor ID Resource in which this object belongs. The ID of a `oneuptime_monitor`.
- `monitor_status_id` (String) Relation to Monitor Status ID Resource in which this object belongs. The ID of a `oneuptime_monitor_status`.

### Optional

- `ends_at` (String) When did this status change end?
- `root_cause` (String) What is the root cause of this status change?
- `starts_at` (String) When did this status change?

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of status change?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_change_log` (String) A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing monitor status event by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_monitor_status_event.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_monitor_status_event.example <id>
```
