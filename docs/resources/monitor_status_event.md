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

- `ends_at` (String) When did this status change end? OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent.
- `root_cause` (String) What is the root cause of this status change?
- `starts_at` (String) When did this status change?

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_monitoring_off` (Boolean) Was the monitor's monitoring turned off for this status event? Uptime, status pages and SLOs leave this time out.
- `is_owner_notified` (Boolean) Are owners notified of status change?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_change_log` (String) A JSON value: read it with `jsondecode()`.
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
