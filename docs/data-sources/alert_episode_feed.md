---
page_title: "oneuptime_alert_episode_feed Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Log of the entire alert episode activity. This is a log of all the episode state changes, alerts added/removed, notes, etc.
---

# oneuptime_alert_episode_feed (Data Source)

Log of the entire alert episode activity. This is a log of all the episode state changes, alerts added/removed, notes, etc.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one alert episode feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_alert_episode_feed" "example" {
  alert_episode_id = oneuptime_alert_episode.example.id
}

# Or by id:
data "oneuptime_alert_episode_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_episode_feed_event_type` (String) Alert Episode Feed Event Type.
- `alert_episode_id` (String) Relation to Alert Episode ID in which this resource belongs. The ID of a `oneuptime_alert_episode`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the entire alert episode activity in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) More information in Markdown.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for the alert episode log.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
