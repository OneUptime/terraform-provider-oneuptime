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
  service_id = "123e4567-e89b-12d3-a456-426614174000"
  feed_info_in_markdown = "# Heading

This is **markdown** content"
  service_feed_event_type = "Example short text"
  display_color = jsonencode({
    "_type": "Color",
    "value": "#ff0000"
  })
}
```

## Schema

### Required

- `service_id` (String) A unique identifier for an object, represented as a UUID..
- `feed_info_in_markdown` (String) Log of the service change in Markdown..
- `service_feed_event_type` (String) Service Feed Event..
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
terraform import oneuptime_service_feed.example <id>
```
