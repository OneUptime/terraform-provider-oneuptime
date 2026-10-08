---
page_title: "oneuptime_cloud_resource_feed Resource - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this cloud resource - creation, updates, owner changes and the rules that made them.
---

# oneuptime_cloud_resource_feed (Resource)

Log of everything that happened to this cloud resource - creation, updates, owner changes and the rules that made them.

## Example Usage

```terraform
resource "oneuptime_cloud_resource_feed" "example" {
  cloud_resource_id              = oneuptime_cloud_resource.example.id
  feed_info_in_markdown          = "# Heading\n\nThis is **markdown** content"
  cloud_resource_feed_event_type = "Example short text"
  display_color                  = "#ff0000"
}
```

## Schema

### Required

- `cloud_resource_feed_event_type` (String) Cloud Resource Feed Event.
- `cloud_resource_id` (String) Relation to Cloud Resource ID in which this resource belongs. The ID of a `oneuptime_cloud_resource`.
- `display_color` (String) Display color for this feed item.
- `feed_info_in_markdown` (String) Log of the cloud resource change in Markdown.

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

Import an existing cloud resource feed by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_cloud_resource_feed.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_cloud_resource_feed.example <id>
```
