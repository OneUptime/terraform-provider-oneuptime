---
page_title: "oneuptime_monitor_feed Data Source - oneuptime"
subcategory: "Monitors"
description: |-
  Log of the entire monitor state change. This is a log of all the monitor state changes, public notes, more etc.
---

# oneuptime_monitor_feed (Data Source)

Log of the entire monitor state change. This is a log of all the monitor state changes, public notes, more etc.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one monitor feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_monitor_feed" "example" {
  monitor_id = oneuptime_monitor.example.id
}

# Or by id:
data "oneuptime_monitor_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the entire monitor state change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `monitor_feed_event_type` (String) Monitor Feed Event.
- `monitor_id` (String) Relation to Monitor ID in which this resource belongs. The ID of a `oneuptime_monitor`.
- `more_information_in_markdown` (String) More information in Markdown.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for the monitor log.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
