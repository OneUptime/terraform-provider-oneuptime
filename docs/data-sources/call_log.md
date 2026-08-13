---
page_title: "oneuptime_call_log Data Source - oneuptime"
subcategory: "Logs & Metrics"
description: |-
  Logs of all the Call sent out to all users and subscribers for this project.
---

# oneuptime_call_log (Data Source)

Logs of all the Call sent out to all users and subscribers for this project. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_call_log" "by_name" {
  name = "example-call_log"
}

data "oneuptime_call_log" "by_id" {
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
- `call_data` (String) Content of the data that was sent in the call.. Computed.
- `status_message` (String) Status Message (if any).. Computed.
- `status` (String) Status of the Call sent.. Computed.
- `call_cost_in_usd_cents` (Number) Call Cost in USD Cents.. Computed.
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
- `team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
