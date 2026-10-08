---
page_title: "oneuptime_scheduled_maintenance_reminder_rule Data Source - oneuptime"
subcategory: "Scheduled Maintenance"
description: |-
  Configure reminder rules to periodically notify scheduled maintenance event owners while an event is still not complete
---

# oneuptime_scheduled_maintenance_reminder_rule (Data Source)

Configure reminder rules to periodically notify scheduled maintenance event owners while an event is still not complete

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one scheduled maintenance reminder rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_scheduled_maintenance_reminder_rule" "example" {
  name = "Example scheduled maintenance reminder rule"
}

# Or by id:
data "oneuptime_scheduled_maintenance_reminder_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this reminder rule.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this reminder rule is enabled.
- `name` (String) Name of this reminder rule.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `remind_while_scheduled` (Boolean) Send reminders while the event is still scheduled (before it starts). When disabled, reminders only begin once the event has started.
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind scheduled maintenance event owners while the event is still not complete. For example, set to 30 to remind owners every 30 minutes.
- `stop_reminders_on_state` (String) Stop sending reminders once the scheduled maintenance event reaches this state. Select Ongoing to stop reminders when the event starts, or Completed to keep reminding until the event is completed.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `labels` (Set of String) Only apply this reminder rule to scheduled maintenance events with these labels. Leave empty to match all events. IDs of `oneuptime_label` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
