---
page_title: "oneuptime_sms_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all the SMS sent out to all users and subscribers for this project.
---

# oneuptime_sms_log (Data Source)

Logs of all the SMS sent out to all users and subscribers for this project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_sms_log" "by_name" {
  name = "example-sms_log"
}

data "oneuptime_sms_log" "by_id" {
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
- `to_number` (String) Phone object. Computed.
- `from_number` (String) Phone object. Computed.
- `sms_text` (String) Text content of the message.. Computed.
- `status_message` (String) Status Message (if any).. Computed.
- `status` (String) Status of the SMS sent.. Computed.
- `error_code` (String) Error code returned by the SMS provider (e.g. Twilio error code 30007 for carrier filtering) when the message could not be delivered... Computed.
- `sms_cost_in_usd_cents` (Number) SMS Cost in USD Cents.. Computed.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `monitor_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `scheduled_maintenance_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_announcement_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_escalation_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `on_call_duty_policy_schedule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `user_on_call_log_timeline_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
