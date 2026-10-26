---
page_title: "oneuptime_ai_rule Data Source - oneuptime"
subcategory: "Other"
description: |-
  Which new incidents or alerts OneUptime AI investigates and fixes, and how. A rule investigates what it matches (right away, or only if still open after a wait), and fixes it or not, asking first or not. A fix may run the project's runbooks (only the rule's, when it names some), Bash commands on a Runner the rule picks, the tools of the project's external MCP servers (only the rule's, when it names some), and open fix and telemetry pull requests in the project's code repositories (only the rule's, when it names some). With no rule, OneUptime AI investigates every one, and fixes every one while automatic fixing is on. These were auto remediation rules, and investigation rules are AI rules now too. A rule whose Trigger Entity Type is Profile, Trace, Exception, Log or Metric decides the same about the problems OneUptime finds in a service's profiles, traces, exceptions, logs or metrics, matched by their service; their fix is a pull request in one of the project's code repositories.
---

# oneuptime_ai_rule (Data Source)

Which new incidents or alerts OneUptime AI investigates and fixes, and how. A rule investigates what it matches (right away, or only if still open after a wait), and fixes it or not, asking first or not. A fix may run the project's runbooks (only the rule's, when it names some), Bash commands on a Runner the rule picks, the tools of the project's external MCP servers (only the rule's, when it names some), and open fix and telemetry pull requests in the project's code repositories (only the rule's, when it names some). With no rule, OneUptime AI investigates every one, and fixes every one while automatic fixing is on. These were auto remediation rules, and investigation rules are AI rules now too. A rule whose Trigger Entity Type is Profile, Trace, Exception, Log or Metric decides the same about the problems OneUptime finds in a service's profiles, traces, exceptions, logs or metrics, matched by their service; their fix is a pull request in one of the project's code repositories.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one ai rule may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

~> **Renamed:** this data source was called `oneuptime_auto_remediation_rule` before. The old name still works, but is deprecated.

## Example Usage

```terraform
data "oneuptime_ai_rule" "example" {
  name = "Example ai rule"
}

# Or by id:
data "oneuptime_ai_rule" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_composes_commands` (Boolean) Kept for rules saved before rules were simplified: the AI investigates the incident/alert and composes Bash commands for the rule's Runner and Command Runners, and the rule does nothing else. Suggest mode proposes a command plan for one-click approval; FullAuto mode may execute commands inline, but only ones matching the command allowlist. A OneUptime AI rule now composes commands on its Runner when it picks one, so new rules leave this off.
- `ai_context` (String) Notes OneUptime AI reads after the project's AI context when it investigates or fixes what this rule matches. Markdown, at most 5,000 characters.
- `ai_selects_runbook` (Boolean) Kept for rules saved before rules were simplified: an AI planning run reads the incident/alert context and picks the most applicable runbook (from the attached candidates, or all enabled runbooks when none are attached), always proposing it for approval, and the rule does nothing else. Every OneUptime AI rule now lets AI pick a runbook on its own, so new rules leave this off.
- `auto_resolve_on_verified_recovery` (Boolean) When verification confirms the monitors recovered inside the window, automatically resolve the incident/alert. Off by default — the timeline note is posted either way.
- `created_by_user_id` (String) User ID who created this object. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this AI rule.
- `description_pattern` (String) Case-insensitive regex matched against the entity's description. Leave empty to match any description.
- `execution_mode` (String) How a rule that fixes asks. Suggest asks before fixing: every fix the rule starts waits for one-click human approval. FullAuto fixes without asking: OneUptime AI fixes on its own where the cluster's or resource's AI agent page allows, a runbook it picks from the ones the rule names starts right away (one it picks from all of the project's runbooks still waits for approval), a change through an MCP server the rule names is made right away (through any of the project's MCP servers, it still waits for approval), and a rule that runs its runbooks every time starts them immediately.
- `fixes` (Boolean) Whether OneUptime AI fixes the incidents or alerts this rule matches, the way its Execution Mode says. Off, the rule only investigates them. With enabled rules, an incident is fixed only when a rule that fixes matches it. On by default.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `investigates` (Boolean) Whether OneUptime AI investigates the incidents or alerts this rule matches. With enabled rules that investigate, only their matches are investigated. Off, the rule only fixes, and does not change which incidents are investigated. The dashboard turns it on for every rule; left out of an API write, it is off.
- `is_enabled` (Boolean) Whether this rule is enabled.
- `llm_provider_id` (String) ID of the LLM provider OneUptime AI uses for everything this rule starts. Empty: the project's default provider, else the global one. The ID of a `oneuptime_llm_provider`.
- `name` (String) Name of this AI rule.
- `only_if_open_after_minutes` (Number) Investigate only incidents or alerts still open this many minutes after they were created: OneUptime AI waits that long, then investigates if it is still open. Leave empty to investigate right away. When several rules match, the shortest wait decides. At most 1440 minutes (a day).
- `open_code_fix_pull_requests` (Boolean) When an investigation of an incident or alert this rule matches finds the cause in your code, OneUptime AI opens a pull request with the fix, for review. Nothing is merged. Acts only while fixing new incidents (or alerts) automatically is on. On by default.
- `open_telemetry_pull_requests` (Boolean) When missing logs, traces or metrics stop an investigation of an incident or alert this rule matches, OneUptime AI opens a pull request that adds them, for review. Nothing is merged. Acts only while fixing new incidents (or alerts) automatically is on. On by default.
- `remediation_action` (String) OneUptimeAI (the default): OneUptime AI fixes the matched incident or alert - on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, with the project's runbooks (only the rule's Runbooks when it names some) and with commands on the rule's Runner when it picks one. Runbooks: the rule starts every one of its Runbooks for every match, without OneUptime AI deciding. Whether a person approves first is the rule's Execution Mode.
- `runner_id` (String) ID of the Runner OneUptime AI runs this rule's Bash commands on: one of the project's own Runners. Empty means none: OneUptime AI runs no commands of its own for the rule. Commands on the affected Kubernetes cluster or host always go through its AI agent. The ID of a `oneuptime_runner`.
- `service_name_pattern` (String) Profile, Trace, Exception, Log and Metric AI rules only: case-insensitive regex matched against the name of the service OneUptime found the problem in. Leave empty to match any name.
- `title_pattern` (String) Case-insensitive regex matched against the entity's title. Leave empty to match any title.
- `trigger_entity_type` (String) Which kind of signal this rule decides about: Incident or Alert - or Profile, the slow code OneUptime finds in a service's profiles, or Trace, the slow and failing operations OneUptime finds in a service's traces, or Exception, the new and rising exceptions OneUptime finds in a service, or Log, the error, missing and noisy logs OneUptime finds in a service's logs, or Metric, the problems OneUptime finds in a service's metrics.
- `verification_window_minutes` (Number) How long after the runbook starts the subject's monitors get to recover before verification fails. Defaults to 15 minutes.

### Read-Only

- `alert_severities` (Set of String) Only trigger for alerts with these severities (alert rules only). Leave empty to match any severity. IDs of `oneuptime_alert_severity` resources.
- `code_repositories` (Set of String) The code repositories OneUptime AI may read while it investigates and fixes what this rule matches, and open the rule's pull requests in. Leave empty to let it use every repository the project has connected. IDs of `oneuptime_code_repository` resources.
- `command_allowlist` (String) Glob patterns for commands the AI may execute WITHOUT human approval under FullAuto (for example: systemctl restart *). Commands that do not match are proposed for one-click approval instead. Destructive commands are always refused by the built-in policy. A JSON value: write it with `jsonencode()`.
- `command_runners` (Set of String) Kept for rules saved before rules were simplified: more Runners the AI may target with composed commands, besides the rule's Runner. A Runner runs a rule's commands only when the rule names it. IDs of `oneuptime_runner` resources.
- `created_at` (String) Date and Time when the object was created.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `external_mcp_servers` (Set of String) The external MCP servers OneUptime AI may use while it investigates and fixes what this rule matches (an investigation uses only their tools that read). Leave empty to let it use all of the project's MCP servers; changes through them are then always proposed for approval first. IDs of `oneuptime_external_mcp_server` resources.
- `incident_severities` (Set of String) Only trigger for incidents with these severities (incident rules only). Leave empty to match any severity. IDs of `oneuptime_incident_severity` resources.
- `labels` (Set of String) Only trigger for incidents/alerts that carry at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.
- `monitor_labels` (Set of String) Only trigger when the incident/alert's monitor carries at least one of these labels — the natural way to scope rules to environments (e.g. staging vs production). Leave empty to match any monitor label. IDs of `oneuptime_label` resources.
- `monitors` (Set of String) Only trigger for incidents/alerts from these monitors. Leave empty to match any monitor. IDs of `oneuptime_monitor` resources.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `runbooks` (Set of String) The runbooks OneUptime AI may pick and run for this rule. Leave empty to let it pick any of the project's runbooks. A rule whose Fix With is Runbooks starts every one of them for every match instead. IDs of `oneuptime_runbook` resources.
- `service_labels` (Set of String) Profile, Trace, Exception, Log and Metric AI rules only: match what OneUptime found in a service that carries at least one of these labels. Leave empty to match any label. IDs of `oneuptime_label` resources.
- `services` (Set of String) Profile, Trace, Exception, Log and Metric AI rules only: match what OneUptime found in one of these services. Leave empty to match any service. IDs of `oneuptime_service` resources.
- `updated_at` (String) Date and Time when the object was updated.
