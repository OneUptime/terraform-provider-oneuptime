---
page_title: "oneuptime_project Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  OneUptime Project, and everything happens inside it
---

# oneuptime_project (Data Source)

OneUptime Project, and everything happens inside it Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_project" "by_name" {
  name = "example-project"
}

data "oneuptime_project" "by_id" {
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
- `slug` (String) Friendly globally unique name for your object.. Computed.
- `payment_provider_plan_id` (String) Project payment_provider_plan_id. Computed.
- `payment_provider_subscription_id` (String) Project payment_provider_subscription_id. Computed.
- `payment_provider_metered_subscription_id` (String) Project payment_provider_metered_subscription_id. Computed.
- `payment_provider_subscription_seats` (Number) Project payment_provider_subscription_seats. Computed.
- `trial_ends_at` (String) A date time object.. Computed.
- `payment_provider_customer_id` (String) Project payment_provider_customer_id. Computed.
- `business_details` (String) Business legal name, address and any tax information to appear on invoices... Computed.
- `business_details_country` (String) Two-letter ISO country code for billing address (e.g., US, GB, DE)... Computed.
- `finance_accounting_email` (String) Invoices, receipts and billing related notifications will be sent to these emails in addition to project owner. Separate multiple emails with a comma... Computed.
- `payment_provider_subscription_status` (String) Project payment_provider_subscription_status. Computed.
- `payment_provider_metered_subscription_status` (String) Project payment_provider_metered_subscription_status. Computed.
- `payment_provider_promo_code` (String) Project payment_provider_promo_code. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_feature_flag_monitor_groups_enabled` (Bool) Is Feature Flag Monitor Groups Enabled.. Computed.
- `workflow_runs_in_last30_days` (Number) Project workflow_runs_in_last30_days. Computed.
- `require_sso_for_login` (Bool) Project require_sso_for_login. Computed.
- `require_sso_with_sso_provider_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_number_prefix` (String) Custom prefix for incident numbers (e.g., 'INC-'). If empty, '#' is used... Computed.
- `alert_number_prefix` (String) Custom prefix for alert numbers (e.g., 'ALT-'). If empty, '#' is used... Computed.
- `scheduled_maintenance_number_prefix` (String) Custom prefix for scheduled maintenance numbers (e.g., 'SM-'). If empty, '#' is used... Computed.
- `incident_episode_number_prefix` (String) Custom prefix for incident episode numbers (e.g., 'IE-'). If empty, '#' is used... Computed.
- `alert_episode_number_prefix` (String) Custom prefix for alert episode numbers (e.g., 'AE-'). If empty, '#' is used... Computed.
- `sms_or_call_current_balance_in_usd_cents` (Number) Balance in USD for SMS, Call, and WhatsApp.. Computed.
- `auto_recharge_sms_or_call_by_balance_in_usd` (Number) Auto recharge amount in USD for SMS, Call, and WhatsApp.. Computed.
- `auto_recharge_sms_or_call_when_current_balance_falls_in_usd` (Number) Auto recharge is triggered when current balance falls to this amount in USD for SMS, Call, and WhatsApp.. Computed.
- `enable_sms_notifications` (Bool) Enable SMS notifications for this project... Computed.
- `enable_whats_app_notifications` (Bool) Enable WhatsApp notifications for this project... Computed.
- `enable_telegram_notifications` (Bool) Enable Telegram notifications for this project... Computed.
- `enable_call_notifications` (Bool) Enable call notifications for this project... Computed.
- `disable_on_call_notification_fallback` (Bool) When enabled, a page routed to a responder with no matching notification rule fails instead of falling back to their verified notification methods... Computed.
- `enable_auto_recharge_sms_or_call_balance` (Bool) Enable auto recharge for SMS, Call, and WhatsApp balance for this project... Computed.
- `ai_current_balance_in_usd_cents` (Number) Balance in USD for AI services.. Computed.
- `auto_ai_recharge_by_balance_in_usd` (Number) Auto recharge amount in USD for AI services.. Computed.
- `auto_recharge_ai_when_current_balance_falls_in_usd` (Number) Auto recharge is triggered when current balance falls to this amount in USD for AI services.. Computed.
- `enable_ai` (Bool) Enable AI services for this project... Computed.
- `enable_auto_remediation` (Bool) Kill switch for auto-remediation: when disabled, no auto-remediation rule fires in this project... Computed.
- `enable_ai_command_execution` (Bool) When enabled, auto-remediation rules may let the AI compose and run commands on opted-in Runners (with an operator allowlist for auto-execution, and one-click approval for everything else). Off by default... Computed.
- `enable_automatic_incident_investigation` (Bool) When enabled, OneUptime's AI SRE automatically investigates every new incident and posts a cited root cause analysis to the incident timeline. Requires AI to be enabled and an LLM provider to be configured... Computed.
- `enable_automatic_alert_investigation` (Bool) When enabled, OneUptime's AI SRE automatically investigates every new alert and posts a cited root cause analysis to the alert timeline. Requires AI to be enabled and an LLM provider to be configured... Computed.
- `enable_incident_instrumentation_fix_tasks` (Bool) When enabled, an incident AI investigation that ends inconclusive (telemetry was insufficient to determine a root cause) automatically queues an AI agent task that opens a pull request adding the missing instrumentation to the implicated code paths. Requires a repository connected through the GitHub App. Pull requests are always human-reviewed — nothing merges automatically... Computed.
- `enable_alert_instrumentation_fix_tasks` (Bool) When enabled, an alert AI investigation that ends inconclusive (telemetry was insufficient to determine a root cause) automatically queues an AI agent task that opens a pull request adding the missing instrumentation to the implicated code paths. Requires a repository connected through the GitHub App. Pull requests are always human-reviewed — nothing merges automatically... Computed.
- `enable_automatic_incident_code_fixes` (Bool) When enabled, an incident AI investigation that ends with a confident, evidenced root cause analysis and recommends a repository code change automatically queues an AI agent task that opens a fix pull request, ready for review, from that analysis — the automatic form of the 'Open Fix PR from this analysis' button. Operational, infrastructure, external, user-error and inconclusive findings do not offer or open code-fix pull requests. Requires a repository connected through the GitHub App and a Runner with the code-fix capability. Pull requests are always human-reviewed — nothing merges automatically... Computed.
- `enable_automatic_alert_code_fixes` (Bool) When enabled, an alert AI investigation that ends with a confident, evidenced root cause analysis and recommends a repository code change automatically queues an AI agent task that opens a fix pull request, ready for review, from that analysis — the automatic form of the 'Open Fix PR from this analysis' button. Operational, infrastructure, external, user-error and inconclusive findings do not offer or open code-fix pull requests. Requires a repository connected through the GitHub App and a Runner with the code-fix capability. Pull requests are always human-reviewed — nothing merges automatically... Computed.
- `enable_ai_insights` (Bool) When enabled, OneUptime AI continuously watches this project's telemetry with deterministic statistical sensors (error-log spikes, exception novelty and spikes, trace-latency regressions, week-over-week metric drift) and files quiet Insights — never pages, never opens incidents. Each new insight also gets a budgeted, read-only AI triage analysis when an LLM provider is configured... Computed.
- `enable_insight_fix_tasks` (Bool) When enabled, insights whose deterministic evidence points at code (new or spiking exceptions with a resolvable repository, trace-latency regressions with span-tree findings) automatically queue an AI agent task that opens a pull request with a proposed fix, ready for review. Honors the daily fix task budget and per-repository open-PR caps. Pull requests are always human-reviewed — nothing merges automatically... Computed.
- `auto_archive_non_actionable_exceptions` (Bool) When enabled, exception groups the AI triage classifies as expected denials (auth failures, plan/paywall rejections, scanner probes tripping intentional validation) are automatically archived so they stop surfacing in the unresolved list and never queue AI fix tasks. Groups classified as user errors or infrastructure conditions are NOT auto-archived — only clear expected denials are. Archiving is reversible from the Archived tab... Computed.
- `alert_investigation_minimum_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `ai_daily_autonomous_token_limit` (Number) Fallback maximum tokens per UTC day for autonomous AI work that is not associated with an incident or alert. When the limit is reached, new autonomous work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit... Computed.
- `incident_ai_daily_autonomous_token_limit` (Number) Maximum tokens per UTC day that autonomous incident-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new incident-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit... Computed.
- `alert_ai_daily_autonomous_token_limit` (Number) Maximum tokens per UTC day that autonomous alert-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new alert-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit... Computed.
- `ai_daily_fix_task_limit` (Number) Fallback maximum AI fix tasks (agent runs that open pull requests) that may be created per UTC day for work not associated with an incident or alert, across every fix recipe and trigger. Unset means the default of 25 per day; 0 pauses these AI fix tasks entirely... Computed.
- `incident_ai_daily_fix_task_limit` (Number) Maximum AI fix tasks derived from incidents that may be created per UTC day for this project. Unset means the default of 25 per day; 0 pauses incident AI fix tasks entirely... Computed.
- `alert_ai_daily_fix_task_limit` (Number) Maximum AI fix tasks derived from alerts that may be created per UTC day for this project. Unset means the default of 25 per day; 0 pauses alert AI fix tasks entirely... Computed.
- `alert_investigation_dedupe_window_minutes` (Number) Repeat alerts from the same monitor within this many minutes are not re-investigated by AI — the first analysis stands. Unset means the default of 30 minutes; 0 disables the cooldown... Computed.
- `incident_investigation_minimum_severity_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `incident_investigation_dedupe_window_minutes` (Number) Incidents affecting a monitor that AI investigated within this many minutes are not re-investigated — the first analysis stands. Unset means the default of 30 minutes; 0 disables the cooldown... Computed.
- `ai_max_concurrent_investigations` (Number) Fallback maximum number of non-incident and non-alert AI investigations that may run at the same time for this project. Unset means the default of 3. Minimum 1 — pause autonomous work with its opt-in toggle or a daily token limit of 0 instead... Computed.
- `incident_ai_max_concurrent_investigations` (Number) How many incident AI investigations may run at the same time for this project. Unset means the default of 3. Minimum 1 — pause incident investigations with the opt-in toggle or a daily token limit of 0 instead... Computed.
- `alert_ai_max_concurrent_investigations` (Number) How many alert AI investigations may run at the same time for this project. Unset means the default of 3. Minimum 1 — pause alert investigations with the opt-in toggle or a daily token limit of 0 instead... Computed.
- `enable_auto_recharge_ai_balance` (Bool) Enable auto recharge for AI balance for this project... Computed.
- `send_invoices_by_email` (Bool) When enabled, invoices will be automatically sent to the finance/accounting email when they are generated... Computed.
- `plan_name` (String) Name of the plan this project is subscribed to... Computed.
- `reseller_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `reseller_plan_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `let_customer_support_access_project` (Bool) OneUptime customer support can access this project. This is used for debugging purposes... Computed.
- `do_not_add_global_probes_by_default_on_new_monitors` (Bool) If enabled, global probes will NOT be automatically added to new monitors. Enable this only if you are using ONLY custom probes to monitor your resources... Computed.
- `git_hub_app_installation_id` (String) The GitHub App installation ID for this project. This is set when the GitHub App is installed on the organization... Computed.
- `default_metric_cardinality_budget` (Number) Project-wide default max distinct series per metric. Services without a per-service override use this value... Computed.
- `default_telemetry_retention_in_days` (Number) Project-wide default number of days to retain telemetry data (logs, traces, metrics). Services without a per-service override use this value... Computed.
- `telemetry_retention_config` (String) Project-wide per-pillar retention overrides for telemetry data (logs by severity, traces by status, metrics, profiles). Falls back to defaultTelemetryRetentionInDays when a pillar or bucket is not set... Computed.
- `default_metric_downsampling_retention_days` (String) Project-wide default retention for each downsampling tier (raw, 1m, 5m, 1h, 1d) in days... Computed.
- `enable_audit_logs` (Bool) When enabled, changes to resources in this project are recorded as audit log entries... Computed.
- `is_session_replay_allowed` (Bool) When enabled, RUM applications in this project may record session replays if they are individually enabled too. On by default; switch it off here to stop session replay across the entire project in one place... Computed.
- `audit_logs_retention_in_days` (Number) Number of days to retain audit log entries. Minimum 7, maximum 180... Computed.
- `store_system_events_in_audit_logs` (Bool) When enabled, audit logs will also include events triggered by the system. By default, only events triggered by users are recorded... Computed.
