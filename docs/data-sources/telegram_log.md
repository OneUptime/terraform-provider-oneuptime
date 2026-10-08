---
page_title: "oneuptime_telegram_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all the Telegram messages sent out to all users and subscribers for this project.
---

# oneuptime_telegram_log (Data Source)

Logs of all the Telegram messages sent out to all users and subscribers for this project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one telegram log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_telegram_log" "example" {
  to_chat_id = "example-to-chat-id"
}

# Or by id:
data "oneuptime_telegram_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of Alert associated with this message (if any). The ID of a `oneuptime_alert`.
- `from_bot_username` (String) OneUptime Telegram bot username the message was sent from.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of Incident associated with this message (if any). The ID of a `oneuptime_incident`.
- `message_text` (String) Text content of the Telegram message.
- `monitor_id` (String) ID of Monitor associated with this message (if any). The ID of a `oneuptime_monitor`.
- `on_call_duty_policy_escalation_rule_id` (String) ID of On-Call Duty Policy Escalation Rule associated with this message (if any). The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of On-Call Duty Policy associated with this message (if any). The ID of a `oneuptime_on_call_policy`.
- `on_call_duty_policy_schedule_id` (String) ID of On-Call Duty Policy Schedule associated with this message (if any). The ID of a `oneuptime_on_call_policy_schedule`.
- `scheduled_maintenance_id` (String) ID of Scheduled Maintenance associated with this message (if any). The ID of a `oneuptime_scheduled_maintenance_event`.
- `status` (String) Status of the Telegram message sent.
- `status_message` (String) Status Message (if any).
- `status_page_announcement_id` (String) ID of Status Page Announcement associated with this message (if any). The ID of a `oneuptime_status_page_announcement`.
- `status_page_id` (String) ID of Status Page associated with this message (if any). The ID of a `oneuptime_status_page`.
- `team_id` (String) ID of Team associated with this message (if any). The ID of a `oneuptime_team`.
- `telegram_cost_in_usd_cents` (Number) Telegram Message Cost in USD Cents.
- `telegram_message_id` (String) Message ID returned by Telegram Bot API.
- `to_chat_id` (String) Telegram Chat ID the message was sent to.
- `user_id` (String) ID of User who initiated this message (if any). The ID of a `oneuptime_user` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
