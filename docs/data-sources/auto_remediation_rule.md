---
page_title: "oneuptime_auto_remediation_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Which new incidents or alerts OneUptime AI fixes automatically, and whether it asks first. It may run the project's runbooks (only the rule's, when it names some), commands on a Runner the rule picks, and the tools of the project's external MCP servers (only the rule's, when it names some). With no rule, OneUptime AI fixes every one while automatic fixing is on.
---

# oneuptime_auto_remediation_rule (Data Source)

Which new incidents or alerts OneUptime AI fixes automatically, and whether it asks first. It may run the project's runbooks (only the rule's, when it names some), commands on a Runner the rule picks, and the tools of the project's external MCP servers (only the rule's, when it names some). With no rule, OneUptime AI fixes every one while automatic fixing is on.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one auto remediation rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_auto_remediation_rule" "example" {
  name = "Example auto remediation rule"
}

# Or by id:
data "oneuptime_auto_remediation_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_composes_commands` (Boolean) Kept for rules saved before rules were simplified: the AI investigates the incident/alert and composes Bash/SSH commands for the rule's Runner and Command Runners, and the rule does nothing else. Suggest mode proposes a command plan for one-click approval; FullAuto mode may execute commands inline, but only ones matching the command allowlist. A OneUptime AI rule now composes commands on its Runner when it picks one, so new rules leave this off.
- `ai_selects_runbook` (Boolean) Kept for rules saved before rules were simplified: an AI planning run reads the incident/alert context and picks the most applicable runbook (from the attached candidates, or all enabled runbooks when none are attached), always proposing it for approval, and the rule does nothing else. Every OneUptime AI rule now lets AI pick a runbook on its own, so new rules leave this off.
- `auto_resolve_on_verified_recovery` (Boolean) When verification confirms the monitors recovered inside the window, automatically resolve the incident/alert. Off by default — the timeline note is posted either way.
- `created_by_user_id` (String) User ID who created this object. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this auto-remediation rule.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description.
- `execution_mode` (String) Suggest asks before fixing: every fix the rule starts waits for one-click human approval. FullAuto fixes without asking: OneUptime AI fixes on its own where the cluster's or resource's AI agent page allows, a runbook it picks from the ones the rule names starts right away (one it picks from all of the project's runbooks still waits for approval), a change through an MCP server the rule names is made right away (through any of the project's MCP servers, it still waits for approval), and a rule that runs its runbooks every time starts them immediately.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `name` (String) Name of this auto-remediation rule.
- `remediation_action` (String) OneUptimeAI (the default): OneUptime AI fixes the matched incident or alert - on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, with the project's runbooks (only the rule's Runbooks when it names some) and with commands on the rule's Runner when it picks one. Runbooks: the rule starts every one of its Runbooks for every match, without OneUptime AI deciding. Whether a person approves first is the rule's Execution Mode.
- `runner_id` (String) ID of the Runner OneUptime AI runs this rule's commands on: one of the project's own Runners. Empty means none: OneUptime AI fixes through the AI agent on the affected Kubernetes cluster or host. The ID of a `oneuptime_runner`.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title.
- `trigger_entity_type` (String) Entity type that triggers this rule on creation: Incident or Alert.
- `verification_window_minutes` (Number) How long after the runbook starts the subject's monitors get to recover before verification fails. Defaults to 15 minutes.

### Read-Only

- `alert_severities` (Set of String) Only trigger for alerts with these severities (alert rules only). Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.
- `command_allowlist` (String) Glob patterns for commands the AI may execute WITHOUT human approval under FullAuto (for example: systemctl restart *). Commands that do not match are proposed for one-click approval instead. Destructive commands are always refused by the built-in policy. A JSON value: write it with `jsonencode()`.
- `command_runners` (Set of String) Kept for rules saved before rules were simplified: more Runners the AI may target with composed commands, besides the rule's Runner. A Runner runs a rule's commands only when the rule names it. IDs of `oneuptime_runner` resources.
- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `external_mcp_servers` (Set of String) The external MCP servers OneUptime AI may use while it fixes what this rule matches. Leave empty to let it use all of the project's MCP servers; changes through them are then always proposed for approval first. IDs of `oneuptime_external_mcp_server` resources.
- `incident_severities` (Set of String) Only trigger for incidents with these severities (incident rules only). Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.
- `labels` (Set of String) Only trigger for incidents/alerts that carry at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.
- `monitor_labels` (Set of String) Only trigger when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only trigger for incidents/alerts from these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `runbooks` (Set of String) The runbooks OneUptime AI may pick and run for this rule. Leave empty to let it pick any of the project's runbooks. A rule whose Fix With is Runbooks starts every one of them for every match instead. IDs of `oneuptime_runbook` resources.
- `updated_at` (String) Date and Time when the object was updated.
