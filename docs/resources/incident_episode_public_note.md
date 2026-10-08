---
page_title: "oneuptime_incident_episode_public_note Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage public notes for your incident episode
---

# oneuptime_incident_episode_public_note (Resource)

Manage public notes for your incident episode

## Example Usage

```terraform
resource "oneuptime_incident_episode_public_note" "example" {
  incident_episode_id = oneuptime_incident_episode.example.id
}
```

## Schema

### Required

- `incident_episode_id` (String) Relation to Incident Episode ID in which this resource belongs. The ID of a `oneuptime_incident_episode`.

### Optional

- `attachments` (Set of String) Files attached to this note. IDs of `oneuptime_file` resources.
- `note` (String) Notes in markdown.
- `posted_at` (String) Date and time when the note was posted.
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.
- `should_status_page_subscribers_be_notified_on_note_created` (Boolean) Should subscribers be notified about this note? If left out, this follows the episode: true when subscribers were notified that the episode was created, false when it was created without notifying them.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_message_on_note_updated` (String) Status message for the notification sent to subscribers when this note was last updated - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_on_note_created` (String) Status of notification sent to subscribers about this note.
- `subscriber_notification_status_on_note_updated` (String) Status of the notification sent to subscribers when this note was last updated. Empty until an update notification is requested.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing incident episode public note by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_episode_public_note.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_episode_public_note.example <id>
```
