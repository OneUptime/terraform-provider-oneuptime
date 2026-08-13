---
page_title: "oneuptime_auto_remediation_suggestion Data Source - oneuptime"
subcategory: "Other"
description: |-
  A proposed or executed remediation runbook attached to an incident or alert by an auto-remediation rule.
---

# oneuptime_auto_remediation_suggestion (Data Source)

A proposed or executed remediation runbook attached to an incident or alert by an auto-remediation rule. Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_auto_remediation_suggestion" "by_name" {
  name = "example-auto_remediation_suggestion"
}

data "oneuptime_auto_remediation_suggestion" "by_id" {
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
- `auto_remediation_rule_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `rule_name_snapshot` (String) Name of the rule when this suggestion was created — survives rule deletion... Computed.
- `incident_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `alert_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `runbook_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `runbook_name_snapshot` (String) Name of the proposed runbook when this suggestion was created — survives runbook deletion... Computed.
- `status` (String) Lifecycle status: Planning, Suggested, Approved, AutoExecuted, Dismissed or NoneApplicable... Computed.
- `execution_mode` (String) The rule's execution mode when this suggestion was created (Suggest or FullAuto)... Computed.
- `suggestion_type` (String) Runbook suggestions propose starting a pre-authored runbook; CommandPlan suggestions carry an AI-composed command plan... Computed.
- `command_plan` (String) The AI-composed command plan for CommandPlan suggestions, including per-command execution results once run... Computed.
- `rationale_markdown` (String) Why this runbook was proposed — the AI planning run's reasoning for AI rules, or a short note for deterministic rules... Computed.
- `ai_run_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `runbook_execution_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `approved_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `approved_at` (String) A date time object.. Computed.
- `dismissed_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `dismissed_at` (String) A date time object.. Computed.
- `verification_status` (String) Outcome verification after execution: Pending, Verified, Failed or Skipped. Empty until a runbook is started... Computed.
- `verification_deadline_at` (String) A date time object.. Computed.
- `verification_completed_at` (String) A date time object.. Computed.
- `verification_note` (String) Why verification ended the way it did... Computed.
- `verification_window_minutes` (Number) Snapshot of the rule's verification window when this suggestion was created... Computed.
- `auto_resolve_on_recovery` (Bool) Snapshot of the rule's auto-resolve-on-verified-recovery setting when this suggestion was created... Computed.
