---
page_title: "oneuptime_alert_episode_feed Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Log of the entire alert episode activity. This is a log of all the episode state changes, alerts added/removed, notes, etc.
---

# oneuptime_alert_episode_feed (Resource)

Log of the entire alert episode activity. This is a log of all the episode state changes, alerts added/removed, notes, etc.

## Example Usage

```terraform
resource "oneuptime_alert_episode_feed" "example" {
  alert_episode_id = "123e4567-e89b-12d3-a456-426614174000"
  feed_info_in_markdown = "# Heading

This is **markdown** content"
  alert_episode_feed_event_type = "Example short text"
  display_color = jsonencode({
    "_type": "Color",
    "value": "#ff0000"
  })
}
```

## Schema

### Required

- `alert_episode_id` (String) A unique identifier for an object, represented as a UUID..
- `feed_info_in_markdown` (String) Log of the entire alert episode activity in Markdown..
- `alert_episode_feed_event_type` (String) Alert Episode Feed Event Type..
- `display_color` (String) Color object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `more_information_in_markdown` (String) More information in Markdown..
- `user_id` (String) A unique identifier for an object, represented as a UUID..
- `posted_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_alert_episode_feed.example <id>
```
