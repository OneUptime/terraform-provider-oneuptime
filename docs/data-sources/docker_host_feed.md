---
page_title: "oneuptime_docker_host_feed Data Source - oneuptime"
subcategory: "Other"
description: |-
  Log of everything that happened to this Docker host - creation, updates, owner changes and the rules that made them.
---

# oneuptime_docker_host_feed (Data Source)

Log of everything that happened to this Docker host - creation, updates, owner changes and the rules that made them. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_docker_host_feed" "by_name" {
  name = "example-docker_host_feed"
}

data "oneuptime_docker_host_feed" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `docker_host_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `feed_info_in_markdown` (String) Log of the Docker host change in Markdown.. Computed.
- `more_information_in_markdown` (String) More information in Markdown.. Computed.
- `docker_host_feed_event_type` (String) Docker Host Feed Event.. Computed.
- `display_color` (String) Color object. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `posted_at` (String) A date time object.. Computed.
