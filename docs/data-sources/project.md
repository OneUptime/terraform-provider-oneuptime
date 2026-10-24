---
page_title: "oneuptime_project Data Source - oneuptime"
subcategory: "Teams & Access"
description: |-
  OneUptime Project, and everything happens inside it
---

# oneuptime_project (Data Source)

OneUptime Project, and everything happens inside it

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one project may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_project" "example" {
  name = "Example project"
}

# Or by id:
data "oneuptime_project" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `acknowledge_linked_alerts_when_incident_acknowledged` (Boolean) When enabled, acknowledging an incident also acknowledges every alert linked to it. This stops those alerts' on-call escalations, and their reminders only when the alert reminder rule is set to stop on Acknowledged. Alerts linked to an incident that is already acknowledged are acknowledged as they are linked. On for new projects created in OneUptime; projects that existed before keep their setting.
- `ai_current_balance_in_usd_cents` (Number) Balance in USD for AI services.
- `ai_daily_autonomous_token_limit` (Number) Legacy setting, no longer enforced: autonomous AI work that is not associated with an incident or alert has no daily token limit. Use the Daily Incident AI Token Limit and Daily Alert AI Token Limit instead.
- `ai_daily_fix_task_limit` (Number) Legacy setting, no longer enforced: AI fix tasks that are not associated with an incident or alert have no daily limit. Use the Daily Incident AI Fix Task Limit and Daily Alert AI Fix Task Limit instead.
- `ai_daily_spend_limit_in_usd` (Number) OneUptime Cloud: the most AI credits, in whole US dollars, OneUptime AI may spend in this project each UTC day. Only calls billed to the project's AI credits count, so it never stops AI that runs on the project's own LLM provider. Once it is reached, billed AI work is refused until midnight UTC. Ignored where AI is not billed (self-hosted). Unset means no limit; a limit is at least 1 (to turn AI off, use Enable AI).
- `ai_daily_token_limit` (Number) The most tokens OneUptime AI may use in this project each UTC day, across every AI feature: Ask AI, investigations, postmortem drafts, fix pull requests, insight triage, workflows, runbooks and Slack or Microsoft Teams questions. Once it is reached, new AI work is refused until midnight UTC. The incident and alert daily limits still apply under it. Unset means no limit; a limit is a whole number of at least 1 (to turn AI off, use Enable AI).
- `ai_max_concurrent_investigations` (Number) Legacy setting, no longer enforced: AI investigations that are not associated with an incident or alert have no concurrency limit. Use the Max Concurrent Incident Investigations and Max Concurrent Alert Investigations instead.
- `alert_ai_daily_autonomous_token_limit` (Number) Maximum tokens per UTC day that autonomous alert-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new alert-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit.
- `alert_ai_daily_fix_task_limit` (Number) Maximum AI fix tasks derived from alerts that may be created per UTC day for this project. Unset means no limit; 0 pauses alert AI fix tasks entirely.
- `alert_ai_investigation_time_limit_in_minutes` (Number) Stop an alert AI investigation after this many minutes and report what it found. Unset (the default) means no time limit — the investigation runs until it is done.
- `alert_ai_max_concurrent_investigations` (Number) How many alert AI investigations may run at the same time for this project. Unset means no limit — every alert investigation starts right away. Minimum 1 — pause alert investigations with the Enable Automatic Alert Investigation toggle or a daily token limit of 0 instead.
- `alert_episode_number_prefix` (String) Custom prefix for alert episode numbers (e.g., 'AE-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber alert episodes that already exist.
- `alert_episode_numbers_start_at` (Number) Where alert episode numbers start: the next alert episode gets this number if no alert episode has reached it yet, and later ones count up from it. Use it to carry on from another tool's numbers. A new value must be higher than every alert episode number used so far, so a number is never used twice. Empty: numbers carry on from the highest one used. Existing alert episodes keep their numbers.
- `alert_investigation_dedupe_window_minutes` (Number) Repeat alerts from the same monitor within this many minutes are not re-investigated by AI — the first analysis stands. Unset or 0 means no cooldown, so every alert is investigated; at most 1440 minutes (a day).
- `alert_investigation_minimum_severity_id` (String) ID of the minimum AlertSeverity that triggers automatic investigation. The ID of a `oneuptime_alert_severity`.
- `alert_number_prefix` (String) Custom prefix for alert numbers (e.g., 'ALT-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber alerts that already exist.
- `alert_numbers_start_at` (Number) Where alert numbers start: the next alert gets this number if no alert has reached it yet, and later ones count up from it. Use it to carry on from another tool's numbers. A new value must be higher than every alert number used so far, so a number is never used twice. Empty: numbers carry on from the highest one used. Existing alerts keep their numbers.
- `audit_logs_retention_in_days` (Number) Number of days to retain audit log entries. Minimum 7, maximum 180.
- `auto_ai_recharge_by_balance_in_usd` (Number) Auto recharge amount in USD for AI services.
- `auto_archive_non_actionable_exceptions` (Boolean) When enabled, exception groups the AI triage classifies as expected denials (auth failures, plan/paywall rejections, scanner probes tripping intentional validation) are automatically archived so they stop surfacing in the unresolved list and never queue AI fix tasks. Groups classified as user errors or infrastructure conditions are NOT auto-archived — only clear expected denials are. Archiving is reversible from the Archived tab. On for new projects created in OneUptime; projects that existed before keep their setting.
- `auto_recharge_ai_when_current_balance_falls_in_usd` (Number) Auto recharge is triggered when current balance falls to this amount in USD for AI services.
- `auto_recharge_sms_or_call_by_balance_in_usd` (Number) Auto recharge amount in USD for SMS, Call, and WhatsApp.
- `auto_recharge_sms_or_call_when_current_balance_falls_in_usd` (Number) Auto recharge is triggered when current balance falls to this amount in USD for SMS, Call, and WhatsApp.
- `billing_marketplace` (String) The cloud marketplace that bills this project (AWS, Azure or GoogleCloud). Empty when the project is billed by a card on file.
- `business_details` (String) Business legal name, address and any tax information to appear on invoices.
- `business_details_country` (String) Two-letter ISO country code for billing address (e.g., US, GB, DE).
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `data_residency` (String) Where this project's data is hosted. Set by OneUptime staff on OneUptime Cloud.
- `default_metric_cardinality_budget` (Number) Project-wide default max distinct series per metric. Services without a per-service override use this value.
- `default_telemetry_retention_in_days` (Number) Project-wide default number of days to retain telemetry data (logs, traces, metrics). Services without a per-service override use this value.
- `disable_on_call_notification_fallback` (Boolean) When enabled, a page routed to a responder with no matching notification rule fails instead of falling back to their verified notification methods.
- `do_not_add_global_probes_by_default_on_new_monitors` (Boolean) If enabled, global probes will NOT be automatically added to new monitors. Enable this only if you are using ONLY custom probes to monitor your resources.
- `enable_ai` (Boolean) Master switch for AI in this project. When disabled, every AI feature stops: Ask AI, investigations, postmortem drafts, auto-remediation and AI commands on Runners.
- `enable_ai_insights` (Boolean) When enabled, OneUptime AI continuously watches this project's telemetry with deterministic statistical sensors (error-log spikes, exception novelty and spikes, trace-latency regressions, week-over-week metric drift) and files quiet Insights — never pages, never opens incidents. Each new insight also gets a budgeted, read-only AI triage analysis when an LLM provider is configured. On for new projects created in OneUptime; projects that existed before keep their setting.
- `enable_audit_logs` (Boolean) When enabled, changes to resources in this project are recorded as audit log entries.
- `enable_auto_recharge_ai_balance` (Boolean) Enable auto recharge for AI balance for this project.
- `enable_auto_recharge_sms_or_call_balance` (Boolean) Enable auto recharge for SMS, Call, and WhatsApp balance for this project.
- `enable_automatic_alert_investigation` (Boolean) When enabled, OneUptime's AI SRE automatically investigates every new alert and posts a cited root cause analysis to the alert timeline. On for new projects created in OneUptime; projects that existed before keep their setting. Requires AI to be enabled and an LLM provider to be configured.
- `enable_automatic_alert_remediation` (Boolean) When enabled, OneUptime fixes new alerts automatically: OneUptime AI fixes each one on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, and AI rules - when there are any - choose which alerts are fixed, which runbooks run and whether a person approves first. Off by default, for new projects too. Any AI investigation of the alert settles first. Requires AI to be enabled. Fix and telemetry pull requests open only while this is on, as each AI Rule says (openCodeFixPullRequests, openTelemetryPullRequests - on by default, and on with no rule).
- `enable_automatic_incident_investigation` (Boolean) When enabled, OneUptime's AI SRE automatically investigates every new incident and posts a cited root cause analysis to the incident timeline; any auto-remediation for the incident waits until that investigation settles. On for new projects created in OneUptime; projects that existed before keep their setting. Drafting a postmortem when an incident resolves is a separate setting (Enable Automatic Postmortem Draft). Requires AI to be enabled and an LLM provider to be configured.
- `enable_automatic_incident_remediation` (Boolean) When enabled, OneUptime fixes new incidents automatically: OneUptime AI fixes each one on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, and AI rules - when there are any - choose which incidents are fixed, which runbooks run and whether a person approves first. Off by default, for new projects too. Any AI investigation of the incident settles first. Requires AI to be enabled. Fix and telemetry pull requests open only while this is on, as each AI Rule says (openCodeFixPullRequests, openTelemetryPullRequests - on by default, and on with no rule).
- `enable_automatic_postmortem_draft` (Boolean) When enabled, OneUptime's AI SRE drafts a postmortem from the incident's timeline and telemetry when an incident is resolved, for a human to review and edit. It never overwrites a postmortem that already exists. On for new projects created in OneUptime; projects that existed before keep their setting. Requires AI to be enabled and an LLM provider to be configured.
- `enable_automatic_security_case_triage` (Boolean) When enabled, OneUptime AI triages every new security case a detection opens: it posts what happened, how likely it is to be real, the priority it deserves and the playbook that fits to the case's timeline. It changes nothing: a person decides the verdict and priority and runs the playbook. On for new projects created in OneUptime; projects that existed before keep their setting. Requires AI to be enabled and an LLM provider to be configured.
- `enable_call_notifications` (Boolean) Enable call notifications for this project.
- `enable_insight_fix_tasks` (Boolean) When enabled, insights whose deterministic evidence points at code (new or spiking exceptions with a resolvable repository, trace-latency regressions with span-tree findings) automatically queue an AI agent task that opens a pull request with a proposed fix, ready for review. Honors any open-PR cap set on the repository. Pull requests are always human-reviewed — nothing merges automatically. On for new projects created in OneUptime; projects that existed before keep their setting.
- `enable_sms_notifications` (Boolean) Enable SMS notifications for this project.
- `enable_telegram_notifications` (Boolean) Enable Telegram notifications for this project.
- `enable_whats_app_notifications` (Boolean) Enable WhatsApp notifications for this project.
- `finance_accounting_email` (String) Invoices, receipts and billing related notifications will be sent to these emails in addition to project owner. Separate multiple emails with a comma.
- `git_hub_app_installation_id` (String) The GitHub App installation ID for this project. This is set when the GitHub App is installed on the organization.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_ai_daily_autonomous_token_limit` (Number) Maximum tokens per UTC day that autonomous incident-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new incident-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit.
- `incident_ai_daily_fix_task_limit` (Number) Maximum AI fix tasks derived from incidents that may be created per UTC day for this project. Unset means no limit; 0 pauses incident AI fix tasks entirely.
- `incident_ai_investigation_time_limit_in_minutes` (Number) Stop an incident AI investigation after this many minutes and report what it found. Unset (the default) means no time limit — the investigation runs until it is done.
- `incident_ai_max_concurrent_investigations` (Number) How many incident AI investigations may run at the same time for this project. Unset means no limit — every incident investigation starts right away. Minimum 1 — pause incident investigations with the Enable Automatic Incident Investigation toggle or a daily token limit of 0 instead.
- `incident_episode_number_prefix` (String) Custom prefix for incident episode numbers (e.g., 'IE-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber incident episodes that already exist.
- `incident_episode_numbers_start_at` (Number) Where incident episode numbers start: the next incident episode gets this number if no incident episode has reached it yet, and later ones count up from it. Use it to carry on from another tool's numbers. A new value must be higher than every incident episode number used so far, so a number is never used twice. Empty: numbers carry on from the highest one used. Existing incident episodes keep their numbers.
- `incident_investigation_dedupe_window_minutes` (Number) Incidents affecting a monitor that AI investigated within this many minutes are not re-investigated — the first analysis stands. Unset or 0 means no cooldown, so every incident is investigated; at most 1440 minutes (a day).
- `incident_investigation_minimum_severity_id` (String) ID of the minimum incident severity that is investigated automatically by AI. The ID of a `oneuptime_incident_severity`.
- `incident_number_prefix` (String) Custom prefix for incident numbers (e.g., 'INC-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber incidents that already exist.
- `incident_numbers_start_at` (Number) Where incident numbers start: the next incident gets this number if no incident has reached it yet, and later ones count up from it. Use it to carry on from another tool's numbers. A new value must be higher than every incident number used so far, so a number is never used twice. Empty: numbers carry on from the highest one used. Existing incidents keep their numbers.
- `is_feature_flag_monitor_groups_enabled` (Boolean) Is Feature Flag Monitor Groups Enabled.
- `is_session_replay_allowed` (Boolean) When enabled, RUM applications in this project may record session replays if they are individually enabled too. On by default; switch it off here to stop session replay across the entire project in one place.
- `let_customer_support_access_project` (Boolean) OneUptime customer support can access this project. This is used for debugging purposes.
- `limit_alerts_to_one_episode` (Boolean) When enabled, an alert can be in only one alert episode: adding an alert that is already in an episode to another one is refused, by hand, through the API, workflows and grouping rules alike. Remove it from its episode first to move it. Alerts that are already in several episodes stay in them. Off by default.
- `limit_incidents_to_one_episode` (Boolean) When enabled, an incident can be in only one incident episode: adding an incident that is already in an episode to another one is refused, by hand, through the API, workflows and grouping rules alike. Remove it from its episode first to move it. Incidents that are already in several episodes stay in them. Off by default.
- `name` (String) Any friendly name of this object.
- `payment_provider_customer_id` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]
- `payment_provider_metered_subscription_id` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]
- `payment_provider_metered_subscription_status` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]
- `payment_provider_plan_id` (String) Permissions - Create: [Logged in User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [Project Owner, Manage Billing]
- `payment_provider_promo_code` (String) Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]
- `payment_provider_subscription_id` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]
- `payment_provider_subscription_seats` (Number) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]
- `payment_provider_subscription_status` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]
- `plan_name` (String) Name of the plan this project is subscribed to.
- `record_and_transcribe_incident_calls` (Boolean) When enabled, OneUptime asks Zoom, Google Meet and Microsoft Teams to record and transcribe every incident call it starts, and adds each call's transcript to its incident once the call is over. Slack huddles and meeting links are not recorded; their transcripts can be added by hand. Off by default.
- `require_sso_for_login` (Boolean) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Project Admin, Edit Project]
- `require_sso_with_sso_provider_id` (String) If set, SSO-enforced login for this project is only satisfied by an SSO token issued by this specific provider id (a Project SSO/OIDC or a Global SSO/OIDC). When null, any trusted SSO provider satisfies enforcement.
- `reseller_id` (String) ID of your OneUptime Reseller in which this object belongs.
- `reseller_plan_id` (String) ID of your OneUptime Reseller Plan in which this object belongs.
- `resolve_linked_alerts_when_incident_resolved` (Boolean) When enabled, resolving an incident also resolves every alert linked to it, except alerts that are still linked to another incident that is not resolved yet. Alerts linked to an incident that is already resolved are resolved as they are linked. On for new projects created in OneUptime; projects that existed before keep their setting.
- `scheduled_maintenance_number_prefix` (String) Custom prefix for scheduled maintenance numbers (e.g., 'SM-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber scheduled maintenance events that already exist.
- `scheduled_maintenance_numbers_start_at` (Number) Where scheduled maintenance event numbers start: the next scheduled maintenance event gets this number if no scheduled maintenance event has reached it yet, and later ones count up from it. Use it to carry on from another tool's numbers. A new value must be higher than every scheduled maintenance event number used so far, so a number is never used twice. Empty: numbers carry on from the highest one used. Existing scheduled maintenance events keep their numbers.
- `send_invoices_by_email` (Boolean) When enabled, invoices will be automatically sent to the finance/accounting email when they are generated.
- `slug` (String) Friendly globally unique name for your object.
- `sms_or_call_current_balance_in_usd_cents` (Number) Balance in USD for SMS, Call, and WhatsApp.
- `store_system_events_in_audit_logs` (Boolean) When enabled, audit logs will also include events triggered by the system. By default, only events triggered by users are recorded.
- `use_global_runners` (Boolean) When enabled, runbook steps that do not pick a Runner of their own run on OneUptime's global Runners, in a sandbox that can reach the internet but not your network. Turn it off to keep every step on the Runners you install. On by default.
- `workflow_runs_in_last30_days` (Number) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `default_metric_downsampling_retention_days` (String) Project-wide default retention for each downsampling tier (raw, 1m, 5m, 1h, 1d) in days. A JSON value: write it with `jsonencode()`.
- `telemetry_retention_config` (String) Project-wide per-pillar retention overrides for telemetry data (logs by severity, traces by status, metrics, profiles). Falls back to defaultTelemetryRetentionInDays when a pillar or bucket is not set. A JSON value: write it with `jsonencode()`.
- `trial_ends_at` (String) Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]
- `updated_at` (String) Date and Time when the object was updated.
