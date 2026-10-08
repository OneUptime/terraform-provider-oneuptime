---
page_title: "oneuptime_service_feed Resource - oneuptime"
subcategory: "Telemetry & Dashboards"
description: |-
  Log of everything that happened to this service - creation, updates, owner changes and the rules that made them.
---

# oneuptime_service_feed (Resource)

Log of everything that happened to this service - creation, updates, owner changes and the rules that made them.

## Example Usage

```terraform
resource "oneuptime_service_feed" "example" {
  service_id              = oneuptime_service.example.id
  feed_info_in_markdown   = "# Heading\n\nThis is **markdown** content"
  service_feed_event_type = "Example short text"
  display_color           = "#ff0000"
}
```

## Schema

### Required

- `display_color` (String) Display color for this feed item.
- `feed_info_in_markdown` (String) Log of the service change in Markdown.
- `service_feed_event_type` (String) Service Feed Event.
- `service_id` (String) Relation to Service ID in which this resource belongs. The ID of a `oneuptime_service`.

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

Import an existing service feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_service_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_service_feed.example <id>
```
