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
  alert_episode_id              = oneuptime_alert_episode.example.id
  feed_info_in_markdown         = "# Heading\n\nThis is **markdown** content"
  alert_episode_feed_event_type = "Example short text"
  display_color                 = "#ff0000"
}
```

## Schema

### Required

- `alert_episode_feed_event_type` (String) Alert Episode Feed Event Type.
- `alert_episode_id` (String) Relation to Alert Episode ID in which this resource belongs. The ID of a `oneuptime_alert_episode`.
- `display_color` (String) Display color for the alert episode log.
- `feed_info_in_markdown` (String) Log of the entire alert episode activity in Markdown.

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

Import an existing alert episode feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_episode_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_episode_feed.example <id>
```
