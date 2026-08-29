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
  cloud_resource_id = "123e4567-e89b-12d3-a456-426614174000"
  feed_info_in_markdown = "# Heading

This is **markdown** content"
  cloud_resource_feed_event_type = "Example short text"
  display_color = jsonencode({
    "_type": "Color",
    "value": "#ff0000"
  })
}
```

## Schema

### Required

- `cloud_resource_id` (String) A unique identifier for an object, represented as a UUID..
- `feed_info_in_markdown` (String) Log of the cloud resource change in Markdown..
- `cloud_resource_feed_event_type` (String) Cloud Resource Feed Event..
- `display_color` (String) Color object.

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `more_information_in_markdown` (String) More information in Markdown..
- `user_id` (String) A unique identifier for an object, represented as a UUID..
- `posted_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_cloud_resource_feed.example <id>
```
