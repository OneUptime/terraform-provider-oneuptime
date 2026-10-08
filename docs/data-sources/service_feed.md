---
page_title: "oneuptime_service_feed Data Source - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Log of everything that happened to this service - creation, updates, owner changes and the rules that made them.
---

# oneuptime_service_feed (Data Source)

Log of everything that happened to this service - creation, updates, owner changes and the rules that made them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one service feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_service_feed" "example" {
  service_id = oneuptime_service.example.id
}

# Or by id:
data "oneuptime_service_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the service change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) More information in Markdown.
- `service_feed_event_type` (String) Service Feed Event.
- `service_id` (String) Relation to Service ID in which this resource belongs. The ID of a `oneuptime_service`.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for this feed item.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
