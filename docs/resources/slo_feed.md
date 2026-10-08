---
page_title: "oneuptime_slo_feed Resource - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this Service Level Objective - configuration changes, status transitions, burn rate alerts and incidents, monitor rule changes and owner changes.
---

# oneuptime_slo_feed (Resource)

Log of everything that happened to this Service Level Objective - configuration changes, status transitions, burn rate alerts and incidents, monitor rule changes and owner changes.

## Example Usage

```terraform
resource "oneuptime_slo_feed" "example" {
  service_level_objective_id              = oneuptime_service_level_objective.example.id
  feed_info_in_markdown                   = "# Heading\n\nThis is **markdown** content"
  service_level_objective_feed_event_type = "Example short text"
  display_color                           = "#ff0000"
}
```

## Schema

### Required

- `display_color` (String) Display color for this feed item.
- `feed_info_in_markdown` (String) Log of the Service Level Objective change in Markdown.
- `service_level_objective_feed_event_type` (String) Service Level Objective Feed Event.
- `service_level_objective_id` (String) Relation to Service Level Objective ID in which this resource belongs. The ID of a `oneuptime_service_level_objective`.

### Optional

- `more_information_in_markdown` (String) More information in Markdown.
- `posted_at` (String) Date and time when the feed was posted.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing slo feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_slo_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_slo_feed.example <id>
```
