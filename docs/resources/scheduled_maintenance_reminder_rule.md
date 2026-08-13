---
page_title: "oneuptime_scheduled_maintenance_reminder_rule Resource - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Configure reminder rules to periodically notify scheduled maintenance event owners while an event is still not complete
---

# oneuptime_scheduled_maintenance_reminder_rule (Resource)

Configure reminder rules to periodically notify scheduled maintenance event owners while an event is still not complete

## Example Usage

```terraform
resource "oneuptime_scheduled_maintenance_reminder_rule" "example" {
  name = "Example short text"
  reminder_interval_in_minutes = 42
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this reminder rule..
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind scheduled maintenance event owners while the event is still not complete. For example, set to 30 to remind owners every 30 minutes...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this reminder rule..
- `order` (Number) Order/priority of this rule. Rules are evaluated in order (lowest first). First matching rule wins...
- `is_enabled` (Bool) Whether this reminder rule is enabled..
- `stop_reminders_on_state` (String) Stop sending reminders once the scheduled maintenance event reaches this state. Select Ongoing to stop reminders when the event starts, or Completed to keep reminding until the event is completed...
- `remind_while_scheduled` (Bool) Send reminders while the event is still scheduled (before it starts). When disabled, reminders only begin once the event has started...
- `labels` (Set) Only apply this reminder rule to scheduled maintenance events with these labels. Leave empty to match all events...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_scheduled_maintenance_reminder_rule.example <id>
```
