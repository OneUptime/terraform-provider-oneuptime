---
page_title: "oneuptime_detection_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Sigma detection rules evaluated against security events. Matches create alerts and detection findings.
---

# oneuptime_detection_rule (Resource)

Sigma detection rules evaluated against security events. Matches create alerts and detection findings.

## Example Usage

```terraform
resource "oneuptime_detection_rule" "example" {
  name            = "Example detection rule"
  sigma_rule_yaml = "This is an example of very long text content that might be stored in this field. It can contain a lot of information, such as detailed descriptions, comments, or any other lengthy text data that needs to be stored in the database."
  description     = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Friendly name for this detection rule.
- `sigma_rule_yaml` (String) The Sigma rule to evaluate, in YAML. detection selections and condition are compiled to a ClickHouse query over security events.

### Optional

- `alert_severity_id` (String) ID of the alert severity for alerts opened by this rule. The ID of a `oneuptime_alert_severity`.
- `description` (String) Description of what this detection rule looks for.
- `distinct_count_field` (String) Optional security-event field (e.g. principalUser, principalIp) whose distinct values are counted instead of raw matching events. The match count threshold then applies to that distinct count. Empty values are not counted. Names that are not typed event columns are looked up in the event's attributes map.
- `evaluation_interval_in_minutes` (Number) How often this rule is evaluated, in minutes. The evaluation window covers the time since the previous evaluation. Defaults to `1`.
- `group_by_field` (String) Optional security-event field (e.g. principalHost, principalUser) to group matches by. One alert is opened per distinct value; empty groups all matches into one alert.
- `incident_severity_id` (String) ID of the incident severity for incidents opened by this rule. The ID of a `oneuptime_incident_severity`.
- `is_enabled` (Boolean) Whether this detection rule is evaluated. Defaults to `true`.
- `match_count_threshold` (Number) Fire only when a group's count — distinct values when a distinct count field is set, matching events otherwise — reaches this number within one evaluation window. 1 fires on any match. Defaults to `1`.
- `should_create_alert` (Boolean) Whether matches open OneUptime alerts. Defaults to `true`.
- `should_create_incident` (Boolean) Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per rule. Defaults to `false`.
- `should_write_detection_finding` (Boolean) Whether matches also write a Detection Finding security event back into the events table. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) ID of the user who created this detection rule. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_error` (String) The most recent evaluation error, if any. Cleared on the next successful evaluation.
- `last_evaluated_at` (String) When the detection engine last evaluated this rule. Null means it has never run.
- `last_match_at` (String) When this rule most recently matched security events. Null means it has never matched.
- `project_id` (String) ID of the project this detection rule belongs to. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing detection rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_detection_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_detection_rule.example <id>
```
