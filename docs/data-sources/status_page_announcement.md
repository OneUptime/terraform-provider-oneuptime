---
page_title: "oneuptime_status_page_announcement Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage announcements on your status page
---

# oneuptime_status_page_announcement (Data Source)

Manage announcements on your status page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page announcement may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_announcement" "example" {
  title = "example-title"
}

# Or by id:
data "oneuptime_status_page_announcement" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Text of the announcement. This can be in Markdown format.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this announcement?
- `should_status_page_subscribers_be_notified` (Boolean) Should subscribers be notified about this announcement?
- `subscriber_notification_status` (String) Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Create Status Page Announcement], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Admin, Status Page Member, Status Page Viewer, Read Status Page Announcement], Update: [Project Owner, Project Admin, Project Member, Status Page Admin, Status Page Member, Edit Status Page Announcement]
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_message_on_announcement_updated` (String) Status message for the notification sent to subscribers when this announcement was last updated - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_on_announcement_updated` (String) Status of the notification sent to subscribers when this announcement was last updated. Empty until an update notification is requested.
- `title` (String) Title of this resource.

### Read-Only

- `attachments` (Set of String) Files attached to this announcement. IDs of `oneuptime_file` resources.
- `created_at` (String) Date and Time when the object was created.
- `end_announcement_at` (String) When should this announcement hidden?
- `monitors` (Set of String) List of monitors affected by this announcement. If none are selected, all subscribers will be notified. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `show_announcement_at` (String) When should this announcement be shown?
- `status_pages` (Set of String) Status Pages to show show this announcement on. IDs of `oneuptime_status_page` resources.
- `updated_at` (String) Date and Time when the object was updated.
