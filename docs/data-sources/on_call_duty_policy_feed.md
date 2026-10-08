---
page_title: "oneuptime_on_call_duty_policy_feed Data Source - oneuptime"
subcategory: "On-Call & Escalation"
description: |-
  Log of the entire onCallDutyPolicy state change. This is a log of all the on call duty policy changes.
---

# oneuptime_on_call_duty_policy_feed (Data Source)

Log of the entire onCallDutyPolicy state change. This is a log of all the on call duty policy changes.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one on call duty policy feed may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_on_call_duty_policy_feed" "example" {
  on_call_duty_policy_id = oneuptime_on_call_policy.example.id
}

# Or by id:
data "oneuptime_on_call_duty_policy_feed" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `feed_info_in_markdown` (String) Log of the entire onCallDutyPolicy state change in Markdown.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `more_information_in_markdown` (String) More information in Markdown.
- `on_call_duty_policy_feed_event_type` (String) On Call Duty Policy Feed Event.
- `on_call_duty_policy_id` (String) Relation to OnCallDutyPolicy ID in which this resource belongs. The ID of a `oneuptime_on_call_policy`.
- `user_id` (String) User who this feed belongs to (if this feed belongs to a User). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `display_color` (String) Display color for the onCallDutyPolicy log.
- `posted_at` (String) Date and time when the feed was posted.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
