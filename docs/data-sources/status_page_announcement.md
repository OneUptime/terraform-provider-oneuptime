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
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_status_page_announcement" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Text of the announcement. This can be in Markdown format.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of this announcement?
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `should_status_page_subscribers_be_notified` (Boolean) Should subscribers be notified about this announcement?
- `subscriber_notification_status` (String) Whether the status page's subscribers were told about the announcement: Pending until they are, then Success or Failed. Set by OneUptime.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_message_on_announcement_updated` (String) Status message for the notification sent to subscribers when this announcement was last updated - includes success messages, failure reasons, or skip reasons.
- `subscriber_notification_status_on_announcement_updated` (String) Status of the notification sent to subscribers when this announcement was last updated. Empty until an update notification is requested.
- `title` (String) Title of this resource.

### Read-Only

- `attachments` (Set of String) Files attached to this announcement. IDs of `oneuptime_file` resources.
- `created_at` (String) Date and Time when the object was created.
- `end_announcement_at` (String) When should this announcement hidden?
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `monitors` (Set of String) List of monitors affected by this announcement. If none are selected, all subscribers will be notified. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `show_announcement_at` (String) When should this announcement be shown?
- `status_pages` (Set of String) Status Pages to show show this announcement on. IDs of `oneuptime_status_page` resources.
- `updated_at` (String) Date and Time when the object was updated.
