---
page_title: "oneuptime_status_page_announcement Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage announcements on your status page
---

# oneuptime_status_page_announcement (Resource)

Manage announcements on your status page

## Example Usage

```terraform
resource "oneuptime_status_page_announcement" "example" {
  title                = "Example short text"
  show_announcement_at = "2030-01-01T00:00:00Z"
  description          = "Managed by Terraform"
}
```

## Schema

### Required

- `description` (String) Text of the announcement. This can be in Markdown format.
- `show_announcement_at` (String) When should this announcement be shown?
- `title` (String) Title of this resource.

### Optional

- `attachments` (Set of String) Files attached to this announcement. IDs of `oneuptime_file` resources.
- `end_announcement_at` (String) When should this announcement hidden?
- `monitors` (Set of String) List of monitors affected by this announcement. If none are selected, all subscribers will be notified. IDs of `oneuptime_monitor` resources.
- `should_status_page_subscribers_be_notified` (Boolean) Should subscribers be notified about this announcement? Defaults to `true`.
- `status_pages` (Set of String) Status Pages to show show this announcement on. IDs of `oneuptime_status_page` resources.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_owner_notified` (Boolean) Are owners notified of this announcement?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `subscriber_notification_status` (String)
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_message_on_announcement_updated` (String) Status message for the notification sent to subscribers when this announcement was last updated - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_on_announcement_updated` (String) Status of the notification sent to subscribers when this announcement was last updated. Empty until an update notification is requested.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page announcement by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_announcement.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_announcement.example <id>
```
