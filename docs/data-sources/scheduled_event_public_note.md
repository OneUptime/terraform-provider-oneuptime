---
page_title: "oneuptime_scheduled_event_public_note Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Manage public notes for your scheduled event
---

# oneuptime_scheduled_event_public_note (Data Source)

Manage public notes for your scheduled event

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled event public note may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_event_public_note" "example" {
  scheduled_maintenance_id = oneuptime_scheduled_maintenance_event.example.id
}

# Or by id:
data "oneuptime_scheduled_event_public_note" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `approval_status` (String) Where this note stands in the project's public note approval: WaitingForApproval, Approved or Rejected. Empty for a note that needed no approval. On a published note, WaitingForApproval or Rejected is about an edit of it (pendingNote).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `is_published` (Boolean) Whether this note is published: shown on status pages and sent to subscribers. False while it waits for approval, or after it was rejected, when the project asks for public notes to be approved.
- `note` (String) Notes in markdown.
- `pending_note` (String) An edit of this published note that waits for approval, in markdown. Status pages show the note as it was approved until the edit is approved. Empty when no edit waits.
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.
- `posted_with_scheduled_maintenance_state_id` (String) The state the scheduled maintenance event moved to when this note was posted with that state change. Subscribers are told this state with the note. Empty for a note posted on its own. The ID of a `oneuptime_scheduled_maintenance_state`.
- `scheduled_maintenance_id` (String) ID of Scheduled Maintenance this resource belongs to. The ID of a `oneuptime_scheduled_maintenance_event`.
- `should_status_page_subscribers_be_notified_on_note_created` (Boolean) Should subscribers be notified about this note? If left out, this follows the scheduled maintenance event: true when subscribers were notified that the event was created, false when it was created without notifying them.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_message_on_note_updated` (String) Status message for the notification sent to subscribers when this note was last updated - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_on_note_created` (String) Status of notification sent to subscribers about this note.
- `subscriber_notification_status_on_note_updated` (String) Status of the notification sent to subscribers when this note was last updated. Empty until an update notification is requested.

### Read-Only

- `approval` (String) This note's approval so far: when it was sent for approval and by whom (submittedAt, submittedByUserId), whether it is an edit of the published note (isEdit), each step's approval (approvals: stepId, userId, at, comment) and a rejection (rejection: userId, at, reason). A JSON value: read it with `jsondecode()`.
- `attachments` (Set of String) Files attached to this note. IDs of `oneuptime_file` resources.
- `created_at` (String) Date and Time when the object was created.
- `posted_at` (String) Date and time when the note was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `saved_from_chat_message` (String) When this note was saved from a Slack or Microsoft Teams message - pinned, or given a note emoji - the chat it was posted in (workspaceType), who wrote it (authorName, and authorUserId when their chat account is connected to OneUptime) and when it was posted there (postedAt). Written by OneUptime; the note's creator is whoever saved it. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
