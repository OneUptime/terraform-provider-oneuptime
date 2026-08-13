---
page_title: "oneuptime_auto_remediation_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Automatically propose or start remediation runbooks when matching incidents or alerts are created.
---

# oneuptime_auto_remediation_rule (Data Source)

Automatically propose or start remediation runbooks when matching incidents or alerts are created. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_auto_remediation_rule" "by_name" {
  name = "example-auto_remediation_rule"
}

data "oneuptime_auto_remediation_rule" "by_id" {
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
- `description` (String) Description of this auto-remediation rule... Computed.
- `is_enabled` (Bool) Whether this rule is enabled... Computed.
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident or Alert... Computed.
- `execution_mode` (String) Suggest proposes the runbook and waits for one-click human approval. FullAuto starts it immediately (deterministic rules only)... Computed.
- `ai_selects_runbook` (Bool) When enabled, an AI planning run reads the incident/alert context and picks the most applicable runbook (from the attached candidates, or all enabled runbooks when none are attached). AI-picked runbooks are always suggest-only — never full-auto... Computed.
- `ai_composes_commands` (Bool) When enabled, the AI investigates the incident/alert and composes Bash/SSH commands for opted-in Runners instead of picking a runbook. Suggest mode proposes a command plan for one-click approval; FullAuto mode may execute commands inline, but only ones matching the command allowlist. Requires the project's Enable AI Command Execution setting... Computed.
- `command_allowlist` (String) Glob patterns for commands the AI may execute WITHOUT human approval under FullAuto (for example: systemctl restart *). Commands that do not match are proposed for one-click approval instead. Destructive commands are always refused by the built-in policy... Computed.
- `command_runners` (Set) Runners the AI may target with composed commands. Leave empty to allow any Runner in the project that has AI commands enabled... Computed.
- `monitors` (Set) Only trigger for incidents/alerts from these monitors. Leave empty to match any monitor... Computed.
- `incident_severities` (Set) Only trigger for incidents with these severities (incident rules only). Leave empty to match any severity... Computed.
- `alert_severities` (Set) Only trigger for alerts with these severities (alert rules only). Leave empty to match any severity... Computed.
- `labels` (Set) Only trigger for incidents/alerts that carry at least one of these labels. Leave empty to match any label... Computed.
- `monitor_labels` (Set) Only trigger when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label... Computed.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title... Computed.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description... Computed.
- `runbooks` (Set) Runbook candidates for this rule. Deterministic rules propose or start every attached runbook; AI rules pick the most applicable one... Computed.
- `verification_window_minutes` (Number) How long after the runbook starts the subject's monitors get to recover before verification fails. Defaults to 15 minutes... Computed.
- `auto_resolve_on_verified_recovery` (Bool) When verification confirms the monitors recovered inside the window, automatically resolve the incident/alert. Off by default — the timeline note is posted either way... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
