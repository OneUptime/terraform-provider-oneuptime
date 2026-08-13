---
page_title: "oneuptime_alert_internal_note Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Manage internal notes for your alert
---

# oneuptime_alert_internal_note (Data Source)

Manage internal notes for your alert Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_alert_internal_note" "by_name" {
  name = "example-alert_internal_note"
}

data "oneuptime_alert_internal_note" "by_id" {
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
- `alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `note` (String) Notes in markdown.. Computed.
- `attachments` (Set) Files attached to this note.. Computed.
- `is_owner_notified` (Bool) Are owners notified of this resource ownership?.. Computed.
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message... Computed.
