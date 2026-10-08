---
page_title: "oneuptime_auto_remediation_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Which new incidents or alerts are fixed automatically, and how: by OneUptime AI or with runbooks, asking first or not. With no rule, OneUptime AI fixes every one while automatic fixing is on.
---

# oneuptime_auto_remediation_rule (Resource)

Which new incidents or alerts are fixed automatically, and how: by OneUptime AI or with runbooks, asking first or not. With no rule, OneUptime AI fixes every one while automatic fixing is on.

## Example Usage

```terraform
resource "oneuptime_auto_remediation_rule" "example" {
  name                = "Example auto remediation rule"
  trigger_entity_type = "Example short text"
  description         = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this auto-remediation rule.
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident or Alert.

### Optional

- `ai_composes_commands` (Boolean) When enabled, the AI investigates the incident/alert and composes Bash/SSH commands for opted-in Runners instead of picking a runbook. Suggest mode proposes a command plan for one-click approval; FullAuto mode may execute commands inline, but only ones matching the command allowlist. Requires AI to be enabled for the project. Defaults to `false`.
- `ai_selects_runbook` (Boolean) When enabled, an AI planning run reads the incident/alert context and picks the most applicable runbook (from the attached candidates, or all enabled runbooks when none are attached). AI-picked runbooks are always suggest-only — never full-auto. Defaults to `false`.
- `alert_severities` (Set of String) Only trigger for alerts with these severities (alert rules only). Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.
- `auto_resolve_on_verified_recovery` (Boolean) When verification confirms the monitors recovered inside the window, automatically resolve the incident/alert. Off by default — the timeline note is posted either way. Defaults to `false`.
- `command_allowlist` (String) Glob patterns for commands the AI may execute WITHOUT human approval under FullAuto (for example: systemctl restart *). Commands that do not match are proposed for one-click approval instead. Destructive commands are always refused by the built-in policy. A JSON value: write it with `jsonencode()`.
- `command_runners` (Set of String) Runners the AI may target with composed commands. Leave empty to allow any Runner in the project that has AI commands enabled. IDs of `oneuptime_runner` resources.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of this auto-remediation rule.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description.
- `execution_mode` (String) Suggest asks before fixing: every fix the rule starts waits for one-click human approval. FullAuto fixes without asking: its runbooks start immediately, and OneUptime AI fixes run on their own where the cluster's or resource's AI agent page allows. Defaults to `Suggest`.
- `incident_severities` (Set of String) Only trigger for incidents with these severities (incident rules only). Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `labels` (Set of String) Only trigger for incidents/alerts that carry at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.
- `monitor_labels` (Set of String) Only trigger when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only trigger for incidents/alerts from these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.
- `remediation_action` (String) OneUptimeAI: OneUptime AI fixes the matched incident or alert on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows. Runbooks: the rule's runbooks run. Whether a person approves first is the rule's Execution Mode. Defaults to `OneUptimeAI`.
- `runbooks` (Set of String) Runbook candidates for this rule. Deterministic rules propose or start every attached runbook; AI rules pick the most applicable one. IDs of `oneuptime_runbook` resources.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title.
- `verification_window_minutes` (Number) How long after the runbook starts the subject's monitors get to recover before verification fails. Defaults to 15 minutes.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object. The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing auto remediation rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_auto_remediation_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_auto_remediation_rule.example <id>
```
