---
page_title: "oneuptime_scheduled_maintenance_state_timeline Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Change state of your scheduled maintenance event.
---

# oneuptime_scheduled_maintenance_state_timeline (Resource)

Change state of your scheduled maintenance event.

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_state_timeline" "example" {
  scheduled_maintenance_id = "123e4567-e89b-12d3-a456-426614174000"
  scheduled_maintenance_state_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `scheduled_maintenance_id` (String) A unique identifier for an object, represented as a UUID..
- `scheduled_maintenance_state_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `should_status_page_subscribers_be_notified` (Bool) Should subscribers be notified about this state change?..
- `ends_at` (String) A date time object..
- `starts_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `subscriber_notification_status` (String) Status of notification sent to subscribers about this scheduled maintenance state change..
- `subscriber_notification_status_message` (String) Status message for subscriber notifications - includes success messages, failure reasons, or skip reasons..
- `is_owner_notified` (Bool) Are owners notified of state change?..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_scheduled_maintenance_state_timeline.example <id>
```
