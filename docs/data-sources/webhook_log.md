---
page_title: "oneuptime_webhook_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all the outbound Webhook requests sent for this project.
---

# oneuptime_webhook_log (Data Source)

Logs of all the outbound Webhook requests sent for this project.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one webhook log may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_webhook_log" "example" {
  webhook_url = "example-webhook-url"
}

# Or by id:
data "oneuptime_webhook_log" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of Alert associated with this request (if any). The ID of a `oneuptime_alert`.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of Incident associated with this request (if any). The ID of a `oneuptime_incident`.
- `monitor_id` (String) ID of Monitor associated with this request (if any). The ID of a `oneuptime_monitor`.
- `on_call_duty_policy_escalation_rule_id` (String) ID of On-Call Duty Policy Escalation Rule associated with this request (if any). The ID of a `oneuptime_escalation_rule`.
- `on_call_duty_policy_id` (String) ID of On-Call Duty Policy associated with this request (if any). The ID of a `oneuptime_on_call_policy`.
- `on_call_duty_policy_schedule_id` (String) ID of On-Call Duty Policy Schedule associated with this request (if any). The ID of a `oneuptime_on_call_policy_schedule`.
- `request_body` (String) JSON body that was POSTed to the webhook URL.
- `response_body` (String) Response body returned by the webhook endpoint (truncated).
- `response_status_code` (Number) HTTP status code returned by the webhook endpoint.
- `scheduled_maintenance_id` (String) ID of Scheduled Maintenance associated with this request (if any). The ID of a `oneuptime_scheduled_maintenance_event`.
- `status` (String) Status of the Webhook request.
- `status_message` (String) Status Message (if any).
- `status_page_announcement_id` (String) ID of Status Page Announcement associated with this request (if any). The ID of a `oneuptime_status_page_announcement`.
- `status_page_id` (String) ID of Status Page associated with this request (if any). The ID of a `oneuptime_status_page`.
- `team_id` (String) ID of Team associated with this request (if any). The ID of a `oneuptime_team`.
- `user_id` (String) ID of User who initiated this request (if any). The ID of a `oneuptime_user` (see the data source).
- `webhook_url` (String) URL the request was sent to.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
