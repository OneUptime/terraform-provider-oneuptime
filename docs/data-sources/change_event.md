---
page_title: "oneuptime_change_event Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Change Event
---

# oneuptime_change_event (Data Source)

API endpoints for Change Event Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_change_event" "by_name" {
  name = "example-change_event"
}

data "oneuptime_change_event" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `project_id` (String) Project ID. Computed.
- `primary_entity_id` (String) Service ID. Computed.
- `primary_entity_type` (String) Service Type. Computed.
- `time` (String) Time. Computed.
- `event_type` (String) Event Type. Computed.
- `title` (String) Title. Computed.
- `description` (String) Description. Computed.
- `attributes` (String) Attributes. Computed.
- `attribute_keys` (Set) Attribute Keys. Computed.
