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
  title = "Example short text"
  show_announcement_at = "2030-01-01T00:00:00Z"
  description = "# Heading

This is **markdown** content"
}
```

## Schema

### Required

- `title` (String) Title of this resource..
- `show_announcement_at` (String) A date time object..
- `description` (String) Text of the announcement. This can be in Markdown format...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `status_pages` (Set) Status Pages to show show this announcement on...
- `monitors` (Set) List of monitors affected by this announcement. If none are selected, all subscribers will be notified...
- `end_announcement_at` (String) A date time object..
- `attachments` (Set) Files attached to this announcement..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `should_status_page_subscribers_be_notified` (Bool) Should subscribers be notified about this announcement?..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `subscriber_notification_status` (String) Status page announcement subscriber_notification_status.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons..
- `is_owner_notified` (Bool) Are owners notified of this announcement?..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page_announcement.example <id>
```
