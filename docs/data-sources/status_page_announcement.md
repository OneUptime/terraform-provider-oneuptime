---
page_title: "oneuptime_status_page_announcement Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Manage announcements on your status page
---

# oneuptime_status_page_announcement (Data Source)

Manage announcements on your status page Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page_announcement" "by_name" {
  name = "example-status_page_announcement"
}

data "oneuptime_status_page_announcement" "by_id" {
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
- `status_pages` (Set) Status Pages to show show this announcement on... Computed.
- `monitors` (Set) List of monitors affected by this announcement. If none are selected, all subscribers will be notified... Computed.
- `title` (String) Title of this resource.. Computed.
- `show_announcement_at` (String) A date time object.. Computed.
- `end_announcement_at` (String) A date time object.. Computed.
- `description` (String) Text of the announcement. This can be in Markdown format... Computed.
- `attachments` (Set) Files attached to this announcement.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `subscriber_notification_status` (String) Status page announcement subscriber_notification_status. Computed.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.. Computed.
- `should_status_page_subscribers_be_notified` (Bool) Should subscribers be notified about this announcement?.. Computed.
- `is_owner_notified` (Bool) Are owners notified of this announcement?.. Computed.
