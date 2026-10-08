---
page_title: "oneuptime_slo_feed Data Source - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this Service Level Objective - configuration changes, status transitions, burn rate alerts and incidents, monitor rule changes and owner changes.
---

# oneuptime_slo_feed (Data Source)

Log of everything that happened to this Service Level Objective - configuration changes, status transitions, burn rate alerts and incidents, monitor rule changes and owner changes.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one slo feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_slo_feed" "example" {
  service_level_objective_id = oneuptime_service_level_objective.example.id
}

# Or by id:
data "oneuptime_slo_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the Service Level Objective change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) More information in Markdown.
- `service_level_objective_feed_event_type` (String) Service Level Objective Feed Event.
- `service_level_objective_id` (String) Relation to Service Level Objective ID in which this resource belongs. The ID of a `oneuptime_service_level_objective`.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for this feed item.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
