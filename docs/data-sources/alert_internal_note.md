---
page_title: "oneuptime_alert_internal_note Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Manage internal notes for your alert
---

# oneuptime_alert_internal_note (Data Source)

Manage internal notes for your alert

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one alert internal note may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_alert_internal_note" "example" {
  alert_id = oneuptime_alert.example.id
}

# Or by id:
data "oneuptime_alert_internal_note" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) Relation to Alert ID in which this resource belongs. The ID of a `oneuptime_alert`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `note` (String) Notes in markdown.
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.

### Read-Only

- `attachments` (Set of String) Files attached to this note. IDs of `oneuptime_file` resources.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
