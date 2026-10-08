---
page_title: "oneuptime_scheduled_maintenance_feed Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Log of the entire scheduled maintenance state change. This is a log of all the scheduled maintenance state changes, public notes, more etc.
---

# oneuptime_scheduled_maintenance_feed (Data Source)

Log of the entire scheduled maintenance state change. This is a log of all the scheduled maintenance state changes, public notes, more etc.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled maintenance feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_feed" "example" {
  scheduled_maintenance_id = oneuptime_scheduled_maintenance_event.example.id
}

# Or by id:
data "oneuptime_scheduled_maintenance_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the entire scheduled maintenance state change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) More information in Markdown.
- `scheduled_maintenance_feed_event_type` (String) ScheduledMaintenance Log Event.
- `scheduled_maintenance_id` (String) Relation to ScheduledMaintenance ID in which this resource belongs. The ID of a `oneuptime_scheduled_maintenance_event`.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for this log.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
