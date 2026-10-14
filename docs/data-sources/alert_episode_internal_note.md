---
page_title: "oneuptime_alert_episode_internal_note Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Manage internal notes for your alert episodes
---

# oneuptime_alert_episode_internal_note (Data Source)

Manage internal notes for your alert episodes

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one alert episode internal note may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_alert_episode_internal_note" "example" {
  alert_episode_id = oneuptime_alert_episode.example.id
}

# Or by id:
data "oneuptime_alert_episode_internal_note" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_episode_id` (String) Relation to Alert Episode ID in which this resource belongs. The ID of a `oneuptime_alert_episode`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `note` (String) Notes in markdown.
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.

### Read-Only

- `attachments` (Set of String) Files attached to this note. IDs of `oneuptime_file` resources.
- `created_at` (String) Date and Time when the object was created.
- `posted_at` (String) Date and time the note belongs to in the timeline: when it was written, or - for a note saved from a Slack or Microsoft Teams message - when that message was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `saved_from_chat_message` (String) When this note was saved from a Slack or Microsoft Teams message - pinned, or given a note emoji - the chat it was posted in (workspaceType), who wrote it (authorName, and authorUserId when their chat account is connected to OneUptime) and when it was posted there (postedAt). Written by OneUptime; the note's creator is whoever saved it. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
