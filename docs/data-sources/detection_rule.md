---
page_title: "oneuptime_detection_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Sigma detection rules evaluated against security events. Matches create alerts and detection findings.
---

# oneuptime_detection_rule (Data Source)

Sigma detection rules evaluated against security events. Matches create alerts and detection findings. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_detection_rule" "by_name" {
  name = "example-detection_rule"
}

data "oneuptime_detection_rule" "by_id" {
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
- `description` (String) Description of what this detection rule looks for... Computed.
- `sigma_rule_yaml` (String) The Sigma rule to evaluate, in YAML. detection selections and condition are compiled to a ClickHouse query over security events... Computed.
- `is_enabled` (Bool) Whether this detection rule is evaluated... Computed.
- `evaluation_interval_in_minutes` (Number) How often this rule is evaluated, in minutes. The evaluation window covers the time since the previous evaluation... Computed.
- `group_by_field` (String) Optional security-event field (e.g. principalHost, principalUser) to group matches by. One alert is opened per distinct value; empty groups all matches into one alert... Computed.
- `should_create_alert` (Bool) Whether matches open OneUptime alerts... Computed.
- `should_write_detection_finding` (Bool) Whether matches also write a Detection Finding security event back into the events table... Computed.
- `should_create_incident` (Bool) Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per rule... Computed.
- `alert_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `last_evaluated_at` (String) A date time object.. Computed.
- `last_match_at` (String) A date time object.. Computed.
- `last_error` (String) The most recent evaluation error, if any. Cleared on the next successful evaluation... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
