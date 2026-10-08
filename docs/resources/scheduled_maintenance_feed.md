---
page_title: "oneuptime_scheduled_maintenance_feed Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Log of the entire scheduled maintenance state change. This is a log of all the scheduled maintenance state changes, public notes, more etc.
---

# oneuptime_scheduled_maintenance_feed (Resource)

Log of the entire scheduled maintenance state change. This is a log of all the scheduled maintenance state changes, public notes, more etc.

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_feed" "example" {
  scheduled_maintenance_id              = oneuptime_scheduled_maintenance_event.example.id
  feed_info_in_markdown                 = "# Heading\n\nThis is **markdown** content"
  scheduled_maintenance_feed_event_type = "Example short text"
  display_color                         = "#ff0000"
}
```

## Schema

### Required

- `display_color` (String) Display color for this log.
- `feed_info_in_markdown` (String) Log of the entire scheduled maintenance state change in Markdown.
- `scheduled_maintenance_feed_event_type` (String) ScheduledMaintenance Log Event.
- `scheduled_maintenance_id` (String) Relation to ScheduledMaintenance ID in which this resource belongs. The ID of a `oneuptime_scheduled_maintenance_event`.

### Optional

- `more_information_in_markdown` (String) More information in Markdown.
- `posted_at` (String) Date and time when the feed was posted.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing scheduled maintenance feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_scheduled_maintenance_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_scheduled_maintenance_feed.example <id>
```
