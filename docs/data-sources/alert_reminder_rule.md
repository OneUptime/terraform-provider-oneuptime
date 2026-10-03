---
page_title: "oneuptime_alert_reminder_rule Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Configure reminder rules to periodically notify alert owners while an alert is still open
---

# oneuptime_alert_reminder_rule (Data Source)

Configure reminder rules to periodically notify alert owners while an alert is still open Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_alert_reminder_rule" "by_name" {
  name = "example-alert_reminder_rule"
}

data "oneuptime_alert_reminder_rule" "by_id" {
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
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource... Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this reminder rule.. Computed.
- `order` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first, and the first one that matches wins. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them... Computed.
- `is_enabled` (Bool) Whether this reminder rule is enabled.. Computed.
- `reminder_interval_in_minutes` (Number) How often (in minutes) to remind alert owners while the alert is still open. For example, set to 30 to remind owners every 30 minutes... Computed.
- `stop_reminders_on_state` (String) Stop sending reminders once the alert reaches this state. Select Acknowledged to stop reminders when the alert is acknowledged, or Resolved to keep reminding until the alert is resolved... Computed.
- `alert_severities` (Set) Only apply this reminder rule to alerts with these severities. Leave empty to match alerts of any severity... Computed.
- `labels` (Set) Only apply this reminder rule to alerts with these labels. Leave empty to match alerts with any labels... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
