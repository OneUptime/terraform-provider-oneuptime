---
page_title: "oneuptime_network_site_status_event Data Source - oneuptime"
subcategory: "Other"
description: |-
  History of the rolled-up health status of a Network Site (Operational to Offline for example)
---

# oneuptime_network_site_status_event (Data Source)

History of the rolled-up health status of a Network Site (Operational to Offline for example)

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one network site status event may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_network_site_status_event" "example" {
  site_id = oneuptime_network_site.example.id
}

# Or by id:
data "oneuptime_network_site_status_event" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `monitor_status_id` (String) Relation to Monitor Status ID Resource in which this object belongs. The ID of a `oneuptime_monitor_status`.
- `site_id` (String) Relation to Network Site ID Resource in which this object belongs. The ID of a `oneuptime_network_site`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ends_at` (String) When did this status change end?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `starts_at` (String) When did this status change?
- `updated_at` (String) Date and Time when the object was updated.
