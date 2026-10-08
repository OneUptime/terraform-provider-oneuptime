---
page_title: "oneuptime_incident_episode_public_note Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage public notes for your incident episode
---

# oneuptime_incident_episode_public_note (Data Source)

Manage public notes for your incident episode

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident episode public note may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_episode_public_note" "example" {
  incident_episode_id = oneuptime_incident_episode.example.id
}

# Or by id:
data "oneuptime_incident_episode_public_note" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_episode_id` (String) Relation to Incident Episode ID in which this resource belongs. The ID of a `oneuptime_incident_episode`.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `note` (String) Notes in markdown.
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.
- `should_status_page_subscribers_be_notified_on_note_created` (Boolean) Should subscribers be notified about this note? If left out, this follows the episode: true when subscribers were notified that the episode was created, false when it was created without notifying them.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_message_on_note_updated` (String) Status message for the notification sent to subscribers when this note was last updated - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_on_note_created` (String) Status of notification sent to subscribers about this note.
- `subscriber_notification_status_on_note_updated` (String) Status of the notification sent to subscribers when this note was last updated. Empty until an update notification is requested.

### Read-Only

- `attachments` (Set of String) Files attached to this note. IDs of `oneuptime_file` resources.
- `created_at` (String) Date and Time when the object was created.
- `posted_at` (String) Date and time when the note was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
