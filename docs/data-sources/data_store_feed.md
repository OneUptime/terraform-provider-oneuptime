---
page_title: "oneuptime_data_store_feed Data Source - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this data store - creation, updates, owner changes and the rules that made them.
---

# oneuptime_data_store_feed (Data Source)

Log of everything that happened to this data store - creation, updates, owner changes and the rules that made them.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one data store feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_data_store_feed" "example" {
  data_store_id = oneuptime_data_store.example.id
}

# Or by id:
data "oneuptime_data_store_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_store_feed_event_type` (String) Data Store Feed Event.
- `data_store_id` (String) Relation to Data Store ID in which this resource belongs. The ID of a `oneuptime_data_store`.
- `feed_info_in_markdown` (String) Log of the data store change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) More information in Markdown.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for this feed item.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
