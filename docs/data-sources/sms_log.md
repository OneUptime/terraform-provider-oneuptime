---
page_title: "oneuptime_sms_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all the SMS sent out to all users and subscribers for this project.
---

# oneuptime_sms_log (Data Source)

Logs of all the SMS sent out to all users and subscribers for this project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one sms log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_sms_log" "example" {
  sms_text = "example-sms-text"
}

# Or by id:
data "oneuptime_sms_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of Alert associated with this SMS (if any). The ID of a `oneuptime_alert`.
- `error_code` (String) Error code returned by the SMS provider (e.g. Twilio error code 30007 for carrier filtering) when the message could not be delivered.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of Incident associated with this SMS (if any). The ID of a `oneuptime_incident`.
- `monitor_id` (String) ID of Monitor associated with this SMS (if any). The ID of a `oneuptime_monitor`.
- `on_call_duty_policy_escalation_rule_id` (String) ID of On-Call Duty Policy Escalation Rule associated with this SMS (if any). The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of On-Call Duty Policy associated with this SMS (if any). The ID of a `oneuptime_on_call_policy`.
- `on_call_duty_policy_schedule_id` (String) ID of On-Call Duty Policy Schedule associated with this SMS (if any). The ID of a `oneuptime_on_call_policy_schedule`.
- `scheduled_maintenance_id` (String) ID of Scheduled Maintenance associated with this SMS (if any). The ID of a `oneuptime_scheduled_maintenance_event`.
- `sms_cost_in_usd_cents` (Number) SMS Cost in USD Cents.
- `sms_text` (String) Text content of the message.
- `status` (String) Status of the SMS sent.
- `status_message` (String) Status Message (if any).
- `status_page_announcement_id` (String) ID of Status Page Announcement associated with this SMS (if any). The ID of a `oneuptime_status_page_announcement`.
- `status_page_id` (String) ID of Status Page associated with this SMS (if any). The ID of a `oneuptime_status_page`.
- `team_id` (String) ID of Team associated with this SMS (if any). The ID of a `oneuptime_team`.
- `user_id` (String) ID of User who initiated this SMS (if any). The ID of a `oneuptime_user` (see the data source).
- `user_on_call_log_timeline_id` (String) ID of the on-call notification timeline entry this SMS was sent for (if any). The ID of a `oneuptime_user_on_call_log_timeline` (see the data source).

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `from_number` (String) Phone Number SMS was sent from.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `to_number` (String) Phone Number SMS was sent to.
- `updated_at` (String) Date and Time when the object was updated.
