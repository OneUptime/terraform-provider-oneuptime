---
page_title: "oneuptime_whats_app_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all the WhatsApp messages sent out to all users and subscribers for this project.
---

# oneuptime_whats_app_log (Data Source)

Logs of all the WhatsApp messages sent out to all users and subscribers for this project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_whats_app_log" "by_name" {
  name = "example-whats_app_log"
}

data "oneuptime_whats_app_log" "by_id" {
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
- `message_text` (String) Text content of the WhatsApp message.. Computed.
- `status_message` (String) Status Message (if any).. Computed.
- `whats_app_message_id` (String) Message ID returned by Meta's API.. Computed.
- `status` (String) Status of the WhatsApp message sent.. Computed.
- `whats_app_cost_in_usd_cents` (Number) WhatsApp Message Cost in USD Cents.. Computed.
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
