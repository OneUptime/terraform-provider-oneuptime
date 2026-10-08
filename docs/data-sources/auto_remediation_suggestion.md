---
page_title: "oneuptime_auto_remediation_suggestion Data Source - oneuptime"
subcategory: "Other"
description: |-
  A proposed or executed remediation runbook attached to an incident or alert by an auto-remediation rule.
---

# oneuptime_auto_remediation_suggestion (Data Source)

A proposed or executed remediation runbook attached to an incident or alert by an auto-remediation rule.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one auto remediation suggestion may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_auto_remediation_suggestion" "example" {
  auto_remediation_rule_id = oneuptime_auto_remediation_rule.example.id
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
- `auto_remediation_rule_id` (String) ID of the rule that produced this suggestion. The ID of a `oneuptime_auto_remediation_rule`.
- `auto_resolve_on_recovery` (Boolean) Snapshot of the rule's auto-resolve-on-verified-recovery setting when this suggestion was created.
- `dismissed_by_user_id` (String) ID of the user who dismissed this suggestion. The ID of a `oneuptime_user` (see the data source).
- `execution_mode` (String) The rule's execution mode when this suggestion was created (Suggest or FullAuto).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this suggestion remediates. The ID of a `oneuptime_incident`.
- `kubernetes_cluster_id` (String) ID of the cluster whose AI remediation mode produced this suggestion. The ID of a `oneuptime_kubernetes_cluster`.
- `rationale_markdown` (String) Why this runbook was proposed — the AI planning run's reasoning for AI rules, or a short note for deterministic rules.
- `resource_id` (String) ID of the resource whose AI remediation mode produced this suggestion, in the table its resource type names.
- `resource_type` (String) The kind of resource whose AI remediation mode produced this suggestion (DockerHost, PodmanHost, DockerSwarmCluster, ProxmoxCluster, VMwareVCenter, CephCluster, DatabaseServer or Host; resource-level remediation, no rule).
- `rule_name_snapshot` (String) Name of the rule when this suggestion was created — survives rule deletion.
- `runbook_execution_id` (String) The runbook execution started when this suggestion was approved or auto-executed.
- `runbook_id` (String) ID of the proposed runbook. The ID of a `oneuptime_runbook`.
- `runbook_name_snapshot` (String) Name of the proposed runbook when this suggestion was created — survives runbook deletion.
- `status` (String) Lifecycle status: Planning, Suggested, Approved, AutoExecuted, Dismissed or NoneApplicable.
- `suggestion_type` (String) Runbook suggestions propose starting a pre-authored runbook; CommandPlan suggestions carry an AI-composed command plan.
- `verification_note` (String) Why verification ended the way it did.
- `verification_status` (String) Outcome verification after execution: Pending, Verified, Failed or Skipped. Empty until a runbook is started.
- `verification_window_minutes` (Number) Snapshot of the rule's verification window when this suggestion was created.

### Read-Only

- `approved_at` (String) When this suggestion was approved.
- `command_plan` (String) The AI-composed command plan for CommandPlan suggestions, including per-command execution results once run. A JSON value: write it with `jsonencode()`.
- `created_at` (String) Date and Time when the object was created.
- `dismissed_at` (String) When this suggestion was dismissed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `verification_completed_at` (String) When verification reached a terminal outcome.
- `verification_deadline_at` (String) When the verification window closes — the monitors must be operational by this time.
