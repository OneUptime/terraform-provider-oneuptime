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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
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
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
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
    EnableAutoRemediation types.Bool `tfsdk:"enable_auto_remediation"`
    EnableAiCommandExecution types.Bool `tfsdk:"enable_ai_command_execution"`
    EnableAutomaticIncidentInvestigation types.Bool `tfsdk:"enable_automatic_incident_investigation"`
    EnableAutomaticAlertInvestigation types.Bool `tfsdk:"enable_automatic_alert_investigation"`
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
    EnableAutoRechargeAiBalance types.Bool `tfsdk:"enable_auto_recharge_ai_balance"`
    SendInvoicesByEmail types.Bool `tfsdk:"send_invoices_by_email"`
    PlanName types.String `tfsdk:"plan_name"`
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
        MarkdownDescription: "OneUptime Project, and everything happens inside it Look up an existing project by `id` or by `name`.",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Look up by unique identifier. Exactly one of `id` or `name` must be set.",
                Optional: true,
                Computed: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.",
                Optional: true,
                Computed: true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "payment_provider_plan_id": schema.StringAttribute{
                Computed: true,
            },
            "payment_provider_subscription_id": schema.StringAttribute{
                Computed: true,
            },
            "payment_provider_metered_subscription_id": schema.StringAttribute{
                Computed: true,
            },
            "payment_provider_subscription_seats": schema.NumberAttribute{
                Computed: true,
            },
            "trial_ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "payment_provider_customer_id": schema.StringAttribute{
                Computed: true,
            },
            "business_details": schema.StringAttribute{
                MarkdownDescription: "Business legal name, address and any tax information to appear on invoices..",
                Computed: true,
            },
            "business_details_country": schema.StringAttribute{
                MarkdownDescription: "Two-letter ISO country code for billing address (e.g., US, GB, DE)..",
                Computed: true,
            },
            "finance_accounting_email": schema.StringAttribute{
                MarkdownDescription: "Invoices, receipts and billing related notifications will be sent to these emails in addition to project owner. Separate multiple emails with a comma..",
                Computed: true,
            },
            "payment_provider_subscription_status": schema.StringAttribute{
                Computed: true,
            },
            "payment_provider_metered_subscription_status": schema.StringAttribute{
                Computed: true,
            },
            "payment_provider_promo_code": schema.StringAttribute{
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_feature_flag_monitor_groups_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Feature Flag Monitor Groups Enabled.",
                Computed: true,
            },
            "workflow_runs_in_last30_days": schema.NumberAttribute{
                Computed: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                Computed: true,
            },
            "require_sso_with_sso_provider_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incident_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident numbers (e.g., 'INC-'). If empty, '#' is used..",
                Computed: true,
            },
            "alert_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert numbers (e.g., 'ALT-'). If empty, '#' is used..",
                Computed: true,
            },
            "scheduled_maintenance_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for scheduled maintenance numbers (e.g., 'SM-'). If empty, '#' is used..",
                Computed: true,
            },
            "incident_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident episode numbers (e.g., 'IE-'). If empty, '#' is used..",
                Computed: true,
            },
            "alert_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert episode numbers (e.g., 'AE-'). If empty, '#' is used..",
                Computed: true,
            },
            "sms_or_call_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for SMS, Call, and WhatsApp.",
                Computed: true,
            },
            "auto_recharge_sms_or_call_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for SMS, Call, and WhatsApp.",
                Computed: true,
            },
            "auto_recharge_sms_or_call_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for SMS, Call, and WhatsApp.",
                Computed: true,
            },
            "enable_sms_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable SMS notifications for this project..",
                Computed: true,
            },
            "enable_whats_app_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable WhatsApp notifications for this project..",
                Computed: true,
            },
            "enable_telegram_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable Telegram notifications for this project..",
                Computed: true,
            },
            "enable_call_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable call notifications for this project..",
                Computed: true,
            },
            "disable_on_call_notification_fallback": schema.BoolAttribute{
                MarkdownDescription: "When enabled, a page routed to a responder with no matching notification rule fails instead of falling back to their verified notification methods..",
                Computed: true,
            },
            "enable_auto_recharge_sms_or_call_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for SMS, Call, and WhatsApp balance for this project..",
                Computed: true,
            },
            "ai_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for AI services.",
                Computed: true,
            },
            "auto_ai_recharge_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for AI services.",
                Computed: true,
            },
            "auto_recharge_ai_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for AI services.",
                Computed: true,
            },
            "enable_ai": schema.BoolAttribute{
                MarkdownDescription: "Enable AI services for this project..",
                Computed: true,
            },
            "enable_auto_remediation": schema.BoolAttribute{
                MarkdownDescription: "Kill switch for auto-remediation: when disabled, no auto-remediation rule fires in this project..",
                Computed: true,
            },
            "enable_ai_command_execution": schema.BoolAttribute{
                MarkdownDescription: "When enabled, auto-remediation rules may let the AI compose and run commands on opted-in Runners (with an operator allowlist for auto-execution, and one-click approval for everything else). Off by default..",
                Computed: true,
            },
            "enable_automatic_incident_investigation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime's AI SRE automatically investigates every new incident and posts a cited root cause analysis to the incident timeline. Requires AI to be enabled and an LLM provider to be configured..",
                Computed: true,
            },
            "enable_automatic_alert_investigation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime's AI SRE automatically investigates every new alert and posts a cited root cause analysis to the alert timeline. Requires AI to be enabled and an LLM provider to be configured..",
                Computed: true,
            },
            "enable_incident_instrumentation_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an incident AI investigation that ends inconclusive (telemetry was insufficient to determine a root cause) automatically queues an AI agent task that opens a pull request adding the missing instrumentation to the implicated code paths. Requires a repository connected through the GitHub App. Pull requests are always human-reviewed — nothing merges automatically..",
                Computed: true,
            },
            "enable_alert_instrumentation_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an alert AI investigation that ends inconclusive (telemetry was insufficient to determine a root cause) automatically queues an AI agent task that opens a pull request adding the missing instrumentation to the implicated code paths. Requires a repository connected through the GitHub App. Pull requests are always human-reviewed — nothing merges automatically..",
                Computed: true,
            },
            "enable_automatic_incident_code_fixes": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an incident AI investigation that ends with a confident, evidenced root cause analysis and recommends a repository code change automatically queues an AI agent task that opens a fix pull request, ready for review, from that analysis — the automatic form of the 'Open Fix PR from this analysis' button. Operational, infrastructure, external, user-error and inconclusive findings do not offer or open code-fix pull requests. Requires a repository connected through the GitHub App and a Runner with the code-fix capability. Pull requests are always human-reviewed — nothing merges automatically..",
                Computed: true,
            },
            "enable_automatic_alert_code_fixes": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an alert AI investigation that ends with a confident, evidenced root cause analysis and recommends a repository code change automatically queues an AI agent task that opens a fix pull request, ready for review, from that analysis — the automatic form of the 'Open Fix PR from this analysis' button. Operational, infrastructure, external, user-error and inconclusive findings do not offer or open code-fix pull requests. Requires a repository connected through the GitHub App and a Runner with the code-fix capability. Pull requests are always human-reviewed — nothing merges automatically..",
                Computed: true,
            },
            "enable_ai_insights": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime AI continuously watches this project's telemetry with deterministic statistical sensors (error-log spikes, exception novelty and spikes, trace-latency regressions, week-over-week metric drift) and files quiet Insights — never pages, never opens incidents. Each new insight also gets a budgeted, read-only AI triage analysis when an LLM provider is configured..",
                Computed: true,
            },
            "enable_insight_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "When enabled, insights whose deterministic evidence points at code (new or spiking exceptions with a resolvable repository, trace-latency regressions with span-tree findings) automatically queue an AI agent task that opens a pull request with a proposed fix, ready for review. Honors the daily fix task budget and per-repository open-PR caps. Pull requests are always human-reviewed — nothing merges automatically..",
                Computed: true,
            },
            "auto_archive_non_actionable_exceptions": schema.BoolAttribute{
                MarkdownDescription: "When enabled, exception groups the AI triage classifies as expected denials (auth failures, plan/paywall rejections, scanner probes tripping intentional validation) are automatically archived so they stop surfacing in the unresolved list and never queue AI fix tasks. Groups classified as user errors or infrastructure conditions are NOT auto-archived — only clear expected denials are. Archiving is reversible from the Archived tab..",
                Computed: true,
            },
            "alert_investigation_minimum_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "ai_daily_autonomous_token_limit": schema.NumberAttribute{
                MarkdownDescription: "Fallback maximum tokens per UTC day for autonomous AI work that is not associated with an incident or alert. When the limit is reached, new autonomous work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit..",
                Computed: true,
            },
            "incident_ai_daily_autonomous_token_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum tokens per UTC day that autonomous incident-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new incident-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit..",
                Computed: true,
            },
            "alert_ai_daily_autonomous_token_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum tokens per UTC day that autonomous alert-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new alert-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit..",
                Computed: true,
            },
            "ai_daily_fix_task_limit": schema.NumberAttribute{
                MarkdownDescription: "Fallback maximum AI fix tasks (agent runs that open pull requests) that may be created per UTC day for work not associated with an incident or alert, across every fix recipe and trigger. Unset means the default of 25 per day; 0 pauses these AI fix tasks entirely..",
                Computed: true,
            },
            "incident_ai_daily_fix_task_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum AI fix tasks derived from incidents that may be created per UTC day for this project. Unset means the default of 25 per day; 0 pauses incident AI fix tasks entirely..",
                Computed: true,
            },
            "alert_ai_daily_fix_task_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum AI fix tasks derived from alerts that may be created per UTC day for this project. Unset means the default of 25 per day; 0 pauses alert AI fix tasks entirely..",
                Computed: true,
            },
            "alert_investigation_dedupe_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Repeat alerts from the same monitor within this many minutes are not re-investigated by AI — the first analysis stands. Unset means the default of 30 minutes; 0 disables the cooldown..",
                Computed: true,
            },
            "incident_investigation_minimum_severity_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "incident_investigation_dedupe_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Incidents affecting a monitor that AI investigated within this many minutes are not re-investigated — the first analysis stands. Unset means the default of 30 minutes; 0 disables the cooldown..",
                Computed: true,
            },
            "ai_max_concurrent_investigations": schema.NumberAttribute{
                MarkdownDescription: "Fallback maximum number of non-incident and non-alert AI investigations that may run at the same time for this project. Unset means the default of 3. Minimum 1 — pause autonomous work with its opt-in toggle or a daily token limit of 0 instead..",
                Computed: true,
            },
            "incident_ai_max_concurrent_investigations": schema.NumberAttribute{
                MarkdownDescription: "How many incident AI investigations may run at the same time for this project. Unset means the default of 3. Minimum 1 — pause incident investigations with the opt-in toggle or a daily token limit of 0 instead..",
                Computed: true,
            },
            "alert_ai_max_concurrent_investigations": schema.NumberAttribute{
                MarkdownDescription: "How many alert AI investigations may run at the same time for this project. Unset means the default of 3. Minimum 1 — pause alert investigations with the opt-in toggle or a daily token limit of 0 instead..",
                Computed: true,
            },
            "enable_auto_recharge_ai_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for AI balance for this project..",
                Computed: true,
            },
            "send_invoices_by_email": schema.BoolAttribute{
                MarkdownDescription: "When enabled, invoices will be automatically sent to the finance/accounting email when they are generated..",
                Computed: true,
            },
            "plan_name": schema.StringAttribute{
                MarkdownDescription: "Name of the plan this project is subscribed to..",
                Computed: true,
            },
            "reseller_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "reseller_plan_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "let_customer_support_access_project": schema.BoolAttribute{
                MarkdownDescription: "OneUptime customer support can access this project. This is used for debugging purposes..",
                Computed: true,
            },
            "do_not_add_global_probes_by_default_on_new_monitors": schema.BoolAttribute{
                MarkdownDescription: "If enabled, global probes will NOT be automatically added to new monitors. Enable this only if you are using ONLY custom probes to monitor your resources..",
                Computed: true,
            },
            "git_hub_app_installation_id": schema.StringAttribute{
                MarkdownDescription: "The GitHub App installation ID for this project. This is set when the GitHub App is installed on the organization..",
                Computed: true,
            },
            "default_metric_cardinality_budget": schema.NumberAttribute{
                MarkdownDescription: "Project-wide default max distinct series per metric. Services without a per-service override use this value..",
                Computed: true,
            },
            "default_telemetry_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "Project-wide default number of days to retain telemetry data (logs, traces, metrics). Services without a per-service override use this value..",
                Computed: true,
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Project-wide per-pillar retention overrides for telemetry data (logs by severity, traces by status, metrics, profiles). Falls back to defaultTelemetryRetentionInDays when a pillar or bucket is not set..",
                Computed: true,
            },
            "default_metric_downsampling_retention_days": schema.StringAttribute{
                MarkdownDescription: "Project-wide default retention for each downsampling tier (raw, 1m, 5m, 1h, 1d) in days..",
                Computed: true,
            },
            "enable_audit_logs": schema.BoolAttribute{
                MarkdownDescription: "When enabled, changes to resources in this project are recorded as audit log entries..",
                Computed: true,
            },
            "is_session_replay_allowed": schema.BoolAttribute{
                MarkdownDescription: "When enabled, RUM applications in this project may record session replays if they are individually enabled too. On by default; switch it off here to stop session replay across the entire project in one place..",
                Computed: true,
            },
            "audit_logs_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain audit log entries. Minimum 7, maximum 180..",
                Computed: true,
            },
            "store_system_events_in_audit_logs": schema.BoolAttribute{
                MarkdownDescription: "When enabled, audit logs will also include events triggered by the system. By default, only events triggered by users are recorded..",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a project.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
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
        "deletedByUserId": true,
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
        "enableAutoRemediation": true,
        "enableAiCommandExecution": true,
        "enableAutomaticIncidentInvestigation": true,
        "enableAutomaticAlertInvestigation": true,
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
        "enableAutoRechargeAiBalance": true,
        "sendInvoicesByEmail": true,
        "planName": true,
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No project found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one project matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
    if obj, ok := item["deletedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := item["deletedAt"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
    }
    if val, ok := item["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if obj, ok := item["version"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(float64); ok {
            data.Version = types.NumberValue(big.NewFloat(val))
        } else {
            data.Version = types.NumberNull()
        }
    } else {
        data.Version = types.NumberNull()
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
    if obj, ok := item["deletedByUserId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := item["deletedByUserId"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
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
    if val, ok := item["enableAutoRemediation"].(bool); ok {
        data.EnableAutoRemediation = types.BoolValue(val)
    } else {
        data.EnableAutoRemediation = types.BoolNull()
    }
    if val, ok := item["enableAiCommandExecution"].(bool); ok {
        data.EnableAiCommandExecution = types.BoolValue(val)
    } else {
        data.EnableAiCommandExecution = types.BoolNull()
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
