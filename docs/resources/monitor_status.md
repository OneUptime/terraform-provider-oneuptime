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

<a id="order"></a>
## Order

`priority` is this monitor status's place among the monitor statuses of its project, lowest first. Setting it places this monitor status and moves no other, so a configuration that numbers every monitor status of a project is applied exactly as written, whatever order the provider's writes land in, and the next plan is empty. For example, over a `for_each` of `local.monitor_statuses`:

```terraform
priority = index(local.monitor_statuses, each.key) + 1
```

- Give each monitor status a number of its own: two monitor statuses at one number are shown the older first.
- A monitor status created without `priority` goes just above the offline status: it takes the offline status's number, and the offline status and the monitor statuses right below it step down one. Give it `priority` when the configuration numbers the others, or the next plan shows theirs changed.
- Deleting a monitor status leaves a gap in the numbers, which changes nothing: monitor statuses go by their numbers, not by counting.
- Moving a monitor status in the OneUptime Dashboard numbers the project's monitor statuses 1..n again; the next plan shows that as changes to `priority`.


## Schema

### Required

- `color` (String) Color of this resource in Hex (#32a852 for example).
- `name` (String) Any friendly name of this object.

### Optional

- `description` (String) Friendly description that will help you remember.
- `is_archived` (Boolean) Archived statuses stay on the monitors that are in them, but can't be picked again. The operational and offline statuses can't be archived. Defaults to `false`.
- `is_offline_state` (Boolean) Is this monitor in offline state? Defaults to `false`.
- `is_operational_state` (Boolean) Is this monitor in operational state? Defaults to `false`.
- `priority` (Number) Where this status sits in the project's list of monitor statuses, from the healthiest (the lowest number) down to the worst: where monitors are shown together, as on a status page or in a monitor group, the status furthest down the list wins. A new status without a number goes just above the offline status. Setting a number places the status there and moves no other: two at one number are read the older first. In the dashboard, drag the rows to reorder them. See [Order](#order).

### Read-Only

- `archived_at` (String) When this monitor status was archived. Empty while it is not archived.
- `archived_by_user_id` (String) User ID who archived this object (if this object was archived by a User). The ID of a `oneuptime_user` (see the data source).
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
