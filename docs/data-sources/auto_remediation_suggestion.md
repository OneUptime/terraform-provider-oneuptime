---
page_title: "oneuptime_auto_remediation_suggestion Data Source - oneuptime"
subcategory: "AI"
description: |-
  A proposed or executed remediation runbook attached to an incident or alert by an auto-remediation rule.
---

# oneuptime_auto_remediation_suggestion (Data Source)

A proposed or executed remediation runbook attached to an incident or alert by an auto-remediation rule.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one auto remediation suggestion may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_auto_remediation_suggestion" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_auto_remediation_suggestion" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `ai_run_id` (String) The AI planning run that picked the runbook (AI rules only).
- `alert_id` (String) ID of the alert this suggestion remediates. The ID of a `oneuptime_alert`.
- `approved_by_user_id` (String) ID of the user who approved this suggestion. The ID of a `oneuptime_user` (see the data source).
- `approved_from` (String) Where this suggestion was approved: Dashboard (the AI card), Slack or MicrosoftTeams. Empty until approved.
- `auto_remediation_rule_id` (String) ID of the AI rule that produced this suggestion. The ID of a `oneuptime_ai_rule`.
- `auto_resolve_on_recovery` (Boolean) Snapshot of the rule's auto-resolve-on-verified-recovery setting when this suggestion was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `dismissed_by_user_id` (String) ID of the user who dismissed this suggestion. The ID of a `oneuptime_user` (see the data source).
- `execution_mode` (String) The rule's execution mode when this suggestion was created (Suggest or FullAuto).
- `first_approved_by_user_id` (String) ID of the user who gave the first of the two approvals a riskier fix needs. The ID of a `oneuptime_user` (see the data source).
- `first_approved_from` (String) Where the first of two approvals was given: Dashboard (the AI card), Slack or MicrosoftTeams.
- `first_approved_plan_fingerprint` (String) The fingerprint of what the first of two approvals approved: a second approval runs the fix only while it still runs exactly that.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this suggestion remediates. The ID of a `oneuptime_incident`.
- `kubernetes_cluster_id` (String) ID of the cluster whose AI remediation mode produced this suggestion. The ID of a `oneuptime_kubernetes_cluster`.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `rationale_markdown` (String) Why this runbook was proposed — the AI planning run's reasoning for AI rules, or a short note for deterministic rules.
- `resource_id` (String) ID of the resource whose AI remediation mode produced this suggestion, in the table its resource type names.
- `resource_type` (String) The kind of resource whose AI remediation mode produced this suggestion (DockerHost, PodmanHost, DockerSwarmCluster, ProxmoxCluster, VMwareVCenter, CephCluster, DatabaseServer or Host; resource-level remediation, no rule).
- `rule_name_snapshot` (String) Name of the rule when this suggestion was created — survives rule deletion.
- `runbook_execution_id` (String) The runbook execution started when this suggestion was approved or auto-executed.
- `runbook_id` (String) ID of the proposed runbook. The ID of a `oneuptime_runbook`.
- `runbook_name_snapshot` (String) Name of the proposed runbook when this suggestion was created — survives runbook deletion.
- `saved_as_runbook_id` (String) ID of the runbook this fix was saved as. Empty until someone saves the verified fix as a runbook. The ID of a `oneuptime_runbook`.
- `status` (String) Lifecycle status: Planning, Suggested, Approved, AutoExecuted, Dismissed or NoneApplicable.
- `suggestion_type` (String) Runbook suggestions propose starting a pre-authored runbook; CommandPlan suggestions carry an AI-composed command plan; McpToolCalls suggestions carry the tool calls OneUptime AI makes, or proposes, through the project's external MCP servers.
- `verification_note` (String) Why verification ended the way it did.
- `verification_status` (String) Outcome verification after execution: Pending, Verified, Failed or Skipped. Empty until a runbook is started.
- `verification_window_minutes` (Number) Snapshot of the rule's verification window when this suggestion was created.

### Read-Only

- `approved_at` (String) When this suggestion was approved.
- `command_plan` (String) The AI-composed command plan for CommandPlan suggestions, including per-command execution results once run. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `dismissed_at` (String) When this suggestion was dismissed.
- `first_approved_at` (String) When the first of the two approvals a riskier fix needs was given.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `mcp_tool_call_plan` (String) For McpToolCalls suggestions: the tool calls OneUptime AI made, or proposes, through the project's external MCP servers, with what each server answered once called. A JSON value: read it with `jsondecode()`.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `runbook_track_record` (String) How the proposed runbook went the last times OneUptime AI ran it for other incidents and alerts, when it was proposed. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
- `verification_completed_at` (String) When verification reached a terminal outcome.
- `verification_deadline_at` (String) When the verification window closes — the monitors must be operational by this time.
