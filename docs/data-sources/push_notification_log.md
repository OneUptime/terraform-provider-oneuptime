---
page_title: "oneuptime_push_notification_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all the Push Notifications sent out to all users and subscribers for this project.
---

# oneuptime_push_notification_log (Data Source)

Logs of all the Push Notifications sent out to all users and subscribers for this project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_push_notification_log" "by_name" {
  name = "example-push_notification_log"
}

data "oneuptime_push_notification_log" "by_id" {
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
- `title` (String) Title of the push notification.. Computed.
- `body` (String) Body of the push notification.. Computed.
- `device_type` (String) Type of device this was sent to (e.g., web).. Computed.
- `device_name` (String) Name of the device this was sent to.. Computed.
- `status_message` (String) Status Message (if any).. Computed.
- `status` (String) Status of the push notification.. Computed.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `scheduled_maintenance_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_announcement_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
