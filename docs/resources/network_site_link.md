---
page_title: "oneuptime_network_site_link Resource - oneuptime"
subcategory: "Other"
description: |-
  Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor.
---

# oneuptime_network_site_link (Resource)

Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor.

## Example Usage

```terraform
resource "oneuptime_network_site_link" "example" {
  from_site_id = oneuptime_network_site.example.id
  to_site_id   = oneuptime_network_site.example.id
  name         = "Example network site link"
}
```

## Schema

### Required

- `from_site_id` (String) ID of the Network Site this link starts from. The ID of a `oneuptime_network_site`.
- `to_site_id` (String) ID of the Network Site this link ends at. The ID of a `oneuptime_network_site`.

### Optional

- `monitor_id` (String) ID of the Monitor whose status colors this link on map views. The ID of a `oneuptime_monitor`.
- `name` (String) Friendly name for this link.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing network site link by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_network_site_link.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_network_site_link.example <id>
```
