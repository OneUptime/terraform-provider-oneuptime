---
page_title: "oneuptime_on_call_duty_policy_feed Resource - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Log of the entire onCallDutyPolicy state change. This is a log of all the on call duty policy changes.
---

# oneuptime_on_call_duty_policy_feed (Resource)

Log of the entire onCallDutyPolicy state change. This is a log of all the on call duty policy changes.

## Example Usage

```terraform
resource "oneuptime_on_call_duty_policy_feed" "example" {
  on_call_duty_policy_id = "123e4567-e89b-12d3-a456-426614174000"
  feed_info_in_markdown = "# Heading

This is **markdown** content"
  on_call_duty_policy_feed_event_type = "Example short text"
  display_color = jsonencode({
    "_type": "Color",
    "value": "#ff0000"
  })
}
```

## Schema

### Required

- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID..
- `feed_info_in_markdown` (String) Log of the entire onCallDutyPolicy state change in Markdown..
- `on_call_duty_policy_feed_event_type` (String) On Call Duty Policy Feed Event..
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
terraform import oneuptime_on_call_duty_policy_feed.example <id>
```
