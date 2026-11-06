---
page_title: "oneuptime_openstack_cloud_feed Data Source - oneuptime"
subcategory: "Infrastructure"
description: |-
  Log of everything that happened to this OpenStack cloud - creation, updates, owner changes and the rules that made them.
---

# oneuptime_openstack_cloud_feed (Data Source)

Log of everything that happened to this OpenStack cloud - creation, updates, owner changes and the rules that made them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one openstack cloud feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_open_stack_cloud_feed` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_openstack_cloud_feed" "example" {
  open_stack_cloud_id = oneuptime_openstack_cloud.example.id
  user_email          = "alice@example.com"
}

# Or by id:
data "oneuptime_openstack_cloud_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the OpenStack cloud change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) More information in Markdown.
- `open_stack_cloud_feed_event_type` (String) OpenStack Cloud Feed Event.
- `open_stack_cloud_id` (String) Relation to OpenStack Cloud ID in which this resource belongs. The ID of a `oneuptime_openstack_cloud`.
- `user_email` (String) Look it up by the email address of the person in `user_id`, instead of their ID: someone in this project, or invited to it. Give one of the two.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source). Or look them up by email address, in `user_email`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for this feed item.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
