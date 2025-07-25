---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "oneuptime_incident_feed Resource - oneuptime"
subcategory: ""
description: |-
  Incident feed resource
---

# oneuptime_incident_feed (Resource)

Incident feed resource

## Example Usage

```terraform
resource "oneuptime_incident_feed" "example" {
  incident_id = "123e4567-e89b-12d3-a456-426614174000"
  feed_info_in_markdown = "example-feed_info_in_markdown"
  incident_feed_event_type = "123e4567-e89b-12d3-a456-426614174000"
  display_color = {
    id = "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

## Schema

- `id` (String) Unique identifier for the resource. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Required.
- `feed_info_in_markdown` (String) Log (in Markdown). Required.
- `more_information_in_markdown` (String) More Information (in Markdown). Optional.
- `incident_feed_event_type` (String) Incident Feed Event. Required.
- `display_color` (Map) Color object. Required.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Optional.
- `posted_at` (Map) A date time object.. Optional.
- `created_at` (Map) A date time object.. Computed.
- `updated_at` (Map) A date time object.. Computed.
- `deleted_at` (Map) A date time object.. Computed.
- `version` (Number) Version. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_feed.example <id>
```
