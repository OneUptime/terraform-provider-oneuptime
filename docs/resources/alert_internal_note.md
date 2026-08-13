---
page_title: "oneuptime_alert_internal_note Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Manage internal notes for your alert
---

# oneuptime_alert_internal_note (Resource)

Manage internal notes for your alert

## Example Usage

```terraform
resource "oneuptime_alert_internal_note" "example" {
  alert_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `alert_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `note` (String) Notes in markdown..
- `attachments` (Set) Files attached to this note..
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `is_owner_notified` (Bool) Are owners notified of this resource ownership?..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_alert_internal_note.example <id>
```
