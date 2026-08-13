---
page_title: "oneuptime_webhook_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all the outbound Webhook requests sent for this project.
---

# oneuptime_webhook_log (Data Source)

Logs of all the outbound Webhook requests sent for this project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_webhook_log" "by_name" {
  name = "example-webhook_log"
}

data "oneuptime_webhook_log" "by_id" {
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
- `webhook_url` (String) URL the request was sent to.. Computed.
- `request_body` (String) JSON body that was POSTed to the webhook URL.. Computed.
- `response_status_code` (Number) HTTP status code returned by the webhook endpoint.. Computed.
- `response_body` (String) Response body returned by the webhook endpoint (truncated).. Computed.
- `status_message` (String) Status Message (if any).. Computed.
- `status` (String) Status of the Webhook request.. Computed.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `scheduled_maintenance_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_announcement_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
