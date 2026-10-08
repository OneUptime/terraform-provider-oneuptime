---
page_title: "oneuptime_change_event Data Source - oneuptime"
subcategory: "Other"
description: |-
  API endpoints for Change Event
---

# oneuptime_change_event (Data Source)

API endpoints for Change Event

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one change event may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_change_event" "example" {
  primary_entity_id = "example-primary-entity-id"
}

# Or by id:
data "oneuptime_change_event" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `attributes` (String) Attributes.
- `description` (String) Description.
- `event_type` (String) Event Type.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `primary_entity_id` (String) Service ID.
- `primary_entity_type` (String) Service Type.
- `time` (String) Time.
- `title` (String) Title.

### Read-Only

- `attribute_keys` (Set of String) Attribute Keys.
- `project_id` (String) Project ID.
