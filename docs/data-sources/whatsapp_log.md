---
page_title: "oneuptime_whatsapp_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all the WhatsApp messages sent out to all users and subscribers for this project.
---

# oneuptime_whatsapp_log (Data Source)

Logs of all the WhatsApp messages sent out to all users and subscribers for this project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one whatsapp log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_whats_app_log` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_whatsapp_log" "example" {
  message_text = "example-message-text"
}

# Or by id:
data "oneuptime_whatsapp_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of Alert associated with this message (if any). The ID of a `oneuptime_alert`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of Incident associated with this message (if any). The ID of a `oneuptime_incident`.
- `message_text` (String) Text content of the WhatsApp message.
- `monitor_id` (String) ID of Monitor associated with this message (if any). The ID of a `oneuptime_monitor`.
- `on_call_duty_policy_escalation_rule_id` (String) ID of On-Call Duty Policy Escalation Rule associated with this message (if any). The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of On-Call Duty Policy associated with this message (if any). The ID of a `oneuptime_on_call_policy`.
- `on_call_duty_policy_schedule_id` (String) ID of On-Call Duty Policy Schedule associated with this message (if any). The ID of a `oneuptime_on_call_policy_schedule`.
- `scheduled_maintenance_id` (String) ID of Scheduled Maintenance associated with this message (if any). The ID of a `oneuptime_scheduled_maintenance_event`.
- `status` (String) Status of the WhatsApp message sent.
- `status_message` (String) Status Message (if any).
- `status_page_announcement_id` (String) ID of Status Page Announcement associated with this message (if any). The ID of a `oneuptime_status_page_announcement`.
- `status_page_id` (String) ID of Status Page associated with this message (if any). The ID of a `oneuptime_status_page`.
- `team_id` (String) ID of Team associated with this message (if any). The ID of a `oneuptime_team`.
- `user_id` (String) ID of User who initiated this message (if any). The ID of a `oneuptime_user` (see the data source).
- `whats_app_cost_in_usd_cents` (Number) WhatsApp Message Cost in USD Cents.
- `whats_app_message_id` (String) Message ID returned by Meta's API.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `from_number` (String) Phone Number WhatsApp message was sent from (Business Number ID).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `to_number` (String) Phone Number WhatsApp message was sent to.
- `updated_at` (String) Date and Time when the object was updated.
