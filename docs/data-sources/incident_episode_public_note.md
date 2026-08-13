---
page_title: "oneuptime_incident_episode_public_note Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Manage public notes for your incident episode
---

# oneuptime_incident_episode_public_note (Data Source)

Manage public notes for your incident episode Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_episode_public_note" "by_name" {
  name = "example-incident_episode_public_note"
}

data "oneuptime_incident_episode_public_note" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_episode_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `note` (String) Notes in markdown.. Computed.
- `attachments` (Set) Files attached to this note.. Computed.
- `subscriber_notification_status_on_note_created` (String) Status of notification sent to subscribers about this note.. Computed.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.. Computed.
- `should_status_page_subscribers_be_notified_on_note_created` (Bool) Should subscribers be notified about this note?.. Computed.
- `is_owner_notified` (Bool) Are owners notified of this resource ownership?.. Computed.
- `posted_at` (String) A date time object.. Computed.
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message... Computed.
