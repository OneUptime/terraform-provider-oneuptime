---
page_title: "oneuptime_network_site_link Data Source - oneuptime"
subcategory: "Other"
description: |-
  Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor.
---

# oneuptime_network_site_link (Data Source)

Explicit links between Network Sites (data center to region WAN links for example), optionally colored by the status of a Monitor.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network site link may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_site_link" "example" {
  name = "Example network site link"
}

# Or by id:
data "oneuptime_network_site_link" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `from_site_id` (String) ID of the Network Site this link starts from. The ID of a `oneuptime_network_site`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `monitor_id` (String) ID of the Monitor whose status colors this link on map views. The ID of a `oneuptime_monitor`.
- `name` (String) Friendly name for this link.
- `to_site_id` (String) ID of the Network Site this link ends at. The ID of a `oneuptime_network_site`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
