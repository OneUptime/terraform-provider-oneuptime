---
page_title: "oneuptime_auto_remediation_rule Resource - oneuptime"
subcategory: "Other"
description: |-
  Automatically propose or start remediation runbooks when matching incidents or alerts are created.
---

# oneuptime_auto_remediation_rule (Resource)

Automatically propose or start remediation runbooks when matching incidents or alerts are created.

## Example Usage

```terraform
resource "oneuptime_auto_remediation_rule" "example" {
  name = "Example short text"
  trigger_entity_type = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this auto-remediation rule...
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident or Alert...

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this auto-remediation rule...
- `is_enabled` (Bool) Whether this rule is enabled...
- `execution_mode` (String) Suggest proposes the runbook and waits for one-click human approval. FullAuto starts it immediately (deterministic rules only)...
- `ai_selects_runbook` (Bool) When enabled, an AI planning run reads the incident/alert context and picks the most applicable runbook (from the attached candidates, or all enabled runbooks when none are attached). AI-picked runbooks are always suggest-only — never full-auto...
- `ai_composes_commands` (Bool) When enabled, the AI investigates the incident/alert and composes Bash/SSH commands for opted-in Runners instead of picking a runbook. Suggest mode proposes a command plan for one-click approval; FullAuto mode may execute commands inline, but only ones matching the command allowlist. Requires the project's Enable AI Command Execution setting...
- `command_allowlist` (String) Glob patterns for commands the AI may execute WITHOUT human approval under FullAuto (for example: systemctl restart *). Commands that do not match are proposed for one-click approval instead. Destructive commands are always refused by the built-in policy...
- `command_runners` (Set) Runners the AI may target with composed commands. Leave empty to allow any Runner in the project that has AI commands enabled...
- `monitors` (Set) Only trigger for incidents/alerts from these monitors. Leave empty to match any monitor...
- `incident_severities` (Set) Only trigger for incidents with these severities (incident rules only). Leave empty to match any severity...
- `alert_severities` (Set) Only trigger for alerts with these severities (alert rules only). Leave empty to match any severity...
- `labels` (Set) Only trigger for incidents/alerts that carry at least one of these labels. Leave empty to match any label...
- `monitor_labels` (Set) Only trigger when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label...
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title...
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description...
- `runbooks` (Set) Runbook candidates for this rule. Deterministic rules propose or start every attached runbook; AI rules pick the most applicable one...
- `verification_window_minutes` (Number) How long after the runbook starts the subject's monitors get to recover before verification fails. Defaults to 15 minutes...
- `auto_resolve_on_verified_recovery` (Bool) When verification confirms the monitors recovered inside the window, automatically resolve the incident/alert. Off by default — the timeline note is posted either way...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_auto_remediation_rule.example <id>
```
