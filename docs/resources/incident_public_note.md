---
page_title: "oneuptime_incident_public_note Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Manage public notes for your incident
---

# oneuptime_incident_public_note (Resource)

Manage public notes for your incident

## Example Usage

```terraform
resource "oneuptime_incident_public_note" "example" {
  incident_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `incident_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `note` (String) Notes in markdown..
- `attachments` (Set) Files attached to this note..
- `should_status_page_subscribers_be_notified_on_note_created` (Bool) Should subscribers be notified about this note?..
- `posted_at` (String) A date time object..
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `subscriber_notification_status_on_note_created` (String) Status of notification sent to subscribers about this note..
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons..
- `is_owner_notified` (Bool) Are owners notified of this resource ownership?..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_public_note.example <id>
```
