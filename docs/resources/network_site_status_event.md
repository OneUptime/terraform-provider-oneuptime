---
page_title: "oneuptime_network_site_status_event Resource - oneuptime"
subcategory: "Other"
description: |-
  History of the rolled-up health status of a Network Site (Operational to Offline for example)
---

# oneuptime_network_site_status_event (Resource)

History of the rolled-up health status of a Network Site (Operational to Offline for example)

## Example Usage

```terraform
resource "oneuptime_network_site_status_event" "example" {
  site_id           = oneuptime_network_site.example.id
  monitor_status_id = oneuptime_monitor_status.example.id
}
```

## Schema

### Required

- `monitor_status_id` (String) Relation to Monitor Status ID Resource in which this object belongs. The ID of a `oneuptime_monitor_status`.
- `site_id` (String) Relation to Network Site ID Resource in which this object belongs. The ID of a `oneuptime_network_site`.

### Optional

- `ends_at` (String) When did this status change end?
- `starts_at` (String) When did this status change?

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network site status event by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_site_status_event.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_site_status_event.example <id>
```
