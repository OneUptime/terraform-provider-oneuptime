---
page_title: "oneuptime_detection_rule Resource - oneuptime"
subcategory: "Security"
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
- `is_real_time` (Boolean) Whether events are checked against this rule as they arrive, within seconds, instead of every evaluation interval. A rule that counts, and a correlation, count in windows as long as the evaluation interval and fire at most once in each. Not for a rule with a regular expression. Defaults to `false`.
- `match_count_threshold` (Number) Fire only when a group's count — distinct values when a distinct count field is set, matching events otherwise — reaches this number within one evaluation window. 1 fires on any match. Defaults to `1`.
- `response_steps` (String) What to do when the rule fires: up to 10 steps, each { id, phase, text }, where phase is Confirm, Contain or Fix. A case's Response Plan shows them ahead of its general steps. A JSON value: write it with `jsonencode()`.
- `risk_score` (Number) What a match adds to the risk of the user, host or address it is about, from 0 to 100. Empty: by the rule's level.
- `should_create_alert` (Boolean) Whether matches open OneUptime alerts. Defaults to `true`.
- `should_create_case` (Boolean) Whether matches open security cases (or join the open case about the same entity), where the security team works them and playbooks respond. Defaults to `true`.
- `should_create_incident` (Boolean) Whether matches also open OneUptime incidents. Off by default: incidents drive on-call, SLAs and status pages, so opt in per rule. Defaults to `false`.
- `should_write_detection_finding` (Boolean) Whether matches also write a Detection Finding security event back into the events table. Defaults to `true`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) ID of the user who created this detection rule. The ID of a `oneuptime_user` (see the data source).
- `evaluating_until` (String) While the detection engine records what a run found: the end of the events that run judged. Null when no run is recording.
- `id` (String) Unique identifier for the resource.
- `last_error` (String) The most recent evaluation error, if any. Cleared on the next successful evaluation.
- `last_evaluated_at` (String) When the detection engine last evaluated this rule. Null means it has never run.
- `last_failed_at` (String) When the rule's last scheduled run failed. Null when its last run worked.
- `last_match_at` (String) When this rule most recently matched security events. Null means it has never matched.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
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
