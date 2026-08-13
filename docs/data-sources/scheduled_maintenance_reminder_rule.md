---
page_title: "oneuptime_scheduled_maintenance_reminder_rule Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Configure reminder rules to periodically notify scheduled maintenance event owners while an event is still not complete
---

# oneuptime_scheduled_maintenance_reminder_rule (Data Source)

Configure reminder rules to periodically notify scheduled maintenance event owners while an event is still not complete Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_scheduled_maintenance_reminder_rule" "by_name" {
  name = "example-scheduled_maintenance_reminder_rule"
}

data "oneuptime_scheduled_maintenance_reminder_rule" "by_id" {
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
- `description` (String) Description of this reminder rule.. Computed.
- `order` (Number) Order/priority of this rule. Rules are evaluated in order (lowest first). First matching rule wins... Computed.
- `is_enabled` (Bool) Whether this reminder rule is enabled.. Computed.
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind scheduled maintenance event owners while the event is still not complete. For example, set to 30 to remind owners every 30 minutes... Computed.
- `stop_reminders_on_state` (String) Stop sending reminders once the scheduled maintenance event reaches this state. Select Ongoing to stop reminders when the event starts, or Completed to keep reminding until the event is completed... Computed.
- `remind_while_scheduled` (Bool) Send reminders while the event is still scheduled (before it starts). When disabled, reminders only begin once the event has started... Computed.
- `labels` (Set) Only apply this reminder rule to scheduled maintenance events with these labels. Leave empty to match all events... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
