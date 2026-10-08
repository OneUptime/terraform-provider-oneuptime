---
page_title: "oneuptime_alert_episode_internal_note Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Manage internal notes for your alert episodes
---

# oneuptime_alert_episode_internal_note (Resource)

Manage internal notes for your alert episodes

## Example Usage

```terraform
resource "oneuptime_alert_episode_internal_note" "example" {
  alert_episode_id = oneuptime_alert_episode.example.id
}
```

## Schema

### Required

- `alert_episode_id` (String) Relation to Alert Episode ID in which this resource belongs. The ID of a `oneuptime_alert_episode`.

### Optional

- `attachments` (Set of String) Files attached to this note. IDs of `oneuptime_file` resources.
- `note` (String) Notes in markdown.
- `posted_from_slack_message_id` (String) Unique identifier for the Slack message this note was created from (channel_id:message_ts). Used to prevent duplicate notes when multiple users react to the same message.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this resource ownership?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert episode internal note by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_episode_internal_note.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_episode_internal_note.example <id>
```
