---
page_title: "oneuptime_scheduled_maintenance_state_timeline Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Change state of your scheduled maintenance event.
---

# oneuptime_scheduled_maintenance_state_timeline (Data Source)

Change state of your scheduled maintenance event.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled maintenance state timeline may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_state_timeline" "example" {
  scheduled_maintenance_id = oneuptime_scheduled_maintenance_event.example.id
}

# Or by id:
data "oneuptime_scheduled_maintenance_state_timeline" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_owner_notified` (Boolean) Are owners notified of state change?
- `scheduled_maintenance_id` (String) ID of Scheduled Maintenance this resource belongs to. The ID of a `oneuptime_scheduled_maintenance_event`.
- `scheduled_maintenance_state_id` (String) Scheduled Maintenance State ID. Which state does this event belongs to? The ID of a `oneuptime_scheduled_maintenance_state`.
- `should_status_page_subscribers_be_notified` (Boolean) Should subscribers be notified about this state change?
- `subscriber_notification_status` (String) Status of notification sent to subscribers about this scheduled maintenance state change.
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `ends_at` (String) When did this status change end?
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `starts_at` (String) When did this status change? Correct this when the recorded time is wrong - every measurement derived from this timeline is recomputed from the corrected value.
- `updated_at` (String) Date and Time when the object was updated.
