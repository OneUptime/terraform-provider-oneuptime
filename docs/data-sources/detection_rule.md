---
page_title: "oneuptime_detection_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Sigma detection rules evaluated against security events. Matches create alerts and detection findings.
---

# oneuptime_detection_rule (Data Source)

Sigma detection rules evaluated against security events. Matches create alerts and detection findings.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one detection rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_detection_rule" "example" {
  name = "Example detection rule"
}

# Or by id:
data "oneuptime_detection_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_severity_id` (String) ID of the alert severity for alerts opened by this rule. The ID of a `oneuptime_alert_severity`.
- `created_by_user_id` (String) ID of the user who created this detection rule. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of what this detection rule looks for.
- `distinct_count_field` (String) Optional security-event field (e.g. principalUser, principalIp) whose distinct values are counted instead of raw matching events. The match count threshold then applies to that distinct count. Empty values are not counted. Names that are not typed event columns are looked up in the event's attributes map.
- `evaluation_interval_in_minutes` (Number) How often this rule is evaluated, in minutes. The evaluation window covers the time since the previous evaluation.
- `group_by_field` (String) Optional security-event field (e.g. principalHost, principalUser) to group matches by. One alert is opened per distinct value; empty groups all matches into one alert.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_severity_id` (String) ID of the incident severity for incidents opened by this rule. The ID of a `oneuptime_incident_severity`.
- `is_enabled` (Boolean) Whether this detection rule is evaluated.
- `last_error` (String) The most recent evaluation error, if any. Cleared on the next successful evaluation.
- `match_count_threshold` (Number) Fire only when a group's count — distinct values when a distinct count field is set, matching events otherwise — reaches this number within one evaluation window. 1 fires on any match.
- `name` (String) Friendly name for this detection rule.
- `should_create_alert` (Boolean) Whether matches open OneUptime alerts.
- `should_create_incident` (Boolean) Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per rule.
- `should_write_detection_finding` (Boolean) Whether matches also write a Detection Finding security event back into the events table.
- `sigma_rule_yaml` (String) The Sigma rule to evaluate, in YAML. detection selections and condition are compiled to a ClickHouse query over security events.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_evaluated_at` (String) When the detection engine last evaluated this rule. Null means it has never run.
- `last_match_at` (String) When this rule most recently matched security events. Null means it has never matched.
- `project_id` (String) ID of the project this detection rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
