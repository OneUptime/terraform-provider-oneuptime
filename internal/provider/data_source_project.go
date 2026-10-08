package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProjectDataSource{}

func NewProjectDataSource() datasource.DataSource {
    return &ProjectDataSource{}
}

// ProjectDataSource defines the data source implementation.
type ProjectDataSource struct {
    client *Client
}

// ProjectDataSourceModel describes the data source data model.
type ProjectDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    Name types.String `tfsdk:"name"`
    Slug types.String `tfsdk:"slug"`
    PaymentProviderPlanId types.String `tfsdk:"payment_provider_plan_id"`
    PaymentProviderSubscriptionId types.String `tfsdk:"payment_provider_subscription_id"`
    PaymentProviderMeteredSubscriptionId types.String `tfsdk:"payment_provider_metered_subscription_id"`
    PaymentProviderSubscriptionSeats types.Number `tfsdk:"payment_provider_subscription_seats"`
    TrialEndsAt types.String `tfsdk:"trial_ends_at"`
    PaymentProviderCustomerId types.String `tfsdk:"payment_provider_customer_id"`
    BusinessDetails types.String `tfsdk:"business_details"`
    BusinessDetailsCountry types.String `tfsdk:"business_details_country"`
    FinanceAccountingEmail types.String `tfsdk:"finance_accounting_email"`
    PaymentProviderSubscriptionStatus types.String `tfsdk:"payment_provider_subscription_status"`
    PaymentProviderMeteredSubscriptionStatus types.String `tfsdk:"payment_provider_metered_subscription_status"`
    PaymentProviderPromoCode types.String `tfsdk:"payment_provider_promo_code"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsFeatureFlagMonitorGroupsEnabled types.Bool `tfsdk:"is_feature_flag_monitor_groups_enabled"`
    WorkflowRunsInLast30Days types.Number `tfsdk:"workflow_runs_in_last30_days"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    RequireSsoWithSsoProviderId types.String `tfsdk:"require_sso_with_sso_provider_id"`
    IncidentNumberPrefix types.String `tfsdk:"incident_number_prefix"`
    AlertNumberPrefix types.String `tfsdk:"alert_number_prefix"`
    ScheduledMaintenanceNumberPrefix types.String `tfsdk:"scheduled_maintenance_number_prefix"`
    IncidentEpisodeNumberPrefix types.String `tfsdk:"incident_episode_number_prefix"`
    AlertEpisodeNumberPrefix types.String `tfsdk:"alert_episode_number_prefix"`
    SmsOrCallCurrentBalanceInUsdCents types.Number `tfsdk:"sms_or_call_current_balance_in_usd_cents"`
    AutoRechargeSmsOrCallByBalanceInUsd types.Number `tfsdk:"auto_recharge_sms_or_call_by_balance_in_usd"`
    AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd types.Number `tfsdk:"auto_recharge_sms_or_call_when_current_balance_falls_in_usd"`
    EnableSmsNotifications types.Bool `tfsdk:"enable_sms_notifications"`
    EnableWhatsAppNotifications types.Bool `tfsdk:"enable_whats_app_notifications"`
    EnableTelegramNotifications types.Bool `tfsdk:"enable_telegram_notifications"`
    EnableCallNotifications types.Bool `tfsdk:"enable_call_notifications"`
    DisableOnCallNotificationFallback types.Bool `tfsdk:"disable_on_call_notification_fallback"`
    EnableAutoRechargeSmsOrCallBalance types.Bool `tfsdk:"enable_auto_recharge_sms_or_call_balance"`
    AiCurrentBalanceInUsdCents types.Number `tfsdk:"ai_current_balance_in_usd_cents"`
    AutoAiRechargeByBalanceInUsd types.Number `tfsdk:"auto_ai_recharge_by_balance_in_usd"`
    AutoRechargeAiWhenCurrentBalanceFallsInUsd types.Number `tfsdk:"auto_recharge_ai_when_current_balance_falls_in_usd"`
    EnableAi types.Bool `tfsdk:"enable_ai"`
    AiDailyTokenLimit types.Number `tfsdk:"ai_daily_token_limit"`
    AiDailySpendLimitInUsd types.Number `tfsdk:"ai_daily_spend_limit_in_usd"`
    EnableAutomaticIncidentInvestigation types.Bool `tfsdk:"enable_automatic_incident_investigation"`
    EnableAutomaticAlertInvestigation types.Bool `tfsdk:"enable_automatic_alert_investigation"`
    EnableAutomaticIncidentRemediation types.Bool `tfsdk:"enable_automatic_incident_remediation"`
    EnableAutomaticAlertRemediation types.Bool `tfsdk:"enable_automatic_alert_remediation"`
    EnableAutomaticPostmortemDraft types.Bool `tfsdk:"enable_automatic_postmortem_draft"`
    AcknowledgeLinkedAlertsWhenIncidentAcknowledged types.Bool `tfsdk:"acknowledge_linked_alerts_when_incident_acknowledged"`
    ResolveLinkedAlertsWhenIncidentResolved types.Bool `tfsdk:"resolve_linked_alerts_when_incident_resolved"`
    EnableIncidentInstrumentationFixTasks types.Bool `tfsdk:"enable_incident_instrumentation_fix_tasks"`
    EnableAlertInstrumentationFixTasks types.Bool `tfsdk:"enable_alert_instrumentation_fix_tasks"`
    EnableAutomaticIncidentCodeFixes types.Bool `tfsdk:"enable_automatic_incident_code_fixes"`
    EnableAutomaticAlertCodeFixes types.Bool `tfsdk:"enable_automatic_alert_code_fixes"`
    EnableAiInsights types.Bool `tfsdk:"enable_ai_insights"`
    EnableInsightFixTasks types.Bool `tfsdk:"enable_insight_fix_tasks"`
    AutoArchiveNonActionableExceptions types.Bool `tfsdk:"auto_archive_non_actionable_exceptions"`
    AlertInvestigationMinimumSeverityId types.String `tfsdk:"alert_investigation_minimum_severity_id"`
    AiDailyAutonomousTokenLimit types.Number `tfsdk:"ai_daily_autonomous_token_limit"`
    IncidentAiDailyAutonomousTokenLimit types.Number `tfsdk:"incident_ai_daily_autonomous_token_limit"`
    AlertAiDailyAutonomousTokenLimit types.Number `tfsdk:"alert_ai_daily_autonomous_token_limit"`
    AiDailyFixTaskLimit types.Number `tfsdk:"ai_daily_fix_task_limit"`
    IncidentAiDailyFixTaskLimit types.Number `tfsdk:"incident_ai_daily_fix_task_limit"`
    AlertAiDailyFixTaskLimit types.Number `tfsdk:"alert_ai_daily_fix_task_limit"`
    AlertInvestigationDedupeWindowMinutes types.Number `tfsdk:"alert_investigation_dedupe_window_minutes"`
    IncidentInvestigationMinimumSeverityId types.String `tfsdk:"incident_investigation_minimum_severity_id"`
    IncidentInvestigationDedupeWindowMinutes types.Number `tfsdk:"incident_investigation_dedupe_window_minutes"`
    AiMaxConcurrentInvestigations types.Number `tfsdk:"ai_max_concurrent_investigations"`
    IncidentAiMaxConcurrentInvestigations types.Number `tfsdk:"incident_ai_max_concurrent_investigations"`
    AlertAiMaxConcurrentInvestigations types.Number `tfsdk:"alert_ai_max_concurrent_investigations"`
    IncidentAiInvestigationTimeLimitInMinutes types.Number `tfsdk:"incident_ai_investigation_time_limit_in_minutes"`
    AlertAiInvestigationTimeLimitInMinutes types.Number `tfsdk:"alert_ai_investigation_time_limit_in_minutes"`
    EnableAutoRechargeAiBalance types.Bool `tfsdk:"enable_auto_recharge_ai_balance"`
    SendInvoicesByEmail types.Bool `tfsdk:"send_invoices_by_email"`
    PlanName types.String `tfsdk:"plan_name"`
    DataResidency types.String `tfsdk:"data_residency"`
    ResellerId types.String `tfsdk:"reseller_id"`
    ResellerPlanId types.String `tfsdk:"reseller_plan_id"`
    LetCustomerSupportAccessProject types.Bool `tfsdk:"let_customer_support_access_project"`
    DoNotAddGlobalProbesByDefaultOnNewMonitors types.Bool `tfsdk:"do_not_add_global_probes_by_default_on_new_monitors"`
    GitHubAppInstallationId types.String `tfsdk:"git_hub_app_installation_id"`
    DefaultMetricCardinalityBudget types.Number `tfsdk:"default_metric_cardinality_budget"`
    DefaultTelemetryRetentionInDays types.Number `tfsdk:"default_telemetry_retention_in_days"`
    TelemetryRetentionConfig types.String `tfsdk:"telemetry_retention_config"`
    DefaultMetricDownsamplingRetentionDays types.String `tfsdk:"default_metric_downsampling_retention_days"`
    EnableAuditLogs types.Bool `tfsdk:"enable_audit_logs"`
    IsSessionReplayAllowed types.Bool `tfsdk:"is_session_replay_allowed"`
    AuditLogsRetentionInDays types.Number `tfsdk:"audit_logs_retention_in_days"`
    StoreSystemEventsInAuditLogs types.Bool `tfsdk:"store_system_events_in_audit_logs"`
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "OneUptime Project, and everything happens inside it Look up an existing project by `id`, or by any of its other arguments (`name`, `acknowledge_linked_alerts_when_incident_acknowledged`, `ai_current_balance_in_usd_cents`, ...): each one set must match, and exactly one project may match them all.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Leave unset to look up by the other arguments instead.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Optional: true,
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Optional: true,
                Computed: true,
            },
            "payment_provider_plan_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
            },
            "payment_provider_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "payment_provider_metered_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "payment_provider_subscription_seats": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "trial_ends_at": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_customer_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "business_details": schema.StringAttribute{
                MarkdownDescription: "Business legal name, address and any tax information to appear on invoices.",
                Optional: true,
                Computed: true,
            },
            "business_details_country": schema.StringAttribute{
                MarkdownDescription: "Two-letter ISO country code for billing address (e.g., US, GB, DE).",
                Optional: true,
                Computed: true,
            },
            "finance_accounting_email": schema.StringAttribute{
                MarkdownDescription: "Invoices, receipts and billing related notifications will be sent to these emails in addition to project owner. Separate multiple emails with a comma.",
                Optional: true,
                Computed: true,
            },
            "payment_provider_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "payment_provider_metered_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "payment_provider_promo_code": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_feature_flag_monitor_groups_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Feature Flag Monitor Groups Enabled.",
                Optional: true,
                Computed: true,
            },
            "workflow_runs_in_last30_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
            },
            "require_sso_with_sso_provider_id": schema.StringAttribute{
                MarkdownDescription: "If set, SSO-enforced login for this project is only satisfied by an SSO token issued by this specific provider id (a Project SSO/OIDC or a Global SSO/OIDC). When null, any trusted SSO provider satisfies enforcement.",
                Optional: true,
                Computed: true,
            },
            "incident_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident numbers (e.g., 'INC-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber incidents that already exist.",
                Optional: true,
                Computed: true,
            },
            "alert_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert numbers (e.g., 'ALT-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber alerts that already exist.",
                Optional: true,
                Computed: true,
            },
            "scheduled_maintenance_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for scheduled maintenance numbers (e.g., 'SM-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber scheduled maintenance events that already exist.",
                Optional: true,
                Computed: true,
            },
            "incident_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident episode numbers (e.g., 'IE-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber incident episodes that already exist.",
                Optional: true,
                Computed: true,
            },
            "alert_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert episode numbers (e.g., 'AE-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber alert episodes that already exist.",
                Optional: true,
                Computed: true,
            },
            "sms_or_call_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for SMS, Call, and WhatsApp.",
                Optional: true,
                Computed: true,
            },
            "auto_recharge_sms_or_call_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for SMS, Call, and WhatsApp.",
                Optional: true,
                Computed: true,
            },
            "auto_recharge_sms_or_call_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for SMS, Call, and WhatsApp.",
                Optional: true,
                Computed: true,
            },
            "enable_sms_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable SMS notifications for this project.",
                Optional: true,
                Computed: true,
            },
            "enable_whats_app_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable WhatsApp notifications for this project.",
                Optional: true,
                Computed: true,
            },
            "enable_telegram_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable Telegram notifications for this project.",
                Optional: true,
                Computed: true,
            },
            "enable_call_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable call notifications for this project.",
                Optional: true,
                Computed: true,
            },
            "disable_on_call_notification_fallback": schema.BoolAttribute{
                MarkdownDescription: "When enabled, a page routed to a responder with no matching notification rule fails instead of falling back to their verified notification methods.",
                Optional: true,
                Computed: true,
            },
            "enable_auto_recharge_sms_or_call_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for SMS, Call, and WhatsApp balance for this project.",
                Optional: true,
                Computed: true,
            },
            "ai_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for AI services.",
                Optional: true,
                Computed: true,
            },
            "auto_ai_recharge_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for AI services.",
                Optional: true,
                Computed: true,
            },
            "auto_recharge_ai_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for AI services.",
                Optional: true,
                Computed: true,
            },
            "enable_ai": schema.BoolAttribute{
                MarkdownDescription: "Master switch for AI in this project. When disabled, every AI feature stops: Ask AI, investigations, postmortem drafts, auto-remediation and AI commands on Runners.",
                Optional: true,
                Computed: true,
            },
            "ai_daily_token_limit": schema.NumberAttribute{
                MarkdownDescription: "The most tokens OneUptime AI may use in this project each UTC day, across every AI feature: Ask AI, investigations, postmortem drafts, fix pull requests, insight triage, workflows, runbooks and Slack or Microsoft Teams questions. Once it is reached, new AI work is refused until midnight UTC. The incident and alert daily limits still apply under it. Unset means no limit; a limit is a whole number of at least 1 (to turn AI off, use Enable AI).",
                Optional: true,
                Computed: true,
            },
            "ai_daily_spend_limit_in_usd": schema.NumberAttribute{
                MarkdownDescription: "OneUptime Cloud: the most AI credits, in whole US dollars, OneUptime AI may spend in this project each UTC day. Only calls billed to the project's AI credits count, so it never stops AI that runs on the project's own LLM provider. Once it is reached, billed AI work is refused until midnight UTC. Ignored where AI is not billed (self-hosted). Unset means no limit; a limit is at least 1 (to turn AI off, use Enable AI).",
                Optional: true,
                Computed: true,
            },
            "enable_automatic_incident_investigation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime's AI SRE automatically investigates every new incident and posts a cited root cause analysis to the incident timeline; any auto-remediation for the incident waits until that investigation settles. On for new projects created in OneUptime; projects that existed before keep their setting. Drafting a postmortem when an incident resolves is a separate setting (Enable Automatic Postmortem Draft). Requires AI to be enabled and an LLM provider to be configured.",
                Optional: true,
                Computed: true,
            },
            "enable_automatic_alert_investigation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime's AI SRE automatically investigates every new alert and posts a cited root cause analysis to the alert timeline. On for new projects created in OneUptime; projects that existed before keep their setting. Requires AI to be enabled and an LLM provider to be configured.",
                Optional: true,
                Computed: true,
            },
            "enable_automatic_incident_remediation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime fixes new incidents automatically: OneUptime AI fixes each one on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, and Auto Remediation Rules - when there are any - choose which incidents are fixed, which runbooks run and whether a person approves first. Off by default, for new projects too. Any AI investigation of the incident settles first. Requires AI to be enabled. It also holds the incident pull-request switches, enableAutomaticIncidentCodeFixes and enableIncidentInstrumentationFixTasks: they open pull requests only while this is on. The dashboard turns them on and off with it; through the API, set them in the same request.",
                Optional: true,
                Computed: true,
            },
            "enable_automatic_alert_remediation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime fixes new alerts automatically: OneUptime AI fixes each one on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, and Auto Remediation Rules - when there are any - choose which alerts are fixed, which runbooks run and whether a person approves first. Off by default, for new projects too. Any AI investigation of the alert settles first. Requires AI to be enabled. It also holds the alert pull-request switches, enableAutomaticAlertCodeFixes and enableAlertInstrumentationFixTasks: they open pull requests only while this is on. The dashboard turns them on and off with it; through the API, set them in the same request.",
                Optional: true,
                Computed: true,
            },
            "enable_automatic_postmortem_draft": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime's AI SRE drafts a postmortem from the incident's timeline and telemetry when an incident is resolved, for a human to review and edit. It never overwrites a postmortem that already exists. On for new projects created in OneUptime; projects that existed before keep their setting. Requires AI to be enabled and an LLM provider to be configured.",
                Optional: true,
                Computed: true,
            },
            "acknowledge_linked_alerts_when_incident_acknowledged": schema.BoolAttribute{
                MarkdownDescription: "When enabled, acknowledging an incident also acknowledges every alert linked to it. This stops those alerts' on-call escalations, and their reminders only when the alert reminder rule is set to stop on Acknowledged. Alerts linked to an incident that is already acknowledged are acknowledged as they are linked. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
            },
            "resolve_linked_alerts_when_incident_resolved": schema.BoolAttribute{
                MarkdownDescription: "When enabled, resolving an incident also resolves every alert linked to it, except alerts that are still linked to another incident that is not resolved yet. Alerts linked to an incident that is already resolved are resolved as they are linked. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
            },
            "enable_incident_instrumentation_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an incident AI investigation that ends inconclusive (telemetry was insufficient to determine a root cause) automatically queues an AI agent task that opens a pull request adding the missing instrumentation to the implicated code paths. Requires a repository connected through the GitHub App. Pull requests are always human-reviewed — nothing merges automatically. Part of fixing: it acts only while enableAutomaticIncidentRemediation (Fix new incidents automatically) is on, and the dashboard turns it on and off with that switch. Off for new projects.",
                Optional: true,
                Computed: true,
            },
            "enable_alert_instrumentation_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an alert AI investigation that ends inconclusive (telemetry was insufficient to determine a root cause) automatically queues an AI agent task that opens a pull request adding the missing instrumentation to the implicated code paths. Requires a repository connected through the GitHub App. Pull requests are always human-reviewed — nothing merges automatically. Part of fixing: it acts only while enableAutomaticAlertRemediation (Fix new alerts automatically) is on, and the dashboard turns it on and off with that switch. Off for new projects.",
                Optional: true,
                Computed: true,
            },
            "enable_automatic_incident_code_fixes": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an incident AI investigation that ends with a confident, evidenced root cause analysis and recommends a repository code change automatically queues an AI agent task that opens a fix pull request, ready for review, from that analysis — the automatic form of the 'Open Fix PR from this analysis' button. Operational, infrastructure, external, user-error and inconclusive findings do not offer or open code-fix pull requests. Requires a repository connected through the GitHub App and a Runner with the code-fix capability. Pull requests are always human-reviewed — nothing merges automatically. Part of fixing: it acts only while enableAutomaticIncidentRemediation (Fix new incidents automatically) is on, and the dashboard turns it on and off with that switch. Off for new projects.",
                Optional: true,
                Computed: true,
            },
            "enable_automatic_alert_code_fixes": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an alert AI investigation that ends with a confident, evidenced root cause analysis and recommends a repository code change automatically queues an AI agent task that opens a fix pull request, ready for review, from that analysis — the automatic form of the 'Open Fix PR from this analysis' button. Operational, infrastructure, external, user-error and inconclusive findings do not offer or open code-fix pull requests. Requires a repository connected through the GitHub App and a Runner with the code-fix capability. Pull requests are always human-reviewed — nothing merges automatically. Part of fixing: it acts only while enableAutomaticAlertRemediation (Fix new alerts automatically) is on, and the dashboard turns it on and off with that switch. Off for new projects.",
                Optional: true,
                Computed: true,
            },
            "enable_ai_insights": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime AI continuously watches this project's telemetry with deterministic statistical sensors (error-log spikes, exception novelty and spikes, trace-latency regressions, week-over-week metric drift) and files quiet Insights — never pages, never opens incidents. Each new insight also gets a budgeted, read-only AI triage analysis when an LLM provider is configured. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
            },
            "enable_insight_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "When enabled, insights whose deterministic evidence points at code (new or spiking exceptions with a resolvable repository, trace-latency regressions with span-tree findings) automatically queue an AI agent task that opens a pull request with a proposed fix, ready for review. Honors any open-PR cap set on the repository. Pull requests are always human-reviewed — nothing merges automatically. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
            },
            "auto_archive_non_actionable_exceptions": schema.BoolAttribute{
                MarkdownDescription: "When enabled, exception groups the AI triage classifies as expected denials (auth failures, plan/paywall rejections, scanner probes tripping intentional validation) are automatically archived so they stop surfacing in the unresolved list and never queue AI fix tasks. Groups classified as user errors or infrastructure conditions are NOT auto-archived — only clear expected denials are. Archiving is reversible from the Archived tab. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
            },
            "alert_investigation_minimum_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the minimum AlertSeverity that triggers automatic investigation. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
            },
            "ai_daily_autonomous_token_limit": schema.NumberAttribute{
                MarkdownDescription: "Legacy setting, no longer enforced: autonomous AI work that is not associated with an incident or alert has no daily token limit. Use the Daily Incident AI Token Limit and Daily Alert AI Token Limit instead.",
                Optional: true,
                Computed: true,
            },
            "incident_ai_daily_autonomous_token_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum tokens per UTC day that autonomous incident-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new incident-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit.",
                Optional: true,
                Computed: true,
            },
            "alert_ai_daily_autonomous_token_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum tokens per UTC day that autonomous alert-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new alert-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit.",
                Optional: true,
                Computed: true,
            },
            "ai_daily_fix_task_limit": schema.NumberAttribute{
                MarkdownDescription: "Legacy setting, no longer enforced: AI fix tasks that are not associated with an incident or alert have no daily limit. Use the Daily Incident AI Fix Task Limit and Daily Alert AI Fix Task Limit instead.",
                Optional: true,
                Computed: true,
            },
            "incident_ai_daily_fix_task_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum AI fix tasks derived from incidents that may be created per UTC day for this project. Unset means no limit; 0 pauses incident AI fix tasks entirely.",
                Optional: true,
                Computed: true,
            },
            "alert_ai_daily_fix_task_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum AI fix tasks derived from alerts that may be created per UTC day for this project. Unset means no limit; 0 pauses alert AI fix tasks entirely.",
                Optional: true,
                Computed: true,
            },
            "alert_investigation_dedupe_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Repeat alerts from the same monitor within this many minutes are not re-investigated by AI — the first analysis stands. Unset or 0 means no cooldown, so every alert is investigated; at most 1440 minutes (a day).",
                Optional: true,
                Computed: true,
            },
            "incident_investigation_minimum_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the minimum incident severity that is investigated automatically by AI. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
            },
            "incident_investigation_dedupe_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Incidents affecting a monitor that AI investigated within this many minutes are not re-investigated — the first analysis stands. Unset or 0 means no cooldown, so every incident is investigated; at most 1440 minutes (a day).",
                Optional: true,
                Computed: true,
            },
            "ai_max_concurrent_investigations": schema.NumberAttribute{
                MarkdownDescription: "Legacy setting, no longer enforced: AI investigations that are not associated with an incident or alert have no concurrency limit. Use the Max Concurrent Incident Investigations and Max Concurrent Alert Investigations instead.",
                Optional: true,
                Computed: true,
            },
            "incident_ai_max_concurrent_investigations": schema.NumberAttribute{
                MarkdownDescription: "How many incident AI investigations may run at the same time for this project. Unset means no limit — every incident investigation starts right away. Minimum 1 — pause incident investigations with the Enable Automatic Incident Investigation toggle or a daily token limit of 0 instead.",
                Optional: true,
                Computed: true,
            },
            "alert_ai_max_concurrent_investigations": schema.NumberAttribute{
                MarkdownDescription: "How many alert AI investigations may run at the same time for this project. Unset means no limit — every alert investigation starts right away. Minimum 1 — pause alert investigations with the Enable Automatic Alert Investigation toggle or a daily token limit of 0 instead.",
                Optional: true,
                Computed: true,
            },
            "incident_ai_investigation_time_limit_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Stop an incident AI investigation after this many minutes and report what it found. Unset (the default) means no time limit — the investigation runs until it is done.",
                Optional: true,
                Computed: true,
            },
            "alert_ai_investigation_time_limit_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Stop an alert AI investigation after this many minutes and report what it found. Unset (the default) means no time limit — the investigation runs until it is done.",
                Optional: true,
                Computed: true,
            },
            "enable_auto_recharge_ai_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for AI balance for this project.",
                Optional: true,
                Computed: true,
            },
            "send_invoices_by_email": schema.BoolAttribute{
                MarkdownDescription: "When enabled, invoices will be automatically sent to the finance/accounting email when they are generated.",
                Optional: true,
                Computed: true,
            },
            "plan_name": schema.StringAttribute{
                MarkdownDescription: "Name of the plan this project is subscribed to.",
                Optional: true,
                Computed: true,
            },
            "data_residency": schema.StringAttribute{
                MarkdownDescription: "Where this project's data is hosted. Set by OneUptime staff on OneUptime Cloud.",
                Optional: true,
                Computed: true,
            },
            "reseller_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Reseller in which this object belongs.",
                Optional: true,
                Computed: true,
            },
            "reseller_plan_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Reseller Plan in which this object belongs.",
                Optional: true,
                Computed: true,
            },
            "let_customer_support_access_project": schema.BoolAttribute{
                MarkdownDescription: "OneUptime customer support can access this project. This is used for debugging purposes.",
                Optional: true,
                Computed: true,
            },
            "do_not_add_global_probes_by_default_on_new_monitors": schema.BoolAttribute{
                MarkdownDescription: "If enabled, global probes will NOT be automatically added to new monitors. Enable this only if you are using ONLY custom probes to monitor your resources.",
                Optional: true,
                Computed: true,
            },
            "git_hub_app_installation_id": schema.StringAttribute{
                MarkdownDescription: "The GitHub App installation ID for this project. This is set when the GitHub App is installed on the organization.",
                Optional: true,
                Computed: true,
            },
            "default_metric_cardinality_budget": schema.NumberAttribute{
                MarkdownDescription: "Project-wide default max distinct series per metric. Services without a per-service override use this value.",
                Optional: true,
                Computed: true,
            },
            "default_telemetry_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "Project-wide default number of days to retain telemetry data (logs, traces, metrics). Services without a per-service override use this value.",
                Optional: true,
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Project-wide per-pillar retention overrides for telemetry data (logs by severity, traces by status, metrics, profiles). Falls back to defaultTelemetryRetentionInDays when a pillar or bucket is not set. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "default_metric_downsampling_retention_days": schema.StringAttribute{
                MarkdownDescription: "Project-wide default retention for each downsampling tier (raw, 1m, 5m, 1h, 1d) in days. A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "enable_audit_logs": schema.BoolAttribute{
                MarkdownDescription: "When enabled, changes to resources in this project are recorded as audit log entries.",
                Optional: true,
                Computed: true,
            },
            "is_session_replay_allowed": schema.BoolAttribute{
                MarkdownDescription: "When enabled, RUM applications in this project may record session replays if they are individually enabled too. On by default; switch it off here to stop session replay across the entire project in one place.",
                Optional: true,
                Computed: true,
            },
            "audit_logs_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain audit log entries. Minimum 7, maximum 180.",
                Optional: true,
                Computed: true,
            },
            "store_system_events_in_audit_logs": schema.BoolAttribute{
                MarkdownDescription: "When enabled, audit logs will also include events triggered by the system. By default, only events triggered by users are recorded.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    d.client = client
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProjectDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        filters["name"] = data.Name.ValueString()
        filterNames = append(filterNames, "name = "+fmt.Sprintf("%q", data.Name.ValueString()))
    }
    if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
        filters["slug"] = data.Slug.ValueString()
        filterNames = append(filterNames, "slug = "+fmt.Sprintf("%q", data.Slug.ValueString()))
    }
    if !data.PaymentProviderPlanId.IsNull() && !data.PaymentProviderPlanId.IsUnknown() {
        filters["paymentProviderPlanId"] = data.PaymentProviderPlanId.ValueString()
        filterNames = append(filterNames, "payment_provider_plan_id = "+fmt.Sprintf("%q", data.PaymentProviderPlanId.ValueString()))
    }
    if !data.PaymentProviderSubscriptionId.IsNull() && !data.PaymentProviderSubscriptionId.IsUnknown() {
        filters["paymentProviderSubscriptionId"] = data.PaymentProviderSubscriptionId.ValueString()
        filterNames = append(filterNames, "payment_provider_subscription_id = "+fmt.Sprintf("%q", data.PaymentProviderSubscriptionId.ValueString()))
    }
    if !data.PaymentProviderMeteredSubscriptionId.IsNull() && !data.PaymentProviderMeteredSubscriptionId.IsUnknown() {
        filters["paymentProviderMeteredSubscriptionId"] = data.PaymentProviderMeteredSubscriptionId.ValueString()
        filterNames = append(filterNames, "payment_provider_metered_subscription_id = "+fmt.Sprintf("%q", data.PaymentProviderMeteredSubscriptionId.ValueString()))
    }
    if !data.PaymentProviderSubscriptionSeats.IsNull() && !data.PaymentProviderSubscriptionSeats.IsUnknown() {
        filters["paymentProviderSubscriptionSeats"] = lookupNumber(data.PaymentProviderSubscriptionSeats)
        filterNames = append(filterNames, "payment_provider_subscription_seats = "+data.PaymentProviderSubscriptionSeats.ValueBigFloat().String())
    }
    if !data.PaymentProviderCustomerId.IsNull() && !data.PaymentProviderCustomerId.IsUnknown() {
        filters["paymentProviderCustomerId"] = data.PaymentProviderCustomerId.ValueString()
        filterNames = append(filterNames, "payment_provider_customer_id = "+fmt.Sprintf("%q", data.PaymentProviderCustomerId.ValueString()))
    }
    if !data.BusinessDetails.IsNull() && !data.BusinessDetails.IsUnknown() {
        filters["businessDetails"] = data.BusinessDetails.ValueString()
        filterNames = append(filterNames, "business_details = "+fmt.Sprintf("%q", data.BusinessDetails.ValueString()))
    }
    if !data.BusinessDetailsCountry.IsNull() && !data.BusinessDetailsCountry.IsUnknown() {
        filters["businessDetailsCountry"] = data.BusinessDetailsCountry.ValueString()
        filterNames = append(filterNames, "business_details_country = "+fmt.Sprintf("%q", data.BusinessDetailsCountry.ValueString()))
    }
    if !data.FinanceAccountingEmail.IsNull() && !data.FinanceAccountingEmail.IsUnknown() {
        filters["financeAccountingEmail"] = data.FinanceAccountingEmail.ValueString()
        filterNames = append(filterNames, "finance_accounting_email = "+fmt.Sprintf("%q", data.FinanceAccountingEmail.ValueString()))
    }
    if !data.PaymentProviderSubscriptionStatus.IsNull() && !data.PaymentProviderSubscriptionStatus.IsUnknown() {
        filters["paymentProviderSubscriptionStatus"] = data.PaymentProviderSubscriptionStatus.ValueString()
        filterNames = append(filterNames, "payment_provider_subscription_status = "+fmt.Sprintf("%q", data.PaymentProviderSubscriptionStatus.ValueString()))
    }
    if !data.PaymentProviderMeteredSubscriptionStatus.IsNull() && !data.PaymentProviderMeteredSubscriptionStatus.IsUnknown() {
        filters["paymentProviderMeteredSubscriptionStatus"] = data.PaymentProviderMeteredSubscriptionStatus.ValueString()
        filterNames = append(filterNames, "payment_provider_metered_subscription_status = "+fmt.Sprintf("%q", data.PaymentProviderMeteredSubscriptionStatus.ValueString()))
    }
    if !data.PaymentProviderPromoCode.IsNull() && !data.PaymentProviderPromoCode.IsUnknown() {
        filters["paymentProviderPromoCode"] = data.PaymentProviderPromoCode.ValueString()
        filterNames = append(filterNames, "payment_provider_promo_code = "+fmt.Sprintf("%q", data.PaymentProviderPromoCode.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsFeatureFlagMonitorGroupsEnabled.IsNull() && !data.IsFeatureFlagMonitorGroupsEnabled.IsUnknown() {
        filters["isFeatureFlagMonitorGroupsEnabled"] = data.IsFeatureFlagMonitorGroupsEnabled.ValueBool()
        filterNames = append(filterNames, "is_feature_flag_monitor_groups_enabled = "+fmt.Sprintf("%t", data.IsFeatureFlagMonitorGroupsEnabled.ValueBool()))
    }
    if !data.WorkflowRunsInLast30Days.IsNull() && !data.WorkflowRunsInLast30Days.IsUnknown() {
        filters["workflowRunsInLast30Days"] = lookupNumber(data.WorkflowRunsInLast30Days)
        filterNames = append(filterNames, "workflow_runs_in_last30_days = "+data.WorkflowRunsInLast30Days.ValueBigFloat().String())
    }
    if !data.RequireSsoForLogin.IsNull() && !data.RequireSsoForLogin.IsUnknown() {
        filters["requireSsoForLogin"] = data.RequireSsoForLogin.ValueBool()
        filterNames = append(filterNames, "require_sso_for_login = "+fmt.Sprintf("%t", data.RequireSsoForLogin.ValueBool()))
    }
    if !data.RequireSsoWithSsoProviderId.IsNull() && !data.RequireSsoWithSsoProviderId.IsUnknown() {
        filters["requireSsoWithSsoProviderId"] = data.RequireSsoWithSsoProviderId.ValueString()
        filterNames = append(filterNames, "require_sso_with_sso_provider_id = "+fmt.Sprintf("%q", data.RequireSsoWithSsoProviderId.ValueString()))
    }
    if !data.IncidentNumberPrefix.IsNull() && !data.IncidentNumberPrefix.IsUnknown() {
        filters["incidentNumberPrefix"] = data.IncidentNumberPrefix.ValueString()
        filterNames = append(filterNames, "incident_number_prefix = "+fmt.Sprintf("%q", data.IncidentNumberPrefix.ValueString()))
    }
    if !data.AlertNumberPrefix.IsNull() && !data.AlertNumberPrefix.IsUnknown() {
        filters["alertNumberPrefix"] = data.AlertNumberPrefix.ValueString()
        filterNames = append(filterNames, "alert_number_prefix = "+fmt.Sprintf("%q", data.AlertNumberPrefix.ValueString()))
    }
    if !data.ScheduledMaintenanceNumberPrefix.IsNull() && !data.ScheduledMaintenanceNumberPrefix.IsUnknown() {
        filters["scheduledMaintenanceNumberPrefix"] = data.ScheduledMaintenanceNumberPrefix.ValueString()
        filterNames = append(filterNames, "scheduled_maintenance_number_prefix = "+fmt.Sprintf("%q", data.ScheduledMaintenanceNumberPrefix.ValueString()))
    }
    if !data.IncidentEpisodeNumberPrefix.IsNull() && !data.IncidentEpisodeNumberPrefix.IsUnknown() {
        filters["incidentEpisodeNumberPrefix"] = data.IncidentEpisodeNumberPrefix.ValueString()
        filterNames = append(filterNames, "incident_episode_number_prefix = "+fmt.Sprintf("%q", data.IncidentEpisodeNumberPrefix.ValueString()))
    }
    if !data.AlertEpisodeNumberPrefix.IsNull() && !data.AlertEpisodeNumberPrefix.IsUnknown() {
        filters["alertEpisodeNumberPrefix"] = data.AlertEpisodeNumberPrefix.ValueString()
        filterNames = append(filterNames, "alert_episode_number_prefix = "+fmt.Sprintf("%q", data.AlertEpisodeNumberPrefix.ValueString()))
    }
    if !data.SmsOrCallCurrentBalanceInUsdCents.IsNull() && !data.SmsOrCallCurrentBalanceInUsdCents.IsUnknown() {
        filters["smsOrCallCurrentBalanceInUSDCents"] = lookupNumber(data.SmsOrCallCurrentBalanceInUsdCents)
        filterNames = append(filterNames, "sms_or_call_current_balance_in_usd_cents = "+data.SmsOrCallCurrentBalanceInUsdCents.ValueBigFloat().String())
    }
    if !data.AutoRechargeSmsOrCallByBalanceInUsd.IsNull() && !data.AutoRechargeSmsOrCallByBalanceInUsd.IsUnknown() {
        filters["autoRechargeSmsOrCallByBalanceInUSD"] = lookupNumber(data.AutoRechargeSmsOrCallByBalanceInUsd)
        filterNames = append(filterNames, "auto_recharge_sms_or_call_by_balance_in_usd = "+data.AutoRechargeSmsOrCallByBalanceInUsd.ValueBigFloat().String())
    }
    if !data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.IsNull() && !data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.IsUnknown() {
        filters["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"] = lookupNumber(data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd)
        filterNames = append(filterNames, "auto_recharge_sms_or_call_when_current_balance_falls_in_usd = "+data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.ValueBigFloat().String())
    }
    if !data.EnableSmsNotifications.IsNull() && !data.EnableSmsNotifications.IsUnknown() {
        filters["enableSmsNotifications"] = data.EnableSmsNotifications.ValueBool()
        filterNames = append(filterNames, "enable_sms_notifications = "+fmt.Sprintf("%t", data.EnableSmsNotifications.ValueBool()))
    }
    if !data.EnableWhatsAppNotifications.IsNull() && !data.EnableWhatsAppNotifications.IsUnknown() {
        filters["enableWhatsAppNotifications"] = data.EnableWhatsAppNotifications.ValueBool()
        filterNames = append(filterNames, "enable_whats_app_notifications = "+fmt.Sprintf("%t", data.EnableWhatsAppNotifications.ValueBool()))
    }
    if !data.EnableTelegramNotifications.IsNull() && !data.EnableTelegramNotifications.IsUnknown() {
        filters["enableTelegramNotifications"] = data.EnableTelegramNotifications.ValueBool()
        filterNames = append(filterNames, "enable_telegram_notifications = "+fmt.Sprintf("%t", data.EnableTelegramNotifications.ValueBool()))
    }
    if !data.EnableCallNotifications.IsNull() && !data.EnableCallNotifications.IsUnknown() {
        filters["enableCallNotifications"] = data.EnableCallNotifications.ValueBool()
        filterNames = append(filterNames, "enable_call_notifications = "+fmt.Sprintf("%t", data.EnableCallNotifications.ValueBool()))
    }
    if !data.DisableOnCallNotificationFallback.IsNull() && !data.DisableOnCallNotificationFallback.IsUnknown() {
        filters["disableOnCallNotificationFallback"] = data.DisableOnCallNotificationFallback.ValueBool()
        filterNames = append(filterNames, "disable_on_call_notification_fallback = "+fmt.Sprintf("%t", data.DisableOnCallNotificationFallback.ValueBool()))
    }
    if !data.EnableAutoRechargeSmsOrCallBalance.IsNull() && !data.EnableAutoRechargeSmsOrCallBalance.IsUnknown() {
        filters["enableAutoRechargeSmsOrCallBalance"] = data.EnableAutoRechargeSmsOrCallBalance.ValueBool()
        filterNames = append(filterNames, "enable_auto_recharge_sms_or_call_balance = "+fmt.Sprintf("%t", data.EnableAutoRechargeSmsOrCallBalance.ValueBool()))
    }
    if !data.AiCurrentBalanceInUsdCents.IsNull() && !data.AiCurrentBalanceInUsdCents.IsUnknown() {
        filters["aiCurrentBalanceInUSDCents"] = lookupNumber(data.AiCurrentBalanceInUsdCents)
        filterNames = append(filterNames, "ai_current_balance_in_usd_cents = "+data.AiCurrentBalanceInUsdCents.ValueBigFloat().String())
    }
    if !data.AutoAiRechargeByBalanceInUsd.IsNull() && !data.AutoAiRechargeByBalanceInUsd.IsUnknown() {
        filters["autoAiRechargeByBalanceInUSD"] = lookupNumber(data.AutoAiRechargeByBalanceInUsd)
        filterNames = append(filterNames, "auto_ai_recharge_by_balance_in_usd = "+data.AutoAiRechargeByBalanceInUsd.ValueBigFloat().String())
    }
    if !data.AutoRechargeAiWhenCurrentBalanceFallsInUsd.IsNull() && !data.AutoRechargeAiWhenCurrentBalanceFallsInUsd.IsUnknown() {
        filters["autoRechargeAiWhenCurrentBalanceFallsInUSD"] = lookupNumber(data.AutoRechargeAiWhenCurrentBalanceFallsInUsd)
        filterNames = append(filterNames, "auto_recharge_ai_when_current_balance_falls_in_usd = "+data.AutoRechargeAiWhenCurrentBalanceFallsInUsd.ValueBigFloat().String())
    }
    if !data.EnableAi.IsNull() && !data.EnableAi.IsUnknown() {
        filters["enableAi"] = data.EnableAi.ValueBool()
        filterNames = append(filterNames, "enable_ai = "+fmt.Sprintf("%t", data.EnableAi.ValueBool()))
    }
    if !data.AiDailyTokenLimit.IsNull() && !data.AiDailyTokenLimit.IsUnknown() {
        filters["aiDailyTokenLimit"] = lookupNumber(data.AiDailyTokenLimit)
        filterNames = append(filterNames, "ai_daily_token_limit = "+data.AiDailyTokenLimit.ValueBigFloat().String())
    }
    if !data.AiDailySpendLimitInUsd.IsNull() && !data.AiDailySpendLimitInUsd.IsUnknown() {
        filters["aiDailySpendLimitInUSD"] = lookupNumber(data.AiDailySpendLimitInUsd)
        filterNames = append(filterNames, "ai_daily_spend_limit_in_usd = "+data.AiDailySpendLimitInUsd.ValueBigFloat().String())
    }
    if !data.EnableAutomaticIncidentInvestigation.IsNull() && !data.EnableAutomaticIncidentInvestigation.IsUnknown() {
        filters["enableAutomaticIncidentInvestigation"] = data.EnableAutomaticIncidentInvestigation.ValueBool()
        filterNames = append(filterNames, "enable_automatic_incident_investigation = "+fmt.Sprintf("%t", data.EnableAutomaticIncidentInvestigation.ValueBool()))
    }
    if !data.EnableAutomaticAlertInvestigation.IsNull() && !data.EnableAutomaticAlertInvestigation.IsUnknown() {
        filters["enableAutomaticAlertInvestigation"] = data.EnableAutomaticAlertInvestigation.ValueBool()
        filterNames = append(filterNames, "enable_automatic_alert_investigation = "+fmt.Sprintf("%t", data.EnableAutomaticAlertInvestigation.ValueBool()))
    }
    if !data.EnableAutomaticIncidentRemediation.IsNull() && !data.EnableAutomaticIncidentRemediation.IsUnknown() {
        filters["enableAutomaticIncidentRemediation"] = data.EnableAutomaticIncidentRemediation.ValueBool()
        filterNames = append(filterNames, "enable_automatic_incident_remediation = "+fmt.Sprintf("%t", data.EnableAutomaticIncidentRemediation.ValueBool()))
    }
    if !data.EnableAutomaticAlertRemediation.IsNull() && !data.EnableAutomaticAlertRemediation.IsUnknown() {
        filters["enableAutomaticAlertRemediation"] = data.EnableAutomaticAlertRemediation.ValueBool()
        filterNames = append(filterNames, "enable_automatic_alert_remediation = "+fmt.Sprintf("%t", data.EnableAutomaticAlertRemediation.ValueBool()))
    }
    if !data.EnableAutomaticPostmortemDraft.IsNull() && !data.EnableAutomaticPostmortemDraft.IsUnknown() {
        filters["enableAutomaticPostmortemDraft"] = data.EnableAutomaticPostmortemDraft.ValueBool()
        filterNames = append(filterNames, "enable_automatic_postmortem_draft = "+fmt.Sprintf("%t", data.EnableAutomaticPostmortemDraft.ValueBool()))
    }
    if !data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.IsNull() && !data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.IsUnknown() {
        filters["acknowledgeLinkedAlertsWhenIncidentAcknowledged"] = data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.ValueBool()
        filterNames = append(filterNames, "acknowledge_linked_alerts_when_incident_acknowledged = "+fmt.Sprintf("%t", data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.ValueBool()))
    }
    if !data.ResolveLinkedAlertsWhenIncidentResolved.IsNull() && !data.ResolveLinkedAlertsWhenIncidentResolved.IsUnknown() {
        filters["resolveLinkedAlertsWhenIncidentResolved"] = data.ResolveLinkedAlertsWhenIncidentResolved.ValueBool()
        filterNames = append(filterNames, "resolve_linked_alerts_when_incident_resolved = "+fmt.Sprintf("%t", data.ResolveLinkedAlertsWhenIncidentResolved.ValueBool()))
    }
    if !data.EnableIncidentInstrumentationFixTasks.IsNull() && !data.EnableIncidentInstrumentationFixTasks.IsUnknown() {
        filters["enableIncidentInstrumentationFixTasks"] = data.EnableIncidentInstrumentationFixTasks.ValueBool()
        filterNames = append(filterNames, "enable_incident_instrumentation_fix_tasks = "+fmt.Sprintf("%t", data.EnableIncidentInstrumentationFixTasks.ValueBool()))
    }
    if !data.EnableAlertInstrumentationFixTasks.IsNull() && !data.EnableAlertInstrumentationFixTasks.IsUnknown() {
        filters["enableAlertInstrumentationFixTasks"] = data.EnableAlertInstrumentationFixTasks.ValueBool()
        filterNames = append(filterNames, "enable_alert_instrumentation_fix_tasks = "+fmt.Sprintf("%t", data.EnableAlertInstrumentationFixTasks.ValueBool()))
    }
    if !data.EnableAutomaticIncidentCodeFixes.IsNull() && !data.EnableAutomaticIncidentCodeFixes.IsUnknown() {
        filters["enableAutomaticIncidentCodeFixes"] = data.EnableAutomaticIncidentCodeFixes.ValueBool()
        filterNames = append(filterNames, "enable_automatic_incident_code_fixes = "+fmt.Sprintf("%t", data.EnableAutomaticIncidentCodeFixes.ValueBool()))
    }
    if !data.EnableAutomaticAlertCodeFixes.IsNull() && !data.EnableAutomaticAlertCodeFixes.IsUnknown() {
        filters["enableAutomaticAlertCodeFixes"] = data.EnableAutomaticAlertCodeFixes.ValueBool()
        filterNames = append(filterNames, "enable_automatic_alert_code_fixes = "+fmt.Sprintf("%t", data.EnableAutomaticAlertCodeFixes.ValueBool()))
    }
    if !data.EnableAiInsights.IsNull() && !data.EnableAiInsights.IsUnknown() {
        filters["enableAiInsights"] = data.EnableAiInsights.ValueBool()
        filterNames = append(filterNames, "enable_ai_insights = "+fmt.Sprintf("%t", data.EnableAiInsights.ValueBool()))
    }
    if !data.EnableInsightFixTasks.IsNull() && !data.EnableInsightFixTasks.IsUnknown() {
        filters["enableInsightFixTasks"] = data.EnableInsightFixTasks.ValueBool()
        filterNames = append(filterNames, "enable_insight_fix_tasks = "+fmt.Sprintf("%t", data.EnableInsightFixTasks.ValueBool()))
    }
    if !data.AutoArchiveNonActionableExceptions.IsNull() && !data.AutoArchiveNonActionableExceptions.IsUnknown() {
        filters["autoArchiveNonActionableExceptions"] = data.AutoArchiveNonActionableExceptions.ValueBool()
        filterNames = append(filterNames, "auto_archive_non_actionable_exceptions = "+fmt.Sprintf("%t", data.AutoArchiveNonActionableExceptions.ValueBool()))
    }
    if !data.AlertInvestigationMinimumSeverityId.IsNull() && !data.AlertInvestigationMinimumSeverityId.IsUnknown() {
        filters["alertInvestigationMinimumSeverityId"] = data.AlertInvestigationMinimumSeverityId.ValueString()
        filterNames = append(filterNames, "alert_investigation_minimum_severity_id = "+fmt.Sprintf("%q", data.AlertInvestigationMinimumSeverityId.ValueString()))
    }
    if !data.AiDailyAutonomousTokenLimit.IsNull() && !data.AiDailyAutonomousTokenLimit.IsUnknown() {
        filters["aiDailyAutonomousTokenLimit"] = lookupNumber(data.AiDailyAutonomousTokenLimit)
        filterNames = append(filterNames, "ai_daily_autonomous_token_limit = "+data.AiDailyAutonomousTokenLimit.ValueBigFloat().String())
    }
    if !data.IncidentAiDailyAutonomousTokenLimit.IsNull() && !data.IncidentAiDailyAutonomousTokenLimit.IsUnknown() {
        filters["incidentAiDailyAutonomousTokenLimit"] = lookupNumber(data.IncidentAiDailyAutonomousTokenLimit)
        filterNames = append(filterNames, "incident_ai_daily_autonomous_token_limit = "+data.IncidentAiDailyAutonomousTokenLimit.ValueBigFloat().String())
    }
    if !data.AlertAiDailyAutonomousTokenLimit.IsNull() && !data.AlertAiDailyAutonomousTokenLimit.IsUnknown() {
        filters["alertAiDailyAutonomousTokenLimit"] = lookupNumber(data.AlertAiDailyAutonomousTokenLimit)
        filterNames = append(filterNames, "alert_ai_daily_autonomous_token_limit = "+data.AlertAiDailyAutonomousTokenLimit.ValueBigFloat().String())
    }
    if !data.AiDailyFixTaskLimit.IsNull() && !data.AiDailyFixTaskLimit.IsUnknown() {
        filters["aiDailyFixTaskLimit"] = lookupNumber(data.AiDailyFixTaskLimit)
        filterNames = append(filterNames, "ai_daily_fix_task_limit = "+data.AiDailyFixTaskLimit.ValueBigFloat().String())
    }
    if !data.IncidentAiDailyFixTaskLimit.IsNull() && !data.IncidentAiDailyFixTaskLimit.IsUnknown() {
        filters["incidentAiDailyFixTaskLimit"] = lookupNumber(data.IncidentAiDailyFixTaskLimit)
        filterNames = append(filterNames, "incident_ai_daily_fix_task_limit = "+data.IncidentAiDailyFixTaskLimit.ValueBigFloat().String())
    }
    if !data.AlertAiDailyFixTaskLimit.IsNull() && !data.AlertAiDailyFixTaskLimit.IsUnknown() {
        filters["alertAiDailyFixTaskLimit"] = lookupNumber(data.AlertAiDailyFixTaskLimit)
        filterNames = append(filterNames, "alert_ai_daily_fix_task_limit = "+data.AlertAiDailyFixTaskLimit.ValueBigFloat().String())
    }
    if !data.AlertInvestigationDedupeWindowMinutes.IsNull() && !data.AlertInvestigationDedupeWindowMinutes.IsUnknown() {
        filters["alertInvestigationDedupeWindowMinutes"] = lookupNumber(data.AlertInvestigationDedupeWindowMinutes)
        filterNames = append(filterNames, "alert_investigation_dedupe_window_minutes = "+data.AlertInvestigationDedupeWindowMinutes.ValueBigFloat().String())
    }
    if !data.IncidentInvestigationMinimumSeverityId.IsNull() && !data.IncidentInvestigationMinimumSeverityId.IsUnknown() {
        filters["incidentInvestigationMinimumSeverityId"] = data.IncidentInvestigationMinimumSeverityId.ValueString()
        filterNames = append(filterNames, "incident_investigation_minimum_severity_id = "+fmt.Sprintf("%q", data.IncidentInvestigationMinimumSeverityId.ValueString()))
    }
    if !data.IncidentInvestigationDedupeWindowMinutes.IsNull() && !data.IncidentInvestigationDedupeWindowMinutes.IsUnknown() {
        filters["incidentInvestigationDedupeWindowMinutes"] = lookupNumber(data.IncidentInvestigationDedupeWindowMinutes)
        filterNames = append(filterNames, "incident_investigation_dedupe_window_minutes = "+data.IncidentInvestigationDedupeWindowMinutes.ValueBigFloat().String())
    }
    if !data.AiMaxConcurrentInvestigations.IsNull() && !data.AiMaxConcurrentInvestigations.IsUnknown() {
        filters["aiMaxConcurrentInvestigations"] = lookupNumber(data.AiMaxConcurrentInvestigations)
        filterNames = append(filterNames, "ai_max_concurrent_investigations = "+data.AiMaxConcurrentInvestigations.ValueBigFloat().String())
    }
    if !data.IncidentAiMaxConcurrentInvestigations.IsNull() && !data.IncidentAiMaxConcurrentInvestigations.IsUnknown() {
        filters["incidentAiMaxConcurrentInvestigations"] = lookupNumber(data.IncidentAiMaxConcurrentInvestigations)
        filterNames = append(filterNames, "incident_ai_max_concurrent_investigations = "+data.IncidentAiMaxConcurrentInvestigations.ValueBigFloat().String())
    }
    if !data.AlertAiMaxConcurrentInvestigations.IsNull() && !data.AlertAiMaxConcurrentInvestigations.IsUnknown() {
        filters["alertAiMaxConcurrentInvestigations"] = lookupNumber(data.AlertAiMaxConcurrentInvestigations)
        filterNames = append(filterNames, "alert_ai_max_concurrent_investigations = "+data.AlertAiMaxConcurrentInvestigations.ValueBigFloat().String())
    }
    if !data.IncidentAiInvestigationTimeLimitInMinutes.IsNull() && !data.IncidentAiInvestigationTimeLimitInMinutes.IsUnknown() {
        filters["incidentAiInvestigationTimeLimitInMinutes"] = lookupNumber(data.IncidentAiInvestigationTimeLimitInMinutes)
        filterNames = append(filterNames, "incident_ai_investigation_time_limit_in_minutes = "+data.IncidentAiInvestigationTimeLimitInMinutes.ValueBigFloat().String())
    }
    if !data.AlertAiInvestigationTimeLimitInMinutes.IsNull() && !data.AlertAiInvestigationTimeLimitInMinutes.IsUnknown() {
        filters["alertAiInvestigationTimeLimitInMinutes"] = lookupNumber(data.AlertAiInvestigationTimeLimitInMinutes)
        filterNames = append(filterNames, "alert_ai_investigation_time_limit_in_minutes = "+data.AlertAiInvestigationTimeLimitInMinutes.ValueBigFloat().String())
    }
    if !data.EnableAutoRechargeAiBalance.IsNull() && !data.EnableAutoRechargeAiBalance.IsUnknown() {
        filters["enableAutoRechargeAiBalance"] = data.EnableAutoRechargeAiBalance.ValueBool()
        filterNames = append(filterNames, "enable_auto_recharge_ai_balance = "+fmt.Sprintf("%t", data.EnableAutoRechargeAiBalance.ValueBool()))
    }
    if !data.SendInvoicesByEmail.IsNull() && !data.SendInvoicesByEmail.IsUnknown() {
        filters["sendInvoicesByEmail"] = data.SendInvoicesByEmail.ValueBool()
        filterNames = append(filterNames, "send_invoices_by_email = "+fmt.Sprintf("%t", data.SendInvoicesByEmail.ValueBool()))
    }
    if !data.PlanName.IsNull() && !data.PlanName.IsUnknown() {
        filters["planName"] = data.PlanName.ValueString()
        filterNames = append(filterNames, "plan_name = "+fmt.Sprintf("%q", data.PlanName.ValueString()))
    }
    if !data.DataResidency.IsNull() && !data.DataResidency.IsUnknown() {
        filters["dataResidency"] = data.DataResidency.ValueString()
        filterNames = append(filterNames, "data_residency = "+fmt.Sprintf("%q", data.DataResidency.ValueString()))
    }
    if !data.ResellerId.IsNull() && !data.ResellerId.IsUnknown() {
        filters["resellerId"] = data.ResellerId.ValueString()
        filterNames = append(filterNames, "reseller_id = "+fmt.Sprintf("%q", data.ResellerId.ValueString()))
    }
    if !data.ResellerPlanId.IsNull() && !data.ResellerPlanId.IsUnknown() {
        filters["resellerPlanId"] = data.ResellerPlanId.ValueString()
        filterNames = append(filterNames, "reseller_plan_id = "+fmt.Sprintf("%q", data.ResellerPlanId.ValueString()))
    }
    if !data.LetCustomerSupportAccessProject.IsNull() && !data.LetCustomerSupportAccessProject.IsUnknown() {
        filters["letCustomerSupportAccessProject"] = data.LetCustomerSupportAccessProject.ValueBool()
        filterNames = append(filterNames, "let_customer_support_access_project = "+fmt.Sprintf("%t", data.LetCustomerSupportAccessProject.ValueBool()))
    }
    if !data.DoNotAddGlobalProbesByDefaultOnNewMonitors.IsNull() && !data.DoNotAddGlobalProbesByDefaultOnNewMonitors.IsUnknown() {
        filters["doNotAddGlobalProbesByDefaultOnNewMonitors"] = data.DoNotAddGlobalProbesByDefaultOnNewMonitors.ValueBool()
        filterNames = append(filterNames, "do_not_add_global_probes_by_default_on_new_monitors = "+fmt.Sprintf("%t", data.DoNotAddGlobalProbesByDefaultOnNewMonitors.ValueBool()))
    }
    if !data.GitHubAppInstallationId.IsNull() && !data.GitHubAppInstallationId.IsUnknown() {
        filters["gitHubAppInstallationId"] = data.GitHubAppInstallationId.ValueString()
        filterNames = append(filterNames, "git_hub_app_installation_id = "+fmt.Sprintf("%q", data.GitHubAppInstallationId.ValueString()))
    }
    if !data.DefaultMetricCardinalityBudget.IsNull() && !data.DefaultMetricCardinalityBudget.IsUnknown() {
        filters["defaultMetricCardinalityBudget"] = lookupNumber(data.DefaultMetricCardinalityBudget)
        filterNames = append(filterNames, "default_metric_cardinality_budget = "+data.DefaultMetricCardinalityBudget.ValueBigFloat().String())
    }
    if !data.DefaultTelemetryRetentionInDays.IsNull() && !data.DefaultTelemetryRetentionInDays.IsUnknown() {
        filters["defaultTelemetryRetentionInDays"] = lookupNumber(data.DefaultTelemetryRetentionInDays)
        filterNames = append(filterNames, "default_telemetry_retention_in_days = "+data.DefaultTelemetryRetentionInDays.ValueBigFloat().String())
    }
    if !data.EnableAuditLogs.IsNull() && !data.EnableAuditLogs.IsUnknown() {
        filters["enableAuditLogs"] = data.EnableAuditLogs.ValueBool()
        filterNames = append(filterNames, "enable_audit_logs = "+fmt.Sprintf("%t", data.EnableAuditLogs.ValueBool()))
    }
    if !data.IsSessionReplayAllowed.IsNull() && !data.IsSessionReplayAllowed.IsUnknown() {
        filters["isSessionReplayAllowed"] = data.IsSessionReplayAllowed.ValueBool()
        filterNames = append(filterNames, "is_session_replay_allowed = "+fmt.Sprintf("%t", data.IsSessionReplayAllowed.ValueBool()))
    }
    if !data.AuditLogsRetentionInDays.IsNull() && !data.AuditLogsRetentionInDays.IsUnknown() {
        filters["auditLogsRetentionInDays"] = lookupNumber(data.AuditLogsRetentionInDays)
        filterNames = append(filterNames, "audit_logs_retention_in_days = "+data.AuditLogsRetentionInDays.ValueBigFloat().String())
    }
    if !data.StoreSystemEventsInAuditLogs.IsNull() && !data.StoreSystemEventsInAuditLogs.IsUnknown() {
        filters["storeSystemEventsInAuditLogs"] = data.StoreSystemEventsInAuditLogs.ValueBool()
        filterNames = append(filterNames, "store_system_events_in_audit_logs = "+fmt.Sprintf("%t", data.StoreSystemEventsInAuditLogs.ValueBool()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the project up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the project up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "name": true,
        "slug": true,
        "paymentProviderPlanId": true,
        "paymentProviderSubscriptionId": true,
        "paymentProviderMeteredSubscriptionId": true,
        "paymentProviderSubscriptionSeats": true,
        "trialEndsAt": true,
        "paymentProviderCustomerId": true,
        "businessDetails": true,
        "businessDetailsCountry": true,
        "financeAccountingEmail": true,
        "paymentProviderSubscriptionStatus": true,
        "paymentProviderMeteredSubscriptionStatus": true,
        "paymentProviderPromoCode": true,
        "createdByUserId": true,
        "isFeatureFlagMonitorGroupsEnabled": true,
        "workflowRunsInLast30Days": true,
        "requireSsoForLogin": true,
        "requireSsoWithSsoProviderId": true,
        "incidentNumberPrefix": true,
        "alertNumberPrefix": true,
        "scheduledMaintenanceNumberPrefix": true,
        "incidentEpisodeNumberPrefix": true,
        "alertEpisodeNumberPrefix": true,
        "smsOrCallCurrentBalanceInUSDCents": true,
        "autoRechargeSmsOrCallByBalanceInUSD": true,
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": true,
        "enableSmsNotifications": true,
        "enableWhatsAppNotifications": true,
        "enableTelegramNotifications": true,
        "enableCallNotifications": true,
        "disableOnCallNotificationFallback": true,
        "enableAutoRechargeSmsOrCallBalance": true,
        "aiCurrentBalanceInUSDCents": true,
        "autoAiRechargeByBalanceInUSD": true,
        "autoRechargeAiWhenCurrentBalanceFallsInUSD": true,
        "enableAi": true,
        "aiDailyTokenLimit": true,
        "aiDailySpendLimitInUSD": true,
        "enableAutomaticIncidentInvestigation": true,
        "enableAutomaticAlertInvestigation": true,
        "enableAutomaticIncidentRemediation": true,
        "enableAutomaticAlertRemediation": true,
        "enableAutomaticPostmortemDraft": true,
        "acknowledgeLinkedAlertsWhenIncidentAcknowledged": true,
        "resolveLinkedAlertsWhenIncidentResolved": true,
        "enableIncidentInstrumentationFixTasks": true,
        "enableAlertInstrumentationFixTasks": true,
        "enableAutomaticIncidentCodeFixes": true,
        "enableAutomaticAlertCodeFixes": true,
        "enableAiInsights": true,
        "enableInsightFixTasks": true,
        "autoArchiveNonActionableExceptions": true,
        "alertInvestigationMinimumSeverityId": true,
        "aiDailyAutonomousTokenLimit": true,
        "incidentAiDailyAutonomousTokenLimit": true,
        "alertAiDailyAutonomousTokenLimit": true,
        "aiDailyFixTaskLimit": true,
        "incidentAiDailyFixTaskLimit": true,
        "alertAiDailyFixTaskLimit": true,
        "alertInvestigationDedupeWindowMinutes": true,
        "incidentInvestigationMinimumSeverityId": true,
        "incidentInvestigationDedupeWindowMinutes": true,
        "aiMaxConcurrentInvestigations": true,
        "incidentAiMaxConcurrentInvestigations": true,
        "alertAiMaxConcurrentInvestigations": true,
        "incidentAiInvestigationTimeLimitInMinutes": true,
        "alertAiInvestigationTimeLimitInMinutes": true,
        "enableAutoRechargeAiBalance": true,
        "sendInvoicesByEmail": true,
        "planName": true,
        "dataResidency": true,
        "resellerId": true,
        "resellerPlanId": true,
        "letCustomerSupportAccessProject": true,
        "doNotAddGlobalProbesByDefaultOnNewMonitors": true,
        "gitHubAppInstallationId": true,
        "defaultMetricCardinalityBudget": true,
        "defaultTelemetryRetentionInDays": true,
        "telemetryRetentionConfig": true,
        "defaultMetricDownsamplingRetentionDays": true,
        "enableAuditLogs": true,
        "isSessionReplayAllowed": true,
        "auditLogsRetentionInDays": true,
        "storeSystemEventsInAuditLogs": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/project/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No project found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read project: %s", err))
            return
        }
        if wrapper, ok := itemResponse["data"].(map[string]interface{}); ok {
            item = wrapper
        } else {
            item = itemResponse
        }
    }
    if !hasId {
        listBody := map[string]interface{}{
            "query":  filters,
            "select": selectParam,
            // limit 2 is enough to detect ambiguity without paging.
            "limit": 2,
        }
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/project/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list project, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list project: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No project matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one project matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for project.")
            return
        }
        item = first
    }

    // Update the model with response data
    if obj, ok := item["_id"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := item["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := item["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := item["createdAt"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := item["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := item["updatedAt"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := item["name"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := item["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := item["slug"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := item["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := item["paymentProviderPlanId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PaymentProviderPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderPlanId = types.StringNull()
        }
    } else if val, ok := item["paymentProviderPlanId"].(string); ok {
        data.PaymentProviderPlanId = types.StringValue(val)
    } else {
        data.PaymentProviderPlanId = types.StringNull()
    }
    if obj, ok := item["paymentProviderSubscriptionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PaymentProviderSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderSubscriptionId = types.StringNull()
        }
    } else if val, ok := item["paymentProviderSubscriptionId"].(string); ok {
        data.PaymentProviderSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionId = types.StringNull()
    }
    if obj, ok := item["paymentProviderMeteredSubscriptionId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderMeteredSubscriptionId = types.StringNull()
        }
    } else if val, ok := item["paymentProviderMeteredSubscriptionId"].(string); ok {
        data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionId = types.StringNull()
    }
    if val, ok := item["paymentProviderSubscriptionSeats"].(float64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["paymentProviderSubscriptionSeats"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
        } else {
            data.PaymentProviderSubscriptionSeats = types.NumberNull()
        }
    } else {
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if obj, ok := item["trialEndsAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrialEndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TrialEndsAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TrialEndsAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TrialEndsAt = types.StringValue(string(jsonBytes))
        } else {
            data.TrialEndsAt = types.StringNull()
        }
    } else if val, ok := item["trialEndsAt"].(string); ok {
        data.TrialEndsAt = types.StringValue(val)
    } else {
        data.TrialEndsAt = types.StringNull()
    }
    if obj, ok := item["paymentProviderCustomerId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderCustomerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PaymentProviderCustomerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderCustomerId = types.StringNull()
        }
    } else if val, ok := item["paymentProviderCustomerId"].(string); ok {
        data.PaymentProviderCustomerId = types.StringValue(val)
    } else {
        data.PaymentProviderCustomerId = types.StringNull()
    }
    if obj, ok := item["businessDetails"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BusinessDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BusinessDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BusinessDetails = types.StringValue(string(jsonBytes))
        } else {
            data.BusinessDetails = types.StringNull()
        }
    } else if val, ok := item["businessDetails"].(string); ok {
        data.BusinessDetails = types.StringValue(val)
    } else {
        data.BusinessDetails = types.StringNull()
    }
    if obj, ok := item["businessDetailsCountry"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BusinessDetailsCountry = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.BusinessDetailsCountry = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
        } else {
            data.BusinessDetailsCountry = types.StringNull()
        }
    } else if val, ok := item["businessDetailsCountry"].(string); ok {
        data.BusinessDetailsCountry = types.StringValue(val)
    } else {
        data.BusinessDetailsCountry = types.StringNull()
    }
    if obj, ok := item["financeAccountingEmail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FinanceAccountingEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.FinanceAccountingEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
        } else {
            data.FinanceAccountingEmail = types.StringNull()
        }
    } else if val, ok := item["financeAccountingEmail"].(string); ok {
        data.FinanceAccountingEmail = types.StringValue(val)
    } else {
        data.FinanceAccountingEmail = types.StringNull()
    }
    if obj, ok := item["paymentProviderSubscriptionStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PaymentProviderSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderSubscriptionStatus = types.StringNull()
        }
    } else if val, ok := item["paymentProviderSubscriptionStatus"].(string); ok {
        data.PaymentProviderSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionStatus = types.StringNull()
    }
    if obj, ok := item["paymentProviderMeteredSubscriptionStatus"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
        }
    } else if val, ok := item["paymentProviderMeteredSubscriptionStatus"].(string); ok {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
    }
    if obj, ok := item["paymentProviderPromoCode"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderPromoCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PaymentProviderPromoCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderPromoCode = types.StringNull()
        }
    } else if val, ok := item["paymentProviderPromoCode"].(string); ok {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if obj, ok := item["createdByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := item["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := item["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    } else {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolNull()
    }
    if val, ok := item["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["workflowRunsInLast30Days"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
        } else {
            data.WorkflowRunsInLast30Days = types.NumberNull()
        }
    } else {
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := item["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    } else {
        data.RequireSsoForLogin = types.BoolNull()
    }
    if obj, ok := item["requireSsoWithSsoProviderId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequireSsoWithSsoProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.RequireSsoWithSsoProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.RequireSsoWithSsoProviderId = types.StringNull()
        }
    } else if val, ok := item["requireSsoWithSsoProviderId"].(string); ok {
        data.RequireSsoWithSsoProviderId = types.StringValue(val)
    } else {
        data.RequireSsoWithSsoProviderId = types.StringNull()
    }
    if obj, ok := item["incidentNumberPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentNumberPrefix = types.StringNull()
        }
    } else if val, ok := item["incidentNumberPrefix"].(string); ok {
        data.IncidentNumberPrefix = types.StringValue(val)
    } else {
        data.IncidentNumberPrefix = types.StringNull()
    }
    if obj, ok := item["alertNumberPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.AlertNumberPrefix = types.StringNull()
        }
    } else if val, ok := item["alertNumberPrefix"].(string); ok {
        data.AlertNumberPrefix = types.StringValue(val)
    } else {
        data.AlertNumberPrefix = types.StringNull()
    }
    if obj, ok := item["scheduledMaintenanceNumberPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceNumberPrefix = types.StringNull()
        }
    } else if val, ok := item["scheduledMaintenanceNumberPrefix"].(string); ok {
        data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceNumberPrefix = types.StringNull()
    }
    if obj, ok := item["incidentEpisodeNumberPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentEpisodeNumberPrefix = types.StringNull()
        }
    } else if val, ok := item["incidentEpisodeNumberPrefix"].(string); ok {
        data.IncidentEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.IncidentEpisodeNumberPrefix = types.StringNull()
    }
    if obj, ok := item["alertEpisodeNumberPrefix"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.AlertEpisodeNumberPrefix = types.StringNull()
        }
    } else if val, ok := item["alertEpisodeNumberPrefix"].(string); ok {
        data.AlertEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.AlertEpisodeNumberPrefix = types.StringNull()
    }
    if val, ok := item["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["smsOrCallCurrentBalanceInUSDCents"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
        }
    } else {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
    }
    if val, ok := item["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["autoRechargeSmsOrCallByBalanceInUSD"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
        }
    } else {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
    }
    if val, ok := item["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberNull()
        }
    } else {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := item["enableSmsNotifications"].(bool); ok {
        data.EnableSmsNotifications = types.BoolValue(val)
    } else {
        data.EnableSmsNotifications = types.BoolNull()
    }
    if val, ok := item["enableWhatsAppNotifications"].(bool); ok {
        data.EnableWhatsAppNotifications = types.BoolValue(val)
    } else {
        data.EnableWhatsAppNotifications = types.BoolNull()
    }
    if val, ok := item["enableTelegramNotifications"].(bool); ok {
        data.EnableTelegramNotifications = types.BoolValue(val)
    } else {
        data.EnableTelegramNotifications = types.BoolNull()
    }
    if val, ok := item["enableCallNotifications"].(bool); ok {
        data.EnableCallNotifications = types.BoolValue(val)
    } else {
        data.EnableCallNotifications = types.BoolNull()
    }
    if val, ok := item["disableOnCallNotificationFallback"].(bool); ok {
        data.DisableOnCallNotificationFallback = types.BoolValue(val)
    } else {
        data.DisableOnCallNotificationFallback = types.BoolNull()
    }
    if val, ok := item["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    } else {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolNull()
    }
    if val, ok := item["aiCurrentBalanceInUSDCents"].(float64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["aiCurrentBalanceInUSDCents"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiCurrentBalanceInUsdCents = types.NumberNull()
        }
    } else {
        data.AiCurrentBalanceInUsdCents = types.NumberNull()
    }
    if val, ok := item["autoAiRechargeByBalanceInUSD"].(float64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["autoAiRechargeByBalanceInUSD"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
        }
    } else {
        data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
    }
    if val, ok := item["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
        }
    } else {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := item["enableAi"].(bool); ok {
        data.EnableAi = types.BoolValue(val)
    } else {
        data.EnableAi = types.BoolNull()
    }
    if val, ok := item["aiDailyTokenLimit"].(float64); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["aiDailyTokenLimit"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyTokenLimit = types.NumberNull()
        }
    } else {
        data.AiDailyTokenLimit = types.NumberNull()
    }
    if val, ok := item["aiDailySpendLimitInUSD"].(float64); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["aiDailySpendLimitInUSD"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailySpendLimitInUsd = types.NumberNull()
        }
    } else {
        data.AiDailySpendLimitInUsd = types.NumberNull()
    }
    if val, ok := item["enableAutomaticIncidentInvestigation"].(bool); ok {
        data.EnableAutomaticIncidentInvestigation = types.BoolValue(val)
    } else {
        data.EnableAutomaticIncidentInvestigation = types.BoolNull()
    }
    if val, ok := item["enableAutomaticAlertInvestigation"].(bool); ok {
        data.EnableAutomaticAlertInvestigation = types.BoolValue(val)
    } else {
        data.EnableAutomaticAlertInvestigation = types.BoolNull()
    }
    if val, ok := item["enableAutomaticIncidentRemediation"].(bool); ok {
        data.EnableAutomaticIncidentRemediation = types.BoolValue(val)
    } else {
        data.EnableAutomaticIncidentRemediation = types.BoolNull()
    }
    if val, ok := item["enableAutomaticAlertRemediation"].(bool); ok {
        data.EnableAutomaticAlertRemediation = types.BoolValue(val)
    } else {
        data.EnableAutomaticAlertRemediation = types.BoolNull()
    }
    if val, ok := item["enableAutomaticPostmortemDraft"].(bool); ok {
        data.EnableAutomaticPostmortemDraft = types.BoolValue(val)
    } else {
        data.EnableAutomaticPostmortemDraft = types.BoolNull()
    }
    if val, ok := item["acknowledgeLinkedAlertsWhenIncidentAcknowledged"].(bool); ok {
        data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged = types.BoolValue(val)
    } else {
        data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged = types.BoolNull()
    }
    if val, ok := item["resolveLinkedAlertsWhenIncidentResolved"].(bool); ok {
        data.ResolveLinkedAlertsWhenIncidentResolved = types.BoolValue(val)
    } else {
        data.ResolveLinkedAlertsWhenIncidentResolved = types.BoolNull()
    }
    if val, ok := item["enableIncidentInstrumentationFixTasks"].(bool); ok {
        data.EnableIncidentInstrumentationFixTasks = types.BoolValue(val)
    } else {
        data.EnableIncidentInstrumentationFixTasks = types.BoolNull()
    }
    if val, ok := item["enableAlertInstrumentationFixTasks"].(bool); ok {
        data.EnableAlertInstrumentationFixTasks = types.BoolValue(val)
    } else {
        data.EnableAlertInstrumentationFixTasks = types.BoolNull()
    }
    if val, ok := item["enableAutomaticIncidentCodeFixes"].(bool); ok {
        data.EnableAutomaticIncidentCodeFixes = types.BoolValue(val)
    } else {
        data.EnableAutomaticIncidentCodeFixes = types.BoolNull()
    }
    if val, ok := item["enableAutomaticAlertCodeFixes"].(bool); ok {
        data.EnableAutomaticAlertCodeFixes = types.BoolValue(val)
    } else {
        data.EnableAutomaticAlertCodeFixes = types.BoolNull()
    }
    if val, ok := item["enableAiInsights"].(bool); ok {
        data.EnableAiInsights = types.BoolValue(val)
    } else {
        data.EnableAiInsights = types.BoolNull()
    }
    if val, ok := item["enableInsightFixTasks"].(bool); ok {
        data.EnableInsightFixTasks = types.BoolValue(val)
    } else {
        data.EnableInsightFixTasks = types.BoolNull()
    }
    if val, ok := item["autoArchiveNonActionableExceptions"].(bool); ok {
        data.AutoArchiveNonActionableExceptions = types.BoolValue(val)
    } else {
        data.AutoArchiveNonActionableExceptions = types.BoolNull()
    }
    if obj, ok := item["alertInvestigationMinimumSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertInvestigationMinimumSeverityId = types.StringNull()
        }
    } else if val, ok := item["alertInvestigationMinimumSeverityId"].(string); ok {
        data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
    } else {
        data.AlertInvestigationMinimumSeverityId = types.StringNull()
    }
    if val, ok := item["aiDailyAutonomousTokenLimit"].(float64); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["aiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        data.AiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := item["incidentAiDailyAutonomousTokenLimit"].(float64); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["incidentAiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := item["alertAiDailyAutonomousTokenLimit"].(float64); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["alertAiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := item["aiDailyFixTaskLimit"].(float64); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["aiDailyFixTaskLimit"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        data.AiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := item["incidentAiDailyFixTaskLimit"].(float64); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["incidentAiDailyFixTaskLimit"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        data.IncidentAiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := item["alertAiDailyFixTaskLimit"].(float64); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["alertAiDailyFixTaskLimit"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        data.AlertAiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := item["alertInvestigationDedupeWindowMinutes"].(float64); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["alertInvestigationDedupeWindowMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertInvestigationDedupeWindowMinutes = types.NumberNull()
        }
    } else {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberNull()
    }
    if obj, ok := item["incidentInvestigationMinimumSeverityId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentInvestigationMinimumSeverityId = types.StringNull()
        }
    } else if val, ok := item["incidentInvestigationMinimumSeverityId"].(string); ok {
        data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
    } else {
        data.IncidentInvestigationMinimumSeverityId = types.StringNull()
    }
    if val, ok := item["incidentInvestigationDedupeWindowMinutes"].(float64); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["incidentInvestigationDedupeWindowMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentInvestigationDedupeWindowMinutes = types.NumberNull()
        }
    } else {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberNull()
    }
    if val, ok := item["aiMaxConcurrentInvestigations"].(float64); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["aiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        data.AiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := item["incidentAiMaxConcurrentInvestigations"].(float64); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["incidentAiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := item["alertAiMaxConcurrentInvestigations"].(float64); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["alertAiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        data.AlertAiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := item["incidentAiInvestigationTimeLimitInMinutes"].(float64); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["incidentAiInvestigationTimeLimitInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberNull()
        }
    } else {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberNull()
    }
    if val, ok := item["alertAiInvestigationTimeLimitInMinutes"].(float64); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["alertAiInvestigationTimeLimitInMinutes"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiInvestigationTimeLimitInMinutes = types.NumberNull()
        }
    } else {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberNull()
    }
    if val, ok := item["enableAutoRechargeAiBalance"].(bool); ok {
        data.EnableAutoRechargeAiBalance = types.BoolValue(val)
    } else {
        data.EnableAutoRechargeAiBalance = types.BoolNull()
    }
    if val, ok := item["sendInvoicesByEmail"].(bool); ok {
        data.SendInvoicesByEmail = types.BoolValue(val)
    } else {
        data.SendInvoicesByEmail = types.BoolNull()
    }
    if obj, ok := item["planName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PlanName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.PlanName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.PlanName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.PlanName = types.StringValue(string(jsonBytes))
        } else {
            data.PlanName = types.StringNull()
        }
    } else if val, ok := item["planName"].(string); ok {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
    }
    if obj, ok := item["dataResidency"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataResidency = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DataResidency = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DataResidency = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DataResidency = types.StringValue(string(jsonBytes))
        } else {
            data.DataResidency = types.StringNull()
        }
    } else if val, ok := item["dataResidency"].(string); ok {
        data.DataResidency = types.StringValue(val)
    } else {
        data.DataResidency = types.StringNull()
    }
    if obj, ok := item["resellerId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResellerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResellerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResellerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResellerId = types.StringValue(string(jsonBytes))
        } else {
            data.ResellerId = types.StringNull()
        }
    } else if val, ok := item["resellerId"].(string); ok {
        data.ResellerId = types.StringValue(val)
    } else {
        data.ResellerId = types.StringNull()
    }
    if obj, ok := item["resellerPlanId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResellerPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ResellerPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ResellerPlanId = types.StringValue(string(jsonBytes))
        } else {
            data.ResellerPlanId = types.StringNull()
        }
    } else if val, ok := item["resellerPlanId"].(string); ok {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := item["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    } else {
        data.LetCustomerSupportAccessProject = types.BoolNull()
    }
    if val, ok := item["doNotAddGlobalProbesByDefaultOnNewMonitors"].(bool); ok {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = types.BoolValue(val)
    } else {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = types.BoolNull()
    }
    if obj, ok := item["gitHubAppInstallationId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
        } else {
            data.GitHubAppInstallationId = types.StringNull()
        }
    } else if val, ok := item["gitHubAppInstallationId"].(string); ok {
        data.GitHubAppInstallationId = types.StringValue(val)
    } else {
        data.GitHubAppInstallationId = types.StringNull()
    }
    if val, ok := item["defaultMetricCardinalityBudget"].(float64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["defaultMetricCardinalityBudget"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
        } else {
            data.DefaultMetricCardinalityBudget = types.NumberNull()
        }
    } else {
        data.DefaultMetricCardinalityBudget = types.NumberNull()
    }
    if val, ok := item["defaultTelemetryRetentionInDays"].(float64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["defaultTelemetryRetentionInDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.DefaultTelemetryRetentionInDays = types.NumberNull()
        }
    } else {
        data.DefaultTelemetryRetentionInDays = types.NumberNull()
    }
    if obj, ok := item["telemetryRetentionConfig"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.TelemetryRetentionConfig = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.TelemetryRetentionConfig = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.TelemetryRetentionConfig = types.StringValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = types.StringNull()
        }
    } else if val, ok := item["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = types.StringValue(val)
    } else {
        data.TelemetryRetentionConfig = types.StringNull()
    }
    if obj, ok := item["defaultMetricDownsamplingRetentionDays"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultMetricDownsamplingRetentionDays = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DefaultMetricDownsamplingRetentionDays = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DefaultMetricDownsamplingRetentionDays = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DefaultMetricDownsamplingRetentionDays = types.StringValue(string(jsonBytes))
        } else {
            data.DefaultMetricDownsamplingRetentionDays = types.StringNull()
        }
    } else if val, ok := item["defaultMetricDownsamplingRetentionDays"].(string); ok {
        data.DefaultMetricDownsamplingRetentionDays = types.StringValue(val)
    } else {
        data.DefaultMetricDownsamplingRetentionDays = types.StringNull()
    }
    if val, ok := item["enableAuditLogs"].(bool); ok {
        data.EnableAuditLogs = types.BoolValue(val)
    } else {
        data.EnableAuditLogs = types.BoolNull()
    }
    if val, ok := item["isSessionReplayAllowed"].(bool); ok {
        data.IsSessionReplayAllowed = types.BoolValue(val)
    } else {
        data.IsSessionReplayAllowed = types.BoolNull()
    }
    if val, ok := item["auditLogsRetentionInDays"].(float64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["auditLogsRetentionInDays"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.AuditLogsRetentionInDays = types.NumberNull()
        }
    } else {
        data.AuditLogsRetentionInDays = types.NumberNull()
    }
    if val, ok := item["storeSystemEventsInAuditLogs"].(bool); ok {
        data.StoreSystemEventsInAuditLogs = types.BoolValue(val)
    } else {
        data.StoreSystemEventsInAuditLogs = types.BoolNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
