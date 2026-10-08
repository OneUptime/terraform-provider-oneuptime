---
page_title: "oneuptime_monitor_status Resource - oneuptime"
subcategory: "Monitors"
description: |-
  Manage monitor status in your project. Monitor Status are Operational, Degraded and Offline for example. Add custom status like Monitoring or more.
---

# oneuptime_monitor_status (Resource)

Manage monitor status in your project. Monitor Status are Operational, Degraded and Offline for example. Add custom status like Monitoring or more.

## Example Usage

```terraform
resource "oneuptime_monitor_status" "example" {
  name        = "Example monitor status"
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
- `is_offline_state` (Boolean) Is this monitor in offline state? Defaults to `false`.
- `is_operational_state` (Boolean) Is this monitor in operational state? Defaults to `false`.
- `priority` (Number) Where this status sits in the project's list of monitor statuses, from the healthiest (the lowest number) down to the worst: where monitors are shown together, as on a status page or in a monitor group, the status furthest down the list wins. A new status without a number goes just above the offline status. Setting a number moves the status to that place, and the ones in between shift by one. In the dashboard, drag the rows to reorder them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `slug` (String) Friendly globally unique name for your object.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing monitor status by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_monitor_status.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_monitor_status.example <id>
```
