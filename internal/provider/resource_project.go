package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
    return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
    client *Client
}

// ProjectResourceModel describes the resource data model.
type ProjectResourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    PaymentProviderPlanId types.String `tfsdk:"payment_provider_plan_id"`
    BusinessDetails types.String `tfsdk:"business_details"`
    BusinessDetailsCountry types.String `tfsdk:"business_details_country"`
    FinanceAccountingEmail types.String `tfsdk:"finance_accounting_email"`
    PaymentProviderPromoCode types.String `tfsdk:"payment_provider_promo_code"`
    IsFeatureFlagMonitorGroupsEnabled types.Bool `tfsdk:"is_feature_flag_monitor_groups_enabled"`
    ActiveMonitorsLimit types.Number `tfsdk:"active_monitors_limit"`
    SeatLimit types.Number `tfsdk:"seat_limit"`
    IncidentNumberPrefix types.String `tfsdk:"incident_number_prefix"`
    AlertNumberPrefix types.String `tfsdk:"alert_number_prefix"`
    ScheduledMaintenanceNumberPrefix types.String `tfsdk:"scheduled_maintenance_number_prefix"`
    IncidentEpisodeNumberPrefix types.String `tfsdk:"incident_episode_number_prefix"`
    AlertEpisodeNumberPrefix types.String `tfsdk:"alert_episode_number_prefix"`
    EnableAutomaticIncidentInvestigation types.Bool `tfsdk:"enable_automatic_incident_investigation"`
    EnableAutomaticAlertInvestigation types.Bool `tfsdk:"enable_automatic_alert_investigation"`
    EnableAutomaticIncidentRemediation types.Bool `tfsdk:"enable_automatic_incident_remediation"`
    EnableAutomaticAlertRemediation types.Bool `tfsdk:"enable_automatic_alert_remediation"`
    EnableAutomaticPostmortemDraft types.Bool `tfsdk:"enable_automatic_postmortem_draft"`
    EnableIncidentInstrumentationFixTasks types.Bool `tfsdk:"enable_incident_instrumentation_fix_tasks"`
    EnableAlertInstrumentationFixTasks types.Bool `tfsdk:"enable_alert_instrumentation_fix_tasks"`
    EnableAutomaticIncidentCodeFixes types.Bool `tfsdk:"enable_automatic_incident_code_fixes"`
    EnableAutomaticAlertCodeFixes types.Bool `tfsdk:"enable_automatic_alert_code_fixes"`
    EnableAiInsights types.Bool `tfsdk:"enable_ai_insights"`
    EnableInsightFixTasks types.Bool `tfsdk:"enable_insight_fix_tasks"`
    AutoArchiveNonActionableExceptions types.Bool `tfsdk:"auto_archive_non_actionable_exceptions"`
    SendInvoicesByEmail types.Bool `tfsdk:"send_invoices_by_email"`
    UtmContent types.String `tfsdk:"utm_content"`
    EnableAuditLogs types.Bool `tfsdk:"enable_audit_logs"`
    IsSessionReplayAllowed types.Bool `tfsdk:"is_session_replay_allowed"`
    AuditLogsRetentionInDays types.Number `tfsdk:"audit_logs_retention_in_days"`
    StoreSystemEventsInAuditLogs types.Bool `tfsdk:"store_system_events_in_audit_logs"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    RequireSsoWithSsoProviderId types.String `tfsdk:"require_sso_with_sso_provider_id"`
    AutoRechargeSmsOrCallByBalanceInUsd types.Number `tfsdk:"auto_recharge_sms_or_call_by_balance_in_usd"`
    AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd types.Number `tfsdk:"auto_recharge_sms_or_call_when_current_balance_falls_in_usd"`
    EnableSmsNotifications types.Bool `tfsdk:"enable_sms_notifications"`
    EnableWhatsAppNotifications types.Bool `tfsdk:"enable_whats_app_notifications"`
    EnableTelegramNotifications types.Bool `tfsdk:"enable_telegram_notifications"`
    EnableCallNotifications types.Bool `tfsdk:"enable_call_notifications"`
    DisableOnCallNotificationFallback types.Bool `tfsdk:"disable_on_call_notification_fallback"`
    EnableAutoRechargeSmsOrCallBalance types.Bool `tfsdk:"enable_auto_recharge_sms_or_call_balance"`
    AutoAiRechargeByBalanceInUsd types.Number `tfsdk:"auto_ai_recharge_by_balance_in_usd"`
    AutoRechargeAiWhenCurrentBalanceFallsInUsd types.Number `tfsdk:"auto_recharge_ai_when_current_balance_falls_in_usd"`
    EnableAi types.Bool `tfsdk:"enable_ai"`
    AiDailyTokenLimit types.Number `tfsdk:"ai_daily_token_limit"`
    AiDailySpendLimitInUsd types.Number `tfsdk:"ai_daily_spend_limit_in_usd"`
    AcknowledgeLinkedAlertsWhenIncidentAcknowledged types.Bool `tfsdk:"acknowledge_linked_alerts_when_incident_acknowledged"`
    ResolveLinkedAlertsWhenIncidentResolved types.Bool `tfsdk:"resolve_linked_alerts_when_incident_resolved"`
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
    DoNotAddGlobalProbesByDefaultOnNewMonitors types.Bool `tfsdk:"do_not_add_global_probes_by_default_on_new_monitors"`
    DefaultMetricCardinalityBudget types.Number `tfsdk:"default_metric_cardinality_budget"`
    DefaultTelemetryRetentionInDays types.Number `tfsdk:"default_telemetry_retention_in_days"`
    TelemetryRetentionConfig JSONSubsetValue `tfsdk:"telemetry_retention_config"`
    DefaultMetricDownsamplingRetentionDays JSONSubsetValue `tfsdk:"default_metric_downsampling_retention_days"`
    CreatedAt RFC3339Value `tfsdk:"created_at"`
    UpdatedAt RFC3339Value `tfsdk:"updated_at"`
    Slug types.String `tfsdk:"slug"`
    PaymentProviderSubscriptionId types.String `tfsdk:"payment_provider_subscription_id"`
    PaymentProviderMeteredSubscriptionId types.String `tfsdk:"payment_provider_metered_subscription_id"`
    PaymentProviderSubscriptionSeats types.Number `tfsdk:"payment_provider_subscription_seats"`
    TrialEndsAt RFC3339Value `tfsdk:"trial_ends_at"`
    PaymentProviderCustomerId types.String `tfsdk:"payment_provider_customer_id"`
    PaymentProviderSubscriptionStatus types.String `tfsdk:"payment_provider_subscription_status"`
    PaymentProviderMeteredSubscriptionStatus types.String `tfsdk:"payment_provider_metered_subscription_status"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    WorkflowRunsInLast30Days types.Number `tfsdk:"workflow_runs_in_last30_days"`
    SmsOrCallCurrentBalanceInUsdCents types.Number `tfsdk:"sms_or_call_current_balance_in_usd_cents"`
    AiCurrentBalanceInUsdCents types.Number `tfsdk:"ai_current_balance_in_usd_cents"`
    PlanName types.String `tfsdk:"plan_name"`
    DataResidency types.String `tfsdk:"data_residency"`
    ResellerId types.String `tfsdk:"reseller_id"`
    ResellerPlanId types.String `tfsdk:"reseller_plan_id"`
    LetCustomerSupportAccessProject types.Bool `tfsdk:"let_customer_support_access_project"`
    GitHubAppInstallationId types.String `tfsdk:"git_hub_app_installation_id"`
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = r.schemaDefinition()
}

func (r *ProjectResource) schemaDefinition() schema.Schema {
    return schema.Schema{
        MarkdownDescription: "OneUptime Project, and everything happens inside it",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object.",
                Required: true,
            },
            "payment_provider_plan_id": schema.StringAttribute{
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "business_details": schema.StringAttribute{
                MarkdownDescription: "Business legal name, address and any tax information to appear on invoices.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "business_details_country": schema.StringAttribute{
                MarkdownDescription: "Two-letter ISO country code for billing address (e.g., US, GB, DE).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "finance_accounting_email": schema.StringAttribute{
                MarkdownDescription: "Invoices, receipts and billing related notifications will be sent to these emails in addition to project owner. Separate multiple emails with a comma.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "payment_provider_promo_code": schema.StringAttribute{
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "is_feature_flag_monitor_groups_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Feature Flag Monitor Groups Enabled.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "active_monitors_limit": schema.NumberAttribute{
                Optional: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.RequiresReplace(),
                },
            },
            "seat_limit": schema.NumberAttribute{
                Optional: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.RequiresReplace(),
                },
            },
            "incident_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident numbers (e.g., 'INC-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber incidents that already exist.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert numbers (e.g., 'ALT-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber alerts that already exist.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "scheduled_maintenance_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for scheduled maintenance numbers (e.g., 'SM-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber scheduled maintenance events that already exist.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident episode numbers (e.g., 'IE-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber incident episodes that already exist.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert episode numbers (e.g., 'AE-'). If empty, '#' is used. Up to 20 letters, numbers or - _ . / : #, not ending in a digit. Changing it does not renumber alert episodes that already exist.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_automatic_incident_investigation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime's AI SRE automatically investigates every new incident and posts a cited root cause analysis to the incident timeline; any auto-remediation for the incident waits until that investigation settles. On for new projects created in OneUptime; projects that existed before keep their setting. Drafting a postmortem when an incident resolves is a separate setting (Enable Automatic Postmortem Draft). Requires AI to be enabled and an LLM provider to be configured.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_automatic_alert_investigation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime's AI SRE automatically investigates every new alert and posts a cited root cause analysis to the alert timeline. On for new projects created in OneUptime; projects that existed before keep their setting. Requires AI to be enabled and an LLM provider to be configured.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_automatic_incident_remediation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime fixes new incidents automatically: OneUptime AI fixes each one on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, and Auto Remediation Rules - when there are any - choose which incidents are fixed, which runbooks run and whether a person approves first. Off by default, for new projects too. Any AI investigation of the incident settles first. Requires AI to be enabled. It also holds the incident pull-request switches, enableAutomaticIncidentCodeFixes and enableIncidentInstrumentationFixTasks: they open pull requests only while this is on. The dashboard turns them on and off with it; through the API, set them in the same request.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_automatic_alert_remediation": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime fixes new alerts automatically: OneUptime AI fixes each one on the Kubernetes clusters and infrastructure it is linked to, the way each one's AI agent page allows, and Auto Remediation Rules - when there are any - choose which alerts are fixed, which runbooks run and whether a person approves first. Off by default, for new projects too. Any AI investigation of the alert settles first. Requires AI to be enabled. It also holds the alert pull-request switches, enableAutomaticAlertCodeFixes and enableAlertInstrumentationFixTasks: they open pull requests only while this is on. The dashboard turns them on and off with it; through the API, set them in the same request.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_automatic_postmortem_draft": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime's AI SRE drafts a postmortem from the incident's timeline and telemetry when an incident is resolved, for a human to review and edit. It never overwrites a postmortem that already exists. On for new projects created in OneUptime; projects that existed before keep their setting. Requires AI to be enabled and an LLM provider to be configured.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_incident_instrumentation_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an incident AI investigation that ends inconclusive (telemetry was insufficient to determine a root cause) automatically queues an AI agent task that opens a pull request adding the missing instrumentation to the implicated code paths. Requires a repository connected through the GitHub App. Pull requests are always human-reviewed — nothing merges automatically. Part of fixing: it acts only while enableAutomaticIncidentRemediation (Fix new incidents automatically) is on, and the dashboard turns it on and off with that switch. Off for new projects.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_alert_instrumentation_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an alert AI investigation that ends inconclusive (telemetry was insufficient to determine a root cause) automatically queues an AI agent task that opens a pull request adding the missing instrumentation to the implicated code paths. Requires a repository connected through the GitHub App. Pull requests are always human-reviewed — nothing merges automatically. Part of fixing: it acts only while enableAutomaticAlertRemediation (Fix new alerts automatically) is on, and the dashboard turns it on and off with that switch. Off for new projects.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_automatic_incident_code_fixes": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an incident AI investigation that ends with a confident, evidenced root cause analysis and recommends a repository code change automatically queues an AI agent task that opens a fix pull request, ready for review, from that analysis — the automatic form of the 'Open Fix PR from this analysis' button. Operational, infrastructure, external, user-error and inconclusive findings do not offer or open code-fix pull requests. Requires a repository connected through the GitHub App and a Runner with the code-fix capability. Pull requests are always human-reviewed — nothing merges automatically. Part of fixing: it acts only while enableAutomaticIncidentRemediation (Fix new incidents automatically) is on, and the dashboard turns it on and off with that switch. Off for new projects.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_automatic_alert_code_fixes": schema.BoolAttribute{
                MarkdownDescription: "When enabled, an alert AI investigation that ends with a confident, evidenced root cause analysis and recommends a repository code change automatically queues an AI agent task that opens a fix pull request, ready for review, from that analysis — the automatic form of the 'Open Fix PR from this analysis' button. Operational, infrastructure, external, user-error and inconclusive findings do not offer or open code-fix pull requests. Requires a repository connected through the GitHub App and a Runner with the code-fix capability. Pull requests are always human-reviewed — nothing merges automatically. Part of fixing: it acts only while enableAutomaticAlertRemediation (Fix new alerts automatically) is on, and the dashboard turns it on and off with that switch. Off for new projects.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_ai_insights": schema.BoolAttribute{
                MarkdownDescription: "When enabled, OneUptime AI continuously watches this project's telemetry with deterministic statistical sensors (error-log spikes, exception novelty and spikes, trace-latency regressions, week-over-week metric drift) and files quiet Insights — never pages, never opens incidents. Each new insight also gets a budgeted, read-only AI triage analysis when an LLM provider is configured. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_insight_fix_tasks": schema.BoolAttribute{
                MarkdownDescription: "When enabled, insights whose deterministic evidence points at code (new or spiking exceptions with a resolvable repository, trace-latency regressions with span-tree findings) automatically queue an AI agent task that opens a pull request with a proposed fix, ready for review. Honors any open-PR cap set on the repository. Pull requests are always human-reviewed — nothing merges automatically. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_archive_non_actionable_exceptions": schema.BoolAttribute{
                MarkdownDescription: "When enabled, exception groups the AI triage classifies as expected denials (auth failures, plan/paywall rejections, scanner probes tripping intentional validation) are automatically archived so they stop surfacing in the unresolved list and never queue AI fix tasks. Groups classified as user errors or infrastructure conditions are NOT auto-archived — only clear expected denials are. Archiving is reversible from the Archived tab. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "send_invoices_by_email": schema.BoolAttribute{
                MarkdownDescription: "When enabled, invoices will be automatically sent to the finance/accounting email when they are generated.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "utm_content": schema.StringAttribute{
                Optional: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "enable_audit_logs": schema.BoolAttribute{
                MarkdownDescription: "When enabled, changes to resources in this project are recorded as audit log entries.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "is_session_replay_allowed": schema.BoolAttribute{
                MarkdownDescription: "When enabled, RUM applications in this project may record session replays if they are individually enabled too. On by default; switch it off here to stop session replay across the entire project in one place.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "audit_logs_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain audit log entries. Minimum 7, maximum 180.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(7)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "store_system_events_in_audit_logs": schema.BoolAttribute{
                MarkdownDescription: "When enabled, audit logs will also include events triggered by the system. By default, only events triggered by users are recorded.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "require_sso_for_login": schema.BoolAttribute{
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "require_sso_with_sso_provider_id": schema.StringAttribute{
                MarkdownDescription: "If set, SSO-enforced login for this project is only satisfied by an SSO token issued by this specific provider id (a Project SSO/OIDC or a Global SSO/OIDC). When null, any trusted SSO provider satisfies enforcement.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_recharge_sms_or_call_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for SMS, Call, and WhatsApp.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(20)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_recharge_sms_or_call_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for SMS, Call, and WhatsApp.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(10)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_sms_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable SMS notifications for this project.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_whats_app_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable WhatsApp notifications for this project.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_telegram_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable Telegram notifications for this project.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_call_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable call notifications for this project.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "disable_on_call_notification_fallback": schema.BoolAttribute{
                MarkdownDescription: "When enabled, a page routed to a responder with no matching notification rule fails instead of falling back to their verified notification methods.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_auto_recharge_sms_or_call_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for SMS, Call, and WhatsApp balance for this project.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_ai_recharge_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for AI services.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(20)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_recharge_ai_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for AI services.",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(10)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_ai": schema.BoolAttribute{
                MarkdownDescription: "Master switch for AI in this project. When disabled, every AI feature stops: Ask AI, investigations, postmortem drafts, auto-remediation and AI commands on Runners.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "ai_daily_token_limit": schema.NumberAttribute{
                MarkdownDescription: "The most tokens OneUptime AI may use in this project each UTC day, across every AI feature: Ask AI, investigations, postmortem drafts, fix pull requests, insight triage, workflows, runbooks and Slack or Microsoft Teams questions. Once it is reached, new AI work is refused until midnight UTC. The incident and alert daily limits still apply under it. Unset means no limit; a limit is a whole number of at least 1 (to turn AI off, use Enable AI).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "ai_daily_spend_limit_in_usd": schema.NumberAttribute{
                MarkdownDescription: "OneUptime Cloud: the most AI credits, in whole US dollars, OneUptime AI may spend in this project each UTC day. Only calls billed to the project's AI credits count, so it never stops AI that runs on the project's own LLM provider. Once it is reached, billed AI work is refused until midnight UTC. Ignored where AI is not billed (self-hosted). Unset means no limit; a limit is at least 1 (to turn AI off, use Enable AI).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "acknowledge_linked_alerts_when_incident_acknowledged": schema.BoolAttribute{
                MarkdownDescription: "When enabled, acknowledging an incident also acknowledges every alert linked to it. This stops those alerts' on-call escalations, and their reminders only when the alert reminder rule is set to stop on Acknowledged. Alerts linked to an incident that is already acknowledged are acknowledged as they are linked. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "resolve_linked_alerts_when_incident_resolved": schema.BoolAttribute{
                MarkdownDescription: "When enabled, resolving an incident also resolves every alert linked to it, except alerts that are still linked to another incident that is not resolved yet. Alerts linked to an incident that is already resolved are resolved as they are linked. On for new projects created in OneUptime; projects that existed before keep their setting.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_investigation_minimum_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the minimum AlertSeverity that triggers automatic investigation. The ID of a `oneuptime_alert_severity`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "ai_daily_autonomous_token_limit": schema.NumberAttribute{
                MarkdownDescription: "Legacy setting, no longer enforced: autonomous AI work that is not associated with an incident or alert has no daily token limit. Use the Daily Incident AI Token Limit and Daily Alert AI Token Limit instead.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_ai_daily_autonomous_token_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum tokens per UTC day that autonomous incident-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new incident-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_ai_daily_autonomous_token_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum tokens per UTC day that autonomous alert-linked AI work may consume for this project, including investigations, remediation, and follow-up fix tasks. When the limit is reached, new alert-linked AI work is skipped until the next day — interactive AI chat is never blocked. Unset means no limit.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "ai_daily_fix_task_limit": schema.NumberAttribute{
                MarkdownDescription: "Legacy setting, no longer enforced: AI fix tasks that are not associated with an incident or alert have no daily limit. Use the Daily Incident AI Fix Task Limit and Daily Alert AI Fix Task Limit instead.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_ai_daily_fix_task_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum AI fix tasks derived from incidents that may be created per UTC day for this project. Unset means no limit; 0 pauses incident AI fix tasks entirely.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_ai_daily_fix_task_limit": schema.NumberAttribute{
                MarkdownDescription: "Maximum AI fix tasks derived from alerts that may be created per UTC day for this project. Unset means no limit; 0 pauses alert AI fix tasks entirely.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_investigation_dedupe_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Repeat alerts from the same monitor within this many minutes are not re-investigated by AI — the first analysis stands. Unset or 0 means no cooldown, so every alert is investigated; at most 1440 minutes (a day).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_investigation_minimum_severity_id": schema.StringAttribute{
                MarkdownDescription: "ID of the minimum incident severity that is investigated automatically by AI. The ID of a `oneuptime_incident_severity`.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_investigation_dedupe_window_minutes": schema.NumberAttribute{
                MarkdownDescription: "Incidents affecting a monitor that AI investigated within this many minutes are not re-investigated — the first analysis stands. Unset or 0 means no cooldown, so every incident is investigated; at most 1440 minutes (a day).",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "ai_max_concurrent_investigations": schema.NumberAttribute{
                MarkdownDescription: "Legacy setting, no longer enforced: AI investigations that are not associated with an incident or alert have no concurrency limit. Use the Max Concurrent Incident Investigations and Max Concurrent Alert Investigations instead.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_ai_max_concurrent_investigations": schema.NumberAttribute{
                MarkdownDescription: "How many incident AI investigations may run at the same time for this project. Unset means no limit — every incident investigation starts right away. Minimum 1 — pause incident investigations with the Enable Automatic Incident Investigation toggle or a daily token limit of 0 instead.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_ai_max_concurrent_investigations": schema.NumberAttribute{
                MarkdownDescription: "How many alert AI investigations may run at the same time for this project. Unset means no limit — every alert investigation starts right away. Minimum 1 — pause alert investigations with the Enable Automatic Alert Investigation toggle or a daily token limit of 0 instead.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_ai_investigation_time_limit_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Stop an incident AI investigation after this many minutes and report what it found. Unset (the default) means no time limit — the investigation runs until it is done.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_ai_investigation_time_limit_in_minutes": schema.NumberAttribute{
                MarkdownDescription: "Stop an alert AI investigation after this many minutes and report what it found. Unset (the default) means no time limit — the investigation runs until it is done.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_auto_recharge_ai_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for AI balance for this project.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "do_not_add_global_probes_by_default_on_new_monitors": schema.BoolAttribute{
                MarkdownDescription: "If enabled, global probes will NOT be automatically added to new monitors. Enable this only if you are using ONLY custom probes to monitor your resources.",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "default_metric_cardinality_budget": schema.NumberAttribute{
                MarkdownDescription: "Project-wide default max distinct series per metric. Services without a per-service override use this value.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "default_telemetry_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "Project-wide default number of days to retain telemetry data (logs, traces, metrics). Services without a per-service override use this value.",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Project-wide per-pillar retention overrides for telemetry data (logs by severity, traces by status, metrics, profiles). Falls back to defaultTelemetryRetentionInDays when a pillar or bucket is not set. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "default_metric_downsampling_retention_days": schema.StringAttribute{
                MarkdownDescription: "Project-wide default retention for each downsampling tier (raw, 1m, 5m, 1h, 1d) in days. A JSON value: write it with `jsonencode()`.",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
                Validators: []validator.String{
                    JSONEnvelopeValidator(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was created.",
                CustomType: RFC3339Type{},
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "Date and Time when the object was updated.",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object.",
                Computed: true,
            },
            "payment_provider_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_metered_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_subscription_seats": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "trial_ends_at": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                CustomType: RFC3339Type{},
                Computed: true,
            },
            "payment_provider_customer_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_metered_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "workflow_runs_in_last30_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User, Billing Admin, Billing Member, Billing Viewer], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "sms_or_call_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for SMS, Call, and WhatsApp.",
                Computed: true,
            },
            "ai_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for AI services.",
                Computed: true,
            },
            "plan_name": schema.StringAttribute{
                MarkdownDescription: "Name of the plan this project is subscribed to.",
                Computed: true,
            },
            "data_residency": schema.StringAttribute{
                MarkdownDescription: "Where this project's data is hosted. Set by OneUptime staff on OneUptime Cloud.",
                Computed: true,
            },
            "reseller_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Reseller in which this object belongs.",
                Computed: true,
            },
            "reseller_plan_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Reseller Plan in which this object belongs.",
                Computed: true,
            },
            "let_customer_support_access_project": schema.BoolAttribute{
                MarkdownDescription: "OneUptime customer support can access this project. This is used for debugging purposes.",
                Computed: true,
            },
            "git_hub_app_installation_id": schema.StringAttribute{
                MarkdownDescription: "The GitHub App installation ID for this project. This is set when the GitHub App is installed on the organization.",
                Computed: true,
            },
        },
    }
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data ProjectResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config ProjectResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data



    // Create API request body. Unset (null/unknown) optional fields are
    // omitted so server-side defaults apply instead of being overwritten
    // with zero values.
    projectRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := projectRequest["data"].(map[string]interface{})

    if !data.Name.IsNull() && !data.Name.IsUnknown() {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.PaymentProviderPlanId.IsNull() && !data.PaymentProviderPlanId.IsUnknown() {
        requestDataMap["paymentProviderPlanId"] = data.PaymentProviderPlanId.ValueString()
    }
    if !data.BusinessDetails.IsNull() && !data.BusinessDetails.IsUnknown() {
        requestDataMap["businessDetails"] = data.BusinessDetails.ValueString()
    }
    if !data.BusinessDetailsCountry.IsNull() && !data.BusinessDetailsCountry.IsUnknown() {
        requestDataMap["businessDetailsCountry"] = data.BusinessDetailsCountry.ValueString()
    }
    if !data.FinanceAccountingEmail.IsNull() && !data.FinanceAccountingEmail.IsUnknown() {
        requestDataMap["financeAccountingEmail"] = data.FinanceAccountingEmail.ValueString()
    }
    if !data.PaymentProviderPromoCode.IsNull() && !data.PaymentProviderPromoCode.IsUnknown() {
        requestDataMap["paymentProviderPromoCode"] = data.PaymentProviderPromoCode.ValueString()
    }
    if !data.IsFeatureFlagMonitorGroupsEnabled.IsNull() && !data.IsFeatureFlagMonitorGroupsEnabled.IsUnknown() {
        requestDataMap["isFeatureFlagMonitorGroupsEnabled"] = data.IsFeatureFlagMonitorGroupsEnabled.ValueBool()
    }
    if !data.ActiveMonitorsLimit.IsNull() && !data.ActiveMonitorsLimit.IsUnknown() {
        requestDataMap["activeMonitorsLimit"] = r.bigFloatToFloat64(data.ActiveMonitorsLimit.ValueBigFloat())
    }
    if !data.SeatLimit.IsNull() && !data.SeatLimit.IsUnknown() {
        requestDataMap["seatLimit"] = r.bigFloatToFloat64(data.SeatLimit.ValueBigFloat())
    }
    if !data.IncidentNumberPrefix.IsNull() && !data.IncidentNumberPrefix.IsUnknown() {
        requestDataMap["incidentNumberPrefix"] = data.IncidentNumberPrefix.ValueString()
    }
    if !data.AlertNumberPrefix.IsNull() && !data.AlertNumberPrefix.IsUnknown() {
        requestDataMap["alertNumberPrefix"] = data.AlertNumberPrefix.ValueString()
    }
    if !data.ScheduledMaintenanceNumberPrefix.IsNull() && !data.ScheduledMaintenanceNumberPrefix.IsUnknown() {
        requestDataMap["scheduledMaintenanceNumberPrefix"] = data.ScheduledMaintenanceNumberPrefix.ValueString()
    }
    if !data.IncidentEpisodeNumberPrefix.IsNull() && !data.IncidentEpisodeNumberPrefix.IsUnknown() {
        requestDataMap["incidentEpisodeNumberPrefix"] = data.IncidentEpisodeNumberPrefix.ValueString()
    }
    if !data.AlertEpisodeNumberPrefix.IsNull() && !data.AlertEpisodeNumberPrefix.IsUnknown() {
        requestDataMap["alertEpisodeNumberPrefix"] = data.AlertEpisodeNumberPrefix.ValueString()
    }
    if !data.EnableAutomaticIncidentInvestigation.IsNull() && !data.EnableAutomaticIncidentInvestigation.IsUnknown() {
        requestDataMap["enableAutomaticIncidentInvestigation"] = data.EnableAutomaticIncidentInvestigation.ValueBool()
    }
    if !data.EnableAutomaticAlertInvestigation.IsNull() && !data.EnableAutomaticAlertInvestigation.IsUnknown() {
        requestDataMap["enableAutomaticAlertInvestigation"] = data.EnableAutomaticAlertInvestigation.ValueBool()
    }
    if !data.EnableAutomaticIncidentRemediation.IsNull() && !data.EnableAutomaticIncidentRemediation.IsUnknown() {
        requestDataMap["enableAutomaticIncidentRemediation"] = data.EnableAutomaticIncidentRemediation.ValueBool()
    }
    if !data.EnableAutomaticAlertRemediation.IsNull() && !data.EnableAutomaticAlertRemediation.IsUnknown() {
        requestDataMap["enableAutomaticAlertRemediation"] = data.EnableAutomaticAlertRemediation.ValueBool()
    }
    if !data.EnableAutomaticPostmortemDraft.IsNull() && !data.EnableAutomaticPostmortemDraft.IsUnknown() {
        requestDataMap["enableAutomaticPostmortemDraft"] = data.EnableAutomaticPostmortemDraft.ValueBool()
    }
    if !data.EnableIncidentInstrumentationFixTasks.IsNull() && !data.EnableIncidentInstrumentationFixTasks.IsUnknown() {
        requestDataMap["enableIncidentInstrumentationFixTasks"] = data.EnableIncidentInstrumentationFixTasks.ValueBool()
    }
    if !data.EnableAlertInstrumentationFixTasks.IsNull() && !data.EnableAlertInstrumentationFixTasks.IsUnknown() {
        requestDataMap["enableAlertInstrumentationFixTasks"] = data.EnableAlertInstrumentationFixTasks.ValueBool()
    }
    if !data.EnableAutomaticIncidentCodeFixes.IsNull() && !data.EnableAutomaticIncidentCodeFixes.IsUnknown() {
        requestDataMap["enableAutomaticIncidentCodeFixes"] = data.EnableAutomaticIncidentCodeFixes.ValueBool()
    }
    if !data.EnableAutomaticAlertCodeFixes.IsNull() && !data.EnableAutomaticAlertCodeFixes.IsUnknown() {
        requestDataMap["enableAutomaticAlertCodeFixes"] = data.EnableAutomaticAlertCodeFixes.ValueBool()
    }
    if !data.EnableAiInsights.IsNull() && !data.EnableAiInsights.IsUnknown() {
        requestDataMap["enableAiInsights"] = data.EnableAiInsights.ValueBool()
    }
    if !data.EnableInsightFixTasks.IsNull() && !data.EnableInsightFixTasks.IsUnknown() {
        requestDataMap["enableInsightFixTasks"] = data.EnableInsightFixTasks.ValueBool()
    }
    if !data.AutoArchiveNonActionableExceptions.IsNull() && !data.AutoArchiveNonActionableExceptions.IsUnknown() {
        requestDataMap["autoArchiveNonActionableExceptions"] = data.AutoArchiveNonActionableExceptions.ValueBool()
    }
    if !data.SendInvoicesByEmail.IsNull() && !data.SendInvoicesByEmail.IsUnknown() {
        requestDataMap["sendInvoicesByEmail"] = data.SendInvoicesByEmail.ValueBool()
    }
    if !data.UtmContent.IsNull() && !data.UtmContent.IsUnknown() {
        requestDataMap["utmContent"] = data.UtmContent.ValueString()
    }
    if !data.EnableAuditLogs.IsNull() && !data.EnableAuditLogs.IsUnknown() {
        requestDataMap["enableAuditLogs"] = data.EnableAuditLogs.ValueBool()
    }
    if !data.IsSessionReplayAllowed.IsNull() && !data.IsSessionReplayAllowed.IsUnknown() {
        requestDataMap["isSessionReplayAllowed"] = data.IsSessionReplayAllowed.ValueBool()
    }
    if !data.AuditLogsRetentionInDays.IsNull() && !data.AuditLogsRetentionInDays.IsUnknown() {
        requestDataMap["auditLogsRetentionInDays"] = r.bigFloatToFloat64(data.AuditLogsRetentionInDays.ValueBigFloat())
    }
    if !data.StoreSystemEventsInAuditLogs.IsNull() && !data.StoreSystemEventsInAuditLogs.IsUnknown() {
        requestDataMap["storeSystemEventsInAuditLogs"] = data.StoreSystemEventsInAuditLogs.ValueBool()
    }

    // Make API call
    httpResp, err := r.client.Post(ctx, "/project", projectRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create project, got error: %s", err))
        return
    }

    var projectResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &projectResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to create project: %s", err))
        return
    }

    // Extract the new resource id from the create response.
    createdId := ""
    if wrapper, ok := projectResponse["data"].(map[string]interface{}); ok {
        if val, ok := wrapper["_id"].(string); ok {
            createdId = val
        }
    } else if val, ok := projectResponse["_id"].(string); ok {
        createdId = val
    }
    if createdId == "" {
        resp.Diagnostics.AddError("OneUptime API Error", "Create response for project did not contain an id. This is a bug in the provider or the API; please report it.")
        return
    }
    data.Id = types.StringValue(createdId)

    /*
     * The server has committed the row. Persist what we know to state BEFORE
     * the read-back: if the read-back fails and we return without setting
     * state, Terraform never learns the resource exists and the created
     * project is orphaned server-side — never refreshed, never
     * destroyed. Delete already refuses to drop state on failure for the
     * same reason; Create must not either.
     */
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Re-read the resource so state reflects server-normalized values.
    selectParam := map[string]interface{}{
        "name": true,
        "paymentProviderPlanId": true,
        "businessDetails": true,
        "businessDetailsCountry": true,
        "financeAccountingEmail": true,
        "paymentProviderPromoCode": true,
        "isFeatureFlagMonitorGroupsEnabled": true,
        "incidentNumberPrefix": true,
        "alertNumberPrefix": true,
        "scheduledMaintenanceNumberPrefix": true,
        "incidentEpisodeNumberPrefix": true,
        "alertEpisodeNumberPrefix": true,
        "enableAutomaticIncidentInvestigation": true,
        "enableAutomaticAlertInvestigation": true,
        "enableAutomaticIncidentRemediation": true,
        "enableAutomaticAlertRemediation": true,
        "enableAutomaticPostmortemDraft": true,
        "enableIncidentInstrumentationFixTasks": true,
        "enableAlertInstrumentationFixTasks": true,
        "enableAutomaticIncidentCodeFixes": true,
        "enableAutomaticAlertCodeFixes": true,
        "enableAiInsights": true,
        "enableInsightFixTasks": true,
        "autoArchiveNonActionableExceptions": true,
        "sendInvoicesByEmail": true,
        "enableAuditLogs": true,
        "isSessionReplayAllowed": true,
        "auditLogsRetentionInDays": true,
        "storeSystemEventsInAuditLogs": true,
        "requireSsoForLogin": true,
        "requireSsoWithSsoProviderId": true,
        "autoRechargeSmsOrCallByBalanceInUSD": true,
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": true,
        "enableSmsNotifications": true,
        "enableWhatsAppNotifications": true,
        "enableTelegramNotifications": true,
        "enableCallNotifications": true,
        "disableOnCallNotificationFallback": true,
        "enableAutoRechargeSmsOrCallBalance": true,
        "autoAiRechargeByBalanceInUSD": true,
        "autoRechargeAiWhenCurrentBalanceFallsInUSD": true,
        "enableAi": true,
        "aiDailyTokenLimit": true,
        "aiDailySpendLimitInUSD": true,
        "acknowledgeLinkedAlertsWhenIncidentAcknowledged": true,
        "resolveLinkedAlertsWhenIncidentResolved": true,
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
        "doNotAddGlobalProbesByDefaultOnNewMonitors": true,
        "defaultMetricCardinalityBudget": true,
        "defaultTelemetryRetentionInDays": true,
        "telemetryRetentionConfig": true,
        "defaultMetricDownsamplingRetentionDays": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "paymentProviderSubscriptionId": true,
        "paymentProviderMeteredSubscriptionId": true,
        "paymentProviderSubscriptionSeats": true,
        "trialEndsAt": true,
        "paymentProviderCustomerId": true,
        "paymentProviderSubscriptionStatus": true,
        "paymentProviderMeteredSubscriptionStatus": true,
        "createdByUserId": true,
        "workflowRunsInLast30Days": true,
        "smsOrCallCurrentBalanceInUSDCents": true,
        "aiCurrentBalanceInUSDCents": true,
        "planName": true,
        "dataResidency": true,
        "resellerId": true,
        "resellerPlanId": true,
        "letCustomerSupportAccessProject": true,
        "gitHubAppInstallationId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/project/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        /*
         * State already owns the id, so the resource is tracked and the next
         * refresh reconciles the remaining attributes. Warn rather than
         * error: erroring here would strand a real resource.
         */
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created project but could not read it back; state is incomplete until the next refresh: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddWarning("Read After Create Failed", fmt.Sprintf("Created project but could not parse the read-back response; state is incomplete until the next refresh: %s", err))
        return
    }

    // Update the model with the authoritative read response
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderPlanId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderPlanId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderPlanId"].(string); ok {
        data.PaymentProviderPlanId = types.StringValue(val)
    } else {
        data.PaymentProviderPlanId = types.StringNull()
    }
    if obj, ok := dataMap["businessDetails"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BusinessDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.BusinessDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.BusinessDetails = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.BusinessDetails = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.BusinessDetails = types.StringValue(string(jsonBytes))
        } else {
            data.BusinessDetails = types.StringNull()
        }
    } else if val, ok := dataMap["businessDetails"].(string); ok {
        data.BusinessDetails = types.StringValue(val)
    } else {
        data.BusinessDetails = types.StringNull()
    }
    if obj, ok := dataMap["businessDetailsCountry"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BusinessDetailsCountry = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.BusinessDetailsCountry = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
        } else {
            data.BusinessDetailsCountry = types.StringNull()
        }
    } else if val, ok := dataMap["businessDetailsCountry"].(string); ok {
        data.BusinessDetailsCountry = types.StringValue(val)
    } else {
        data.BusinessDetailsCountry = types.StringNull()
    }
    if obj, ok := dataMap["financeAccountingEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FinanceAccountingEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FinanceAccountingEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
            } else {
                data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
            } else {
                data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
        } else {
            data.FinanceAccountingEmail = types.StringNull()
        }
    } else if val, ok := dataMap["financeAccountingEmail"].(string); ok {
        data.FinanceAccountingEmail = types.StringValue(val)
    } else {
        data.FinanceAccountingEmail = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderPromoCode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderPromoCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderPromoCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderPromoCode = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderPromoCode"].(string); ok {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if val, ok := dataMap["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["incidentNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["incidentNumberPrefix"].(string); ok {
        data.IncidentNumberPrefix = types.StringValue(val)
    } else {
        data.IncidentNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["alertNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.AlertNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["alertNumberPrefix"].(string); ok {
        data.AlertNumberPrefix = types.StringValue(val)
    } else {
        data.AlertNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["scheduledMaintenanceNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceNumberPrefix"].(string); ok {
        data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["incidentEpisodeNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentEpisodeNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["incidentEpisodeNumberPrefix"].(string); ok {
        data.IncidentEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.IncidentEpisodeNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["alertEpisodeNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.AlertEpisodeNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["alertEpisodeNumberPrefix"].(string); ok {
        data.AlertEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.AlertEpisodeNumberPrefix = types.StringNull()
    }
    if val, ok := dataMap["enableAutomaticIncidentInvestigation"].(bool); ok {
        data.EnableAutomaticIncidentInvestigation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticAlertInvestigation"].(bool); ok {
        data.EnableAutomaticAlertInvestigation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticIncidentRemediation"].(bool); ok {
        data.EnableAutomaticIncidentRemediation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticAlertRemediation"].(bool); ok {
        data.EnableAutomaticAlertRemediation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticPostmortemDraft"].(bool); ok {
        data.EnableAutomaticPostmortemDraft = types.BoolValue(val)
    }
    if val, ok := dataMap["enableIncidentInstrumentationFixTasks"].(bool); ok {
        data.EnableIncidentInstrumentationFixTasks = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAlertInstrumentationFixTasks"].(bool); ok {
        data.EnableAlertInstrumentationFixTasks = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticIncidentCodeFixes"].(bool); ok {
        data.EnableAutomaticIncidentCodeFixes = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticAlertCodeFixes"].(bool); ok {
        data.EnableAutomaticAlertCodeFixes = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAiInsights"].(bool); ok {
        data.EnableAiInsights = types.BoolValue(val)
    }
    if val, ok := dataMap["enableInsightFixTasks"].(bool); ok {
        data.EnableInsightFixTasks = types.BoolValue(val)
    }
    if val, ok := dataMap["autoArchiveNonActionableExceptions"].(bool); ok {
        data.AutoArchiveNonActionableExceptions = types.BoolValue(val)
    }
    if val, ok := dataMap["sendInvoicesByEmail"].(bool); ok {
        data.SendInvoicesByEmail = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAuditLogs"].(bool); ok {
        data.EnableAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["isSessionReplayAllowed"].(bool); ok {
        data.IsSessionReplayAllowed = types.BoolValue(val)
    }
    if val, ok := dataMap["auditLogsRetentionInDays"].(float64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["auditLogsRetentionInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.AuditLogsRetentionInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AuditLogsRetentionInDays = types.NumberNull()
    }
    if val, ok := dataMap["storeSystemEventsInAuditLogs"].(bool); ok {
        data.StoreSystemEventsInAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if obj, ok := dataMap["requireSsoWithSsoProviderId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequireSsoWithSsoProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RequireSsoWithSsoProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.RequireSsoWithSsoProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["requireSsoWithSsoProviderId"].(string); ok {
        data.RequireSsoWithSsoProviderId = types.StringValue(val)
    } else {
        data.RequireSsoWithSsoProviderId = types.StringNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := dataMap["enableSmsNotifications"].(bool); ok {
        data.EnableSmsNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["enableWhatsAppNotifications"].(bool); ok {
        data.EnableWhatsAppNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["enableTelegramNotifications"].(bool); ok {
        data.EnableTelegramNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["enableCallNotifications"].(bool); ok {
        data.EnableCallNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["disableOnCallNotificationFallback"].(bool); ok {
        data.DisableOnCallNotificationFallback = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(float64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoAiRechargeByBalanceInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := dataMap["enableAi"].(bool); ok {
        data.EnableAi = types.BoolValue(val)
    }
    if val, ok := dataMap["aiDailyTokenLimit"].(float64); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailyTokenLimit"].(int); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailyTokenLimit"].(int64); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailyTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailyTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["aiDailySpendLimitInUSD"].(float64); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailySpendLimitInUSD"].(int); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailySpendLimitInUSD"].(int64); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailySpendLimitInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailySpendLimitInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailySpendLimitInUsd = types.NumberNull()
    }
    if val, ok := dataMap["acknowledgeLinkedAlertsWhenIncidentAcknowledged"].(bool); ok {
        data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged = types.BoolValue(val)
    }
    if val, ok := dataMap["resolveLinkedAlertsWhenIncidentResolved"].(bool); ok {
        data.ResolveLinkedAlertsWhenIncidentResolved = types.BoolValue(val)
    }
    if obj, ok := dataMap["alertInvestigationMinimumSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertInvestigationMinimumSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["alertInvestigationMinimumSeverityId"].(string); ok {
        data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
    } else {
        data.AlertInvestigationMinimumSeverityId = types.StringNull()
    }
    if val, ok := dataMap["aiDailyAutonomousTokenLimit"].(float64); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailyAutonomousTokenLimit"].(int); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailyAutonomousTokenLimit"].(int64); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(float64); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(int); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(int64); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(float64); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(int); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(int64); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["aiDailyFixTaskLimit"].(float64); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailyFixTaskLimit"].(int); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailyFixTaskLimit"].(int64); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailyFixTaskLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiDailyFixTaskLimit"].(float64); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiDailyFixTaskLimit"].(int); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiDailyFixTaskLimit"].(int64); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiDailyFixTaskLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := dataMap["alertAiDailyFixTaskLimit"].(float64); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiDailyFixTaskLimit"].(int); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiDailyFixTaskLimit"].(int64); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiDailyFixTaskLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(float64); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(int); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(int64); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertInvestigationDedupeWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertInvestigationDedupeWindowMinutes = types.NumberNull()
    }
    if obj, ok := dataMap["incidentInvestigationMinimumSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentInvestigationMinimumSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentInvestigationMinimumSeverityId"].(string); ok {
        data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
    } else {
        data.IncidentInvestigationMinimumSeverityId = types.StringNull()
    }
    if val, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(float64); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(int); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(int64); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentInvestigationDedupeWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberNull()
    }
    if val, ok := dataMap["aiMaxConcurrentInvestigations"].(float64); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiMaxConcurrentInvestigations"].(int); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiMaxConcurrentInvestigations"].(int64); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(float64); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(int); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(int64); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := dataMap["alertAiMaxConcurrentInvestigations"].(float64); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiMaxConcurrentInvestigations"].(int); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiMaxConcurrentInvestigations"].(int64); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(float64); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(int); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(int64); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(float64); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(int); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(int64); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiInvestigationTimeLimitInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["enableAutoRechargeAiBalance"].(bool); ok {
        data.EnableAutoRechargeAiBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["doNotAddGlobalProbesByDefaultOnNewMonitors"].(bool); ok {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = types.BoolValue(val)
    }
    if val, ok := dataMap["defaultMetricCardinalityBudget"].(float64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["defaultMetricCardinalityBudget"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
        } else {
            data.DefaultMetricCardinalityBudget = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.DefaultMetricCardinalityBudget = types.NumberNull()
    }
    if val, ok := dataMap["defaultTelemetryRetentionInDays"].(float64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["defaultTelemetryRetentionInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.DefaultTelemetryRetentionInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.DefaultTelemetryRetentionInDays = types.NumberNull()
    }
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(string); ok {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
    } else {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderSubscriptionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderSubscriptionId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderSubscriptionId"].(string); ok {
        data.PaymentProviderSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionId = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderMeteredSubscriptionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderMeteredSubscriptionId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionId"].(string); ok {
        data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionSeats"].(float64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["paymentProviderSubscriptionSeats"].(int); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["paymentProviderSubscriptionSeats"].(int64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["paymentProviderSubscriptionSeats"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
        } else {
            data.PaymentProviderSubscriptionSeats = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if obj, ok := dataMap["trialEndsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.TrialEndsAt = NewRFC3339Value(val)
        } else {
            data.TrialEndsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["trialEndsAt"].(string); ok && val != "" {
        data.TrialEndsAt = NewRFC3339Value(val)
    } else {
        data.TrialEndsAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["paymentProviderCustomerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderCustomerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderCustomerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderCustomerId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderCustomerId"].(string); ok {
        data.PaymentProviderCustomerId = types.StringValue(val)
    } else {
        data.PaymentProviderCustomerId = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderSubscriptionStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderSubscriptionStatus = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderSubscriptionStatus"].(string); ok {
        data.PaymentProviderSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionStatus = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(string); ok {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["workflowRunsInLast30Days"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
        } else {
            data.WorkflowRunsInLast30Days = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(float64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiCurrentBalanceInUSDCents"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiCurrentBalanceInUsdCents = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiCurrentBalanceInUsdCents = types.NumberNull()
    }
    if obj, ok := dataMap["planName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PlanName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PlanName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PlanName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PlanName = types.StringValue(string(jsonBytes))
            } else {
                data.PlanName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PlanName = types.StringValue(string(jsonBytes))
            } else {
                data.PlanName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PlanName = types.StringValue(string(jsonBytes))
        } else {
            data.PlanName = types.StringNull()
        }
    } else if val, ok := dataMap["planName"].(string); ok {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
    }
    if obj, ok := dataMap["dataResidency"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataResidency = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataResidency = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataResidency = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataResidency = types.StringValue(string(jsonBytes))
            } else {
                data.DataResidency = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataResidency = types.StringValue(string(jsonBytes))
            } else {
                data.DataResidency = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataResidency = types.StringValue(string(jsonBytes))
        } else {
            data.DataResidency = types.StringNull()
        }
    } else if val, ok := dataMap["dataResidency"].(string); ok {
        data.DataResidency = types.StringValue(val)
    } else {
        data.DataResidency = types.StringNull()
    }
    if obj, ok := dataMap["resellerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResellerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ResellerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ResellerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ResellerId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ResellerId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ResellerId = types.StringValue(string(jsonBytes))
        } else {
            data.ResellerId = types.StringNull()
        }
    } else if val, ok := dataMap["resellerId"].(string); ok {
        data.ResellerId = types.StringValue(val)
    } else {
        data.ResellerId = types.StringNull()
    }
    if obj, ok := dataMap["resellerPlanId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResellerPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ResellerPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ResellerPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ResellerPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ResellerPlanId = types.StringValue(string(jsonBytes))
        } else {
            data.ResellerPlanId = types.StringNull()
        }
    } else if val, ok := dataMap["resellerPlanId"].(string); ok {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := dataMap["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    }
    if obj, ok := dataMap["gitHubAppInstallationId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
            } else {
                data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
            } else {
                data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
        } else {
            data.GitHubAppInstallationId = types.StringNull()
        }
    } else if val, ok := dataMap["gitHubAppInstallationId"].(string); ok {
        data.GitHubAppInstallationId = types.StringValue(val)
    } else {
        data.GitHubAppInstallationId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    // The read response is authoritative, but never let it clobber the id we just received.
    data.Id = types.StringValue(createdId)

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Write logs using the tflog package
    tflog.Trace(ctx, "created a resource")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data ProjectResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "name": true,
        "paymentProviderPlanId": true,
        "businessDetails": true,
        "businessDetailsCountry": true,
        "financeAccountingEmail": true,
        "paymentProviderPromoCode": true,
        "isFeatureFlagMonitorGroupsEnabled": true,
        "incidentNumberPrefix": true,
        "alertNumberPrefix": true,
        "scheduledMaintenanceNumberPrefix": true,
        "incidentEpisodeNumberPrefix": true,
        "alertEpisodeNumberPrefix": true,
        "enableAutomaticIncidentInvestigation": true,
        "enableAutomaticAlertInvestigation": true,
        "enableAutomaticIncidentRemediation": true,
        "enableAutomaticAlertRemediation": true,
        "enableAutomaticPostmortemDraft": true,
        "enableIncidentInstrumentationFixTasks": true,
        "enableAlertInstrumentationFixTasks": true,
        "enableAutomaticIncidentCodeFixes": true,
        "enableAutomaticAlertCodeFixes": true,
        "enableAiInsights": true,
        "enableInsightFixTasks": true,
        "autoArchiveNonActionableExceptions": true,
        "sendInvoicesByEmail": true,
        "enableAuditLogs": true,
        "isSessionReplayAllowed": true,
        "auditLogsRetentionInDays": true,
        "storeSystemEventsInAuditLogs": true,
        "requireSsoForLogin": true,
        "requireSsoWithSsoProviderId": true,
        "autoRechargeSmsOrCallByBalanceInUSD": true,
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": true,
        "enableSmsNotifications": true,
        "enableWhatsAppNotifications": true,
        "enableTelegramNotifications": true,
        "enableCallNotifications": true,
        "disableOnCallNotificationFallback": true,
        "enableAutoRechargeSmsOrCallBalance": true,
        "autoAiRechargeByBalanceInUSD": true,
        "autoRechargeAiWhenCurrentBalanceFallsInUSD": true,
        "enableAi": true,
        "aiDailyTokenLimit": true,
        "aiDailySpendLimitInUSD": true,
        "acknowledgeLinkedAlertsWhenIncidentAcknowledged": true,
        "resolveLinkedAlertsWhenIncidentResolved": true,
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
        "doNotAddGlobalProbesByDefaultOnNewMonitors": true,
        "defaultMetricCardinalityBudget": true,
        "defaultTelemetryRetentionInDays": true,
        "telemetryRetentionConfig": true,
        "defaultMetricDownsamplingRetentionDays": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "paymentProviderSubscriptionId": true,
        "paymentProviderMeteredSubscriptionId": true,
        "paymentProviderSubscriptionSeats": true,
        "trialEndsAt": true,
        "paymentProviderCustomerId": true,
        "paymentProviderSubscriptionStatus": true,
        "paymentProviderMeteredSubscriptionStatus": true,
        "createdByUserId": true,
        "workflowRunsInLast30Days": true,
        "smsOrCallCurrentBalanceInUSDCents": true,
        "aiCurrentBalanceInUSDCents": true,
        "planName": true,
        "dataResidency": true,
        "resellerId": true,
        "resellerPlanId": true,
        "letCustomerSupportAccessProject": true,
        "gitHubAppInstallationId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect(ctx, "/project/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var projectResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &projectResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse project response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := projectResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = projectResponse
    }

    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderPlanId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderPlanId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderPlanId"].(string); ok {
        data.PaymentProviderPlanId = types.StringValue(val)
    } else {
        data.PaymentProviderPlanId = types.StringNull()
    }
    if obj, ok := dataMap["businessDetails"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BusinessDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.BusinessDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.BusinessDetails = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.BusinessDetails = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.BusinessDetails = types.StringValue(string(jsonBytes))
        } else {
            data.BusinessDetails = types.StringNull()
        }
    } else if val, ok := dataMap["businessDetails"].(string); ok {
        data.BusinessDetails = types.StringValue(val)
    } else {
        data.BusinessDetails = types.StringNull()
    }
    if obj, ok := dataMap["businessDetailsCountry"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BusinessDetailsCountry = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.BusinessDetailsCountry = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
        } else {
            data.BusinessDetailsCountry = types.StringNull()
        }
    } else if val, ok := dataMap["businessDetailsCountry"].(string); ok {
        data.BusinessDetailsCountry = types.StringValue(val)
    } else {
        data.BusinessDetailsCountry = types.StringNull()
    }
    if obj, ok := dataMap["financeAccountingEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FinanceAccountingEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FinanceAccountingEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
            } else {
                data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
            } else {
                data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
        } else {
            data.FinanceAccountingEmail = types.StringNull()
        }
    } else if val, ok := dataMap["financeAccountingEmail"].(string); ok {
        data.FinanceAccountingEmail = types.StringValue(val)
    } else {
        data.FinanceAccountingEmail = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderPromoCode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderPromoCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderPromoCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderPromoCode = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderPromoCode"].(string); ok {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if val, ok := dataMap["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["incidentNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["incidentNumberPrefix"].(string); ok {
        data.IncidentNumberPrefix = types.StringValue(val)
    } else {
        data.IncidentNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["alertNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.AlertNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["alertNumberPrefix"].(string); ok {
        data.AlertNumberPrefix = types.StringValue(val)
    } else {
        data.AlertNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["scheduledMaintenanceNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceNumberPrefix"].(string); ok {
        data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["incidentEpisodeNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentEpisodeNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["incidentEpisodeNumberPrefix"].(string); ok {
        data.IncidentEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.IncidentEpisodeNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["alertEpisodeNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.AlertEpisodeNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["alertEpisodeNumberPrefix"].(string); ok {
        data.AlertEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.AlertEpisodeNumberPrefix = types.StringNull()
    }
    if val, ok := dataMap["enableAutomaticIncidentInvestigation"].(bool); ok {
        data.EnableAutomaticIncidentInvestigation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticAlertInvestigation"].(bool); ok {
        data.EnableAutomaticAlertInvestigation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticIncidentRemediation"].(bool); ok {
        data.EnableAutomaticIncidentRemediation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticAlertRemediation"].(bool); ok {
        data.EnableAutomaticAlertRemediation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticPostmortemDraft"].(bool); ok {
        data.EnableAutomaticPostmortemDraft = types.BoolValue(val)
    }
    if val, ok := dataMap["enableIncidentInstrumentationFixTasks"].(bool); ok {
        data.EnableIncidentInstrumentationFixTasks = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAlertInstrumentationFixTasks"].(bool); ok {
        data.EnableAlertInstrumentationFixTasks = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticIncidentCodeFixes"].(bool); ok {
        data.EnableAutomaticIncidentCodeFixes = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticAlertCodeFixes"].(bool); ok {
        data.EnableAutomaticAlertCodeFixes = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAiInsights"].(bool); ok {
        data.EnableAiInsights = types.BoolValue(val)
    }
    if val, ok := dataMap["enableInsightFixTasks"].(bool); ok {
        data.EnableInsightFixTasks = types.BoolValue(val)
    }
    if val, ok := dataMap["autoArchiveNonActionableExceptions"].(bool); ok {
        data.AutoArchiveNonActionableExceptions = types.BoolValue(val)
    }
    if val, ok := dataMap["sendInvoicesByEmail"].(bool); ok {
        data.SendInvoicesByEmail = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAuditLogs"].(bool); ok {
        data.EnableAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["isSessionReplayAllowed"].(bool); ok {
        data.IsSessionReplayAllowed = types.BoolValue(val)
    }
    if val, ok := dataMap["auditLogsRetentionInDays"].(float64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["auditLogsRetentionInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.AuditLogsRetentionInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AuditLogsRetentionInDays = types.NumberNull()
    }
    if val, ok := dataMap["storeSystemEventsInAuditLogs"].(bool); ok {
        data.StoreSystemEventsInAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if obj, ok := dataMap["requireSsoWithSsoProviderId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequireSsoWithSsoProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RequireSsoWithSsoProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.RequireSsoWithSsoProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["requireSsoWithSsoProviderId"].(string); ok {
        data.RequireSsoWithSsoProviderId = types.StringValue(val)
    } else {
        data.RequireSsoWithSsoProviderId = types.StringNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := dataMap["enableSmsNotifications"].(bool); ok {
        data.EnableSmsNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["enableWhatsAppNotifications"].(bool); ok {
        data.EnableWhatsAppNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["enableTelegramNotifications"].(bool); ok {
        data.EnableTelegramNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["enableCallNotifications"].(bool); ok {
        data.EnableCallNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["disableOnCallNotificationFallback"].(bool); ok {
        data.DisableOnCallNotificationFallback = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(float64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoAiRechargeByBalanceInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := dataMap["enableAi"].(bool); ok {
        data.EnableAi = types.BoolValue(val)
    }
    if val, ok := dataMap["aiDailyTokenLimit"].(float64); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailyTokenLimit"].(int); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailyTokenLimit"].(int64); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailyTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailyTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["aiDailySpendLimitInUSD"].(float64); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailySpendLimitInUSD"].(int); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailySpendLimitInUSD"].(int64); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailySpendLimitInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailySpendLimitInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailySpendLimitInUsd = types.NumberNull()
    }
    if val, ok := dataMap["acknowledgeLinkedAlertsWhenIncidentAcknowledged"].(bool); ok {
        data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged = types.BoolValue(val)
    }
    if val, ok := dataMap["resolveLinkedAlertsWhenIncidentResolved"].(bool); ok {
        data.ResolveLinkedAlertsWhenIncidentResolved = types.BoolValue(val)
    }
    if obj, ok := dataMap["alertInvestigationMinimumSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertInvestigationMinimumSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["alertInvestigationMinimumSeverityId"].(string); ok {
        data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
    } else {
        data.AlertInvestigationMinimumSeverityId = types.StringNull()
    }
    if val, ok := dataMap["aiDailyAutonomousTokenLimit"].(float64); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailyAutonomousTokenLimit"].(int); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailyAutonomousTokenLimit"].(int64); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(float64); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(int); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(int64); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(float64); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(int); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(int64); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["aiDailyFixTaskLimit"].(float64); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailyFixTaskLimit"].(int); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailyFixTaskLimit"].(int64); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailyFixTaskLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiDailyFixTaskLimit"].(float64); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiDailyFixTaskLimit"].(int); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiDailyFixTaskLimit"].(int64); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiDailyFixTaskLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := dataMap["alertAiDailyFixTaskLimit"].(float64); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiDailyFixTaskLimit"].(int); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiDailyFixTaskLimit"].(int64); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiDailyFixTaskLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(float64); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(int); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(int64); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertInvestigationDedupeWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertInvestigationDedupeWindowMinutes = types.NumberNull()
    }
    if obj, ok := dataMap["incidentInvestigationMinimumSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentInvestigationMinimumSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentInvestigationMinimumSeverityId"].(string); ok {
        data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
    } else {
        data.IncidentInvestigationMinimumSeverityId = types.StringNull()
    }
    if val, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(float64); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(int); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(int64); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentInvestigationDedupeWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberNull()
    }
    if val, ok := dataMap["aiMaxConcurrentInvestigations"].(float64); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiMaxConcurrentInvestigations"].(int); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiMaxConcurrentInvestigations"].(int64); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(float64); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(int); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(int64); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := dataMap["alertAiMaxConcurrentInvestigations"].(float64); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiMaxConcurrentInvestigations"].(int); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiMaxConcurrentInvestigations"].(int64); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(float64); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(int); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(int64); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(float64); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(int); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(int64); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiInvestigationTimeLimitInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["enableAutoRechargeAiBalance"].(bool); ok {
        data.EnableAutoRechargeAiBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["doNotAddGlobalProbesByDefaultOnNewMonitors"].(bool); ok {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = types.BoolValue(val)
    }
    if val, ok := dataMap["defaultMetricCardinalityBudget"].(float64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["defaultMetricCardinalityBudget"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
        } else {
            data.DefaultMetricCardinalityBudget = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.DefaultMetricCardinalityBudget = types.NumberNull()
    }
    if val, ok := dataMap["defaultTelemetryRetentionInDays"].(float64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["defaultTelemetryRetentionInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.DefaultTelemetryRetentionInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.DefaultTelemetryRetentionInDays = types.NumberNull()
    }
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(string); ok {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
    } else {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderSubscriptionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderSubscriptionId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderSubscriptionId"].(string); ok {
        data.PaymentProviderSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionId = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderMeteredSubscriptionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderMeteredSubscriptionId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionId"].(string); ok {
        data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionSeats"].(float64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["paymentProviderSubscriptionSeats"].(int); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["paymentProviderSubscriptionSeats"].(int64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["paymentProviderSubscriptionSeats"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
        } else {
            data.PaymentProviderSubscriptionSeats = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if obj, ok := dataMap["trialEndsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.TrialEndsAt = NewRFC3339Value(val)
        } else {
            data.TrialEndsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["trialEndsAt"].(string); ok && val != "" {
        data.TrialEndsAt = NewRFC3339Value(val)
    } else {
        data.TrialEndsAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["paymentProviderCustomerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderCustomerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderCustomerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderCustomerId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderCustomerId"].(string); ok {
        data.PaymentProviderCustomerId = types.StringValue(val)
    } else {
        data.PaymentProviderCustomerId = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderSubscriptionStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderSubscriptionStatus = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderSubscriptionStatus"].(string); ok {
        data.PaymentProviderSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionStatus = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(string); ok {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["workflowRunsInLast30Days"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
        } else {
            data.WorkflowRunsInLast30Days = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(float64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiCurrentBalanceInUSDCents"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiCurrentBalanceInUsdCents = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiCurrentBalanceInUsdCents = types.NumberNull()
    }
    if obj, ok := dataMap["planName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PlanName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PlanName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PlanName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PlanName = types.StringValue(string(jsonBytes))
            } else {
                data.PlanName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PlanName = types.StringValue(string(jsonBytes))
            } else {
                data.PlanName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PlanName = types.StringValue(string(jsonBytes))
        } else {
            data.PlanName = types.StringNull()
        }
    } else if val, ok := dataMap["planName"].(string); ok {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
    }
    if obj, ok := dataMap["dataResidency"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataResidency = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataResidency = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataResidency = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataResidency = types.StringValue(string(jsonBytes))
            } else {
                data.DataResidency = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataResidency = types.StringValue(string(jsonBytes))
            } else {
                data.DataResidency = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataResidency = types.StringValue(string(jsonBytes))
        } else {
            data.DataResidency = types.StringNull()
        }
    } else if val, ok := dataMap["dataResidency"].(string); ok {
        data.DataResidency = types.StringValue(val)
    } else {
        data.DataResidency = types.StringNull()
    }
    if obj, ok := dataMap["resellerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResellerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ResellerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ResellerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ResellerId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ResellerId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ResellerId = types.StringValue(string(jsonBytes))
        } else {
            data.ResellerId = types.StringNull()
        }
    } else if val, ok := dataMap["resellerId"].(string); ok {
        data.ResellerId = types.StringValue(val)
    } else {
        data.ResellerId = types.StringNull()
    }
    if obj, ok := dataMap["resellerPlanId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResellerPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ResellerPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ResellerPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ResellerPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ResellerPlanId = types.StringValue(string(jsonBytes))
        } else {
            data.ResellerPlanId = types.StringNull()
        }
    } else if val, ok := dataMap["resellerPlanId"].(string); ok {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := dataMap["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    }
    if obj, ok := dataMap["gitHubAppInstallationId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
            } else {
                data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
            } else {
                data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
        } else {
            data.GitHubAppInstallationId = types.StringNull()
        }
    } else if val, ok := dataMap["gitHubAppInstallationId"].(string); ok {
        data.GitHubAppInstallationId = types.StringValue(val)
    } else {
        data.GitHubAppInstallationId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data ProjectResourceModel
    var state ProjectResourceModel

    // Read Terraform current state data to get the ID
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Read Terraform plan data to get the new values
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use the ID from the current state
    data.Id = state.Id

    // What the configuration sets, and what Terraform planned (see keepPlannedValues).
    var config ProjectResourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan := data

    // Create API request body
    projectRequest := map[string]interface{}{
        "data": map[string]interface{}{},
    }
    requestDataMap := projectRequest["data"].(map[string]interface{})

    if !data.Name.IsUnknown() && !state.Name.IsUnknown() && !data.Name.Equal(state.Name) {
        requestDataMap["name"] = data.Name.ValueString()
    }
    if !data.PaymentProviderPlanId.IsUnknown() && !state.PaymentProviderPlanId.IsUnknown() && !data.PaymentProviderPlanId.Equal(state.PaymentProviderPlanId) {
        requestDataMap["paymentProviderPlanId"] = data.PaymentProviderPlanId.ValueString()
    }
    if !data.BusinessDetails.IsUnknown() && !state.BusinessDetails.IsUnknown() && !data.BusinessDetails.Equal(state.BusinessDetails) {
        requestDataMap["businessDetails"] = data.BusinessDetails.ValueString()
    }
    if !data.BusinessDetailsCountry.IsUnknown() && !state.BusinessDetailsCountry.IsUnknown() && !data.BusinessDetailsCountry.Equal(state.BusinessDetailsCountry) {
        requestDataMap["businessDetailsCountry"] = data.BusinessDetailsCountry.ValueString()
    }
    if !data.FinanceAccountingEmail.IsUnknown() && !state.FinanceAccountingEmail.IsUnknown() && !data.FinanceAccountingEmail.Equal(state.FinanceAccountingEmail) {
        requestDataMap["financeAccountingEmail"] = data.FinanceAccountingEmail.ValueString()
    }
    if !data.IsFeatureFlagMonitorGroupsEnabled.IsUnknown() && !state.IsFeatureFlagMonitorGroupsEnabled.IsUnknown() && !data.IsFeatureFlagMonitorGroupsEnabled.Equal(state.IsFeatureFlagMonitorGroupsEnabled) {
        requestDataMap["isFeatureFlagMonitorGroupsEnabled"] = data.IsFeatureFlagMonitorGroupsEnabled.ValueBool()
    }
    if !data.RequireSsoForLogin.IsUnknown() && !state.RequireSsoForLogin.IsUnknown() && !data.RequireSsoForLogin.Equal(state.RequireSsoForLogin) {
        requestDataMap["requireSsoForLogin"] = data.RequireSsoForLogin.ValueBool()
    }
    if !data.RequireSsoWithSsoProviderId.IsUnknown() && !state.RequireSsoWithSsoProviderId.IsUnknown() && !data.RequireSsoWithSsoProviderId.Equal(state.RequireSsoWithSsoProviderId) {
        requestDataMap["requireSsoWithSsoProviderId"] = data.RequireSsoWithSsoProviderId.ValueString()
    }
    if !data.IncidentNumberPrefix.IsUnknown() && !state.IncidentNumberPrefix.IsUnknown() && !data.IncidentNumberPrefix.Equal(state.IncidentNumberPrefix) {
        requestDataMap["incidentNumberPrefix"] = data.IncidentNumberPrefix.ValueString()
    }
    if !data.AlertNumberPrefix.IsUnknown() && !state.AlertNumberPrefix.IsUnknown() && !data.AlertNumberPrefix.Equal(state.AlertNumberPrefix) {
        requestDataMap["alertNumberPrefix"] = data.AlertNumberPrefix.ValueString()
    }
    if !data.ScheduledMaintenanceNumberPrefix.IsUnknown() && !state.ScheduledMaintenanceNumberPrefix.IsUnknown() && !data.ScheduledMaintenanceNumberPrefix.Equal(state.ScheduledMaintenanceNumberPrefix) {
        requestDataMap["scheduledMaintenanceNumberPrefix"] = data.ScheduledMaintenanceNumberPrefix.ValueString()
    }
    if !data.IncidentEpisodeNumberPrefix.IsUnknown() && !state.IncidentEpisodeNumberPrefix.IsUnknown() && !data.IncidentEpisodeNumberPrefix.Equal(state.IncidentEpisodeNumberPrefix) {
        requestDataMap["incidentEpisodeNumberPrefix"] = data.IncidentEpisodeNumberPrefix.ValueString()
    }
    if !data.AlertEpisodeNumberPrefix.IsUnknown() && !state.AlertEpisodeNumberPrefix.IsUnknown() && !data.AlertEpisodeNumberPrefix.Equal(state.AlertEpisodeNumberPrefix) {
        requestDataMap["alertEpisodeNumberPrefix"] = data.AlertEpisodeNumberPrefix.ValueString()
    }
    if !data.AutoRechargeSmsOrCallByBalanceInUsd.IsUnknown() && !state.AutoRechargeSmsOrCallByBalanceInUsd.IsUnknown() && !data.AutoRechargeSmsOrCallByBalanceInUsd.Equal(state.AutoRechargeSmsOrCallByBalanceInUsd) {
        requestDataMap["autoRechargeSmsOrCallByBalanceInUSD"] = r.bigFloatToFloat64(data.AutoRechargeSmsOrCallByBalanceInUsd.ValueBigFloat())
    }
    if !data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.IsUnknown() && !state.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.IsUnknown() && !data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.Equal(state.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd) {
        requestDataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"] = r.bigFloatToFloat64(data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.ValueBigFloat())
    }
    if !data.EnableSmsNotifications.IsUnknown() && !state.EnableSmsNotifications.IsUnknown() && !data.EnableSmsNotifications.Equal(state.EnableSmsNotifications) {
        requestDataMap["enableSmsNotifications"] = data.EnableSmsNotifications.ValueBool()
    }
    if !data.EnableWhatsAppNotifications.IsUnknown() && !state.EnableWhatsAppNotifications.IsUnknown() && !data.EnableWhatsAppNotifications.Equal(state.EnableWhatsAppNotifications) {
        requestDataMap["enableWhatsAppNotifications"] = data.EnableWhatsAppNotifications.ValueBool()
    }
    if !data.EnableTelegramNotifications.IsUnknown() && !state.EnableTelegramNotifications.IsUnknown() && !data.EnableTelegramNotifications.Equal(state.EnableTelegramNotifications) {
        requestDataMap["enableTelegramNotifications"] = data.EnableTelegramNotifications.ValueBool()
    }
    if !data.EnableCallNotifications.IsUnknown() && !state.EnableCallNotifications.IsUnknown() && !data.EnableCallNotifications.Equal(state.EnableCallNotifications) {
        requestDataMap["enableCallNotifications"] = data.EnableCallNotifications.ValueBool()
    }
    if !data.DisableOnCallNotificationFallback.IsUnknown() && !state.DisableOnCallNotificationFallback.IsUnknown() && !data.DisableOnCallNotificationFallback.Equal(state.DisableOnCallNotificationFallback) {
        requestDataMap["disableOnCallNotificationFallback"] = data.DisableOnCallNotificationFallback.ValueBool()
    }
    if !data.EnableAutoRechargeSmsOrCallBalance.IsUnknown() && !state.EnableAutoRechargeSmsOrCallBalance.IsUnknown() && !data.EnableAutoRechargeSmsOrCallBalance.Equal(state.EnableAutoRechargeSmsOrCallBalance) {
        requestDataMap["enableAutoRechargeSmsOrCallBalance"] = data.EnableAutoRechargeSmsOrCallBalance.ValueBool()
    }
    if !data.AutoAiRechargeByBalanceInUsd.IsUnknown() && !state.AutoAiRechargeByBalanceInUsd.IsUnknown() && !data.AutoAiRechargeByBalanceInUsd.Equal(state.AutoAiRechargeByBalanceInUsd) {
        requestDataMap["autoAiRechargeByBalanceInUSD"] = r.bigFloatToFloat64(data.AutoAiRechargeByBalanceInUsd.ValueBigFloat())
    }
    if !data.AutoRechargeAiWhenCurrentBalanceFallsInUsd.IsUnknown() && !state.AutoRechargeAiWhenCurrentBalanceFallsInUsd.IsUnknown() && !data.AutoRechargeAiWhenCurrentBalanceFallsInUsd.Equal(state.AutoRechargeAiWhenCurrentBalanceFallsInUsd) {
        requestDataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"] = r.bigFloatToFloat64(data.AutoRechargeAiWhenCurrentBalanceFallsInUsd.ValueBigFloat())
    }
    if !data.EnableAi.IsUnknown() && !state.EnableAi.IsUnknown() && !data.EnableAi.Equal(state.EnableAi) {
        requestDataMap["enableAi"] = data.EnableAi.ValueBool()
    }
    if !data.AiDailyTokenLimit.IsUnknown() && !state.AiDailyTokenLimit.IsUnknown() && !data.AiDailyTokenLimit.Equal(state.AiDailyTokenLimit) {
        requestDataMap["aiDailyTokenLimit"] = r.bigFloatToFloat64(data.AiDailyTokenLimit.ValueBigFloat())
    }
    if !data.AiDailySpendLimitInUsd.IsUnknown() && !state.AiDailySpendLimitInUsd.IsUnknown() && !data.AiDailySpendLimitInUsd.Equal(state.AiDailySpendLimitInUsd) {
        requestDataMap["aiDailySpendLimitInUSD"] = r.bigFloatToFloat64(data.AiDailySpendLimitInUsd.ValueBigFloat())
    }
    if !data.EnableAutomaticIncidentInvestigation.IsUnknown() && !state.EnableAutomaticIncidentInvestigation.IsUnknown() && !data.EnableAutomaticIncidentInvestigation.Equal(state.EnableAutomaticIncidentInvestigation) {
        requestDataMap["enableAutomaticIncidentInvestigation"] = data.EnableAutomaticIncidentInvestigation.ValueBool()
    }
    if !data.EnableAutomaticAlertInvestigation.IsUnknown() && !state.EnableAutomaticAlertInvestigation.IsUnknown() && !data.EnableAutomaticAlertInvestigation.Equal(state.EnableAutomaticAlertInvestigation) {
        requestDataMap["enableAutomaticAlertInvestigation"] = data.EnableAutomaticAlertInvestigation.ValueBool()
    }
    if !data.EnableAutomaticIncidentRemediation.IsUnknown() && !state.EnableAutomaticIncidentRemediation.IsUnknown() && !data.EnableAutomaticIncidentRemediation.Equal(state.EnableAutomaticIncidentRemediation) {
        requestDataMap["enableAutomaticIncidentRemediation"] = data.EnableAutomaticIncidentRemediation.ValueBool()
    }
    if !data.EnableAutomaticAlertRemediation.IsUnknown() && !state.EnableAutomaticAlertRemediation.IsUnknown() && !data.EnableAutomaticAlertRemediation.Equal(state.EnableAutomaticAlertRemediation) {
        requestDataMap["enableAutomaticAlertRemediation"] = data.EnableAutomaticAlertRemediation.ValueBool()
    }
    if !data.EnableAutomaticPostmortemDraft.IsUnknown() && !state.EnableAutomaticPostmortemDraft.IsUnknown() && !data.EnableAutomaticPostmortemDraft.Equal(state.EnableAutomaticPostmortemDraft) {
        requestDataMap["enableAutomaticPostmortemDraft"] = data.EnableAutomaticPostmortemDraft.ValueBool()
    }
    if !data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.IsUnknown() && !state.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.IsUnknown() && !data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.Equal(state.AcknowledgeLinkedAlertsWhenIncidentAcknowledged) {
        requestDataMap["acknowledgeLinkedAlertsWhenIncidentAcknowledged"] = data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.ValueBool()
    }
    if !data.ResolveLinkedAlertsWhenIncidentResolved.IsUnknown() && !state.ResolveLinkedAlertsWhenIncidentResolved.IsUnknown() && !data.ResolveLinkedAlertsWhenIncidentResolved.Equal(state.ResolveLinkedAlertsWhenIncidentResolved) {
        requestDataMap["resolveLinkedAlertsWhenIncidentResolved"] = data.ResolveLinkedAlertsWhenIncidentResolved.ValueBool()
    }
    if !data.EnableIncidentInstrumentationFixTasks.IsUnknown() && !state.EnableIncidentInstrumentationFixTasks.IsUnknown() && !data.EnableIncidentInstrumentationFixTasks.Equal(state.EnableIncidentInstrumentationFixTasks) {
        requestDataMap["enableIncidentInstrumentationFixTasks"] = data.EnableIncidentInstrumentationFixTasks.ValueBool()
    }
    if !data.EnableAlertInstrumentationFixTasks.IsUnknown() && !state.EnableAlertInstrumentationFixTasks.IsUnknown() && !data.EnableAlertInstrumentationFixTasks.Equal(state.EnableAlertInstrumentationFixTasks) {
        requestDataMap["enableAlertInstrumentationFixTasks"] = data.EnableAlertInstrumentationFixTasks.ValueBool()
    }
    if !data.EnableAutomaticIncidentCodeFixes.IsUnknown() && !state.EnableAutomaticIncidentCodeFixes.IsUnknown() && !data.EnableAutomaticIncidentCodeFixes.Equal(state.EnableAutomaticIncidentCodeFixes) {
        requestDataMap["enableAutomaticIncidentCodeFixes"] = data.EnableAutomaticIncidentCodeFixes.ValueBool()
    }
    if !data.EnableAutomaticAlertCodeFixes.IsUnknown() && !state.EnableAutomaticAlertCodeFixes.IsUnknown() && !data.EnableAutomaticAlertCodeFixes.Equal(state.EnableAutomaticAlertCodeFixes) {
        requestDataMap["enableAutomaticAlertCodeFixes"] = data.EnableAutomaticAlertCodeFixes.ValueBool()
    }
    if !data.EnableAiInsights.IsUnknown() && !state.EnableAiInsights.IsUnknown() && !data.EnableAiInsights.Equal(state.EnableAiInsights) {
        requestDataMap["enableAiInsights"] = data.EnableAiInsights.ValueBool()
    }
    if !data.EnableInsightFixTasks.IsUnknown() && !state.EnableInsightFixTasks.IsUnknown() && !data.EnableInsightFixTasks.Equal(state.EnableInsightFixTasks) {
        requestDataMap["enableInsightFixTasks"] = data.EnableInsightFixTasks.ValueBool()
    }
    if !data.AutoArchiveNonActionableExceptions.IsUnknown() && !state.AutoArchiveNonActionableExceptions.IsUnknown() && !data.AutoArchiveNonActionableExceptions.Equal(state.AutoArchiveNonActionableExceptions) {
        requestDataMap["autoArchiveNonActionableExceptions"] = data.AutoArchiveNonActionableExceptions.ValueBool()
    }
    if !data.AlertInvestigationMinimumSeverityId.IsUnknown() && !state.AlertInvestigationMinimumSeverityId.IsUnknown() && !data.AlertInvestigationMinimumSeverityId.Equal(state.AlertInvestigationMinimumSeverityId) {
        requestDataMap["alertInvestigationMinimumSeverityId"] = data.AlertInvestigationMinimumSeverityId.ValueString()
    }
    if !data.AiDailyAutonomousTokenLimit.IsUnknown() && !state.AiDailyAutonomousTokenLimit.IsUnknown() && !data.AiDailyAutonomousTokenLimit.Equal(state.AiDailyAutonomousTokenLimit) {
        requestDataMap["aiDailyAutonomousTokenLimit"] = r.bigFloatToFloat64(data.AiDailyAutonomousTokenLimit.ValueBigFloat())
    }
    if !data.IncidentAiDailyAutonomousTokenLimit.IsUnknown() && !state.IncidentAiDailyAutonomousTokenLimit.IsUnknown() && !data.IncidentAiDailyAutonomousTokenLimit.Equal(state.IncidentAiDailyAutonomousTokenLimit) {
        requestDataMap["incidentAiDailyAutonomousTokenLimit"] = r.bigFloatToFloat64(data.IncidentAiDailyAutonomousTokenLimit.ValueBigFloat())
    }
    if !data.AlertAiDailyAutonomousTokenLimit.IsUnknown() && !state.AlertAiDailyAutonomousTokenLimit.IsUnknown() && !data.AlertAiDailyAutonomousTokenLimit.Equal(state.AlertAiDailyAutonomousTokenLimit) {
        requestDataMap["alertAiDailyAutonomousTokenLimit"] = r.bigFloatToFloat64(data.AlertAiDailyAutonomousTokenLimit.ValueBigFloat())
    }
    if !data.AiDailyFixTaskLimit.IsUnknown() && !state.AiDailyFixTaskLimit.IsUnknown() && !data.AiDailyFixTaskLimit.Equal(state.AiDailyFixTaskLimit) {
        requestDataMap["aiDailyFixTaskLimit"] = r.bigFloatToFloat64(data.AiDailyFixTaskLimit.ValueBigFloat())
    }
    if !data.IncidentAiDailyFixTaskLimit.IsUnknown() && !state.IncidentAiDailyFixTaskLimit.IsUnknown() && !data.IncidentAiDailyFixTaskLimit.Equal(state.IncidentAiDailyFixTaskLimit) {
        requestDataMap["incidentAiDailyFixTaskLimit"] = r.bigFloatToFloat64(data.IncidentAiDailyFixTaskLimit.ValueBigFloat())
    }
    if !data.AlertAiDailyFixTaskLimit.IsUnknown() && !state.AlertAiDailyFixTaskLimit.IsUnknown() && !data.AlertAiDailyFixTaskLimit.Equal(state.AlertAiDailyFixTaskLimit) {
        requestDataMap["alertAiDailyFixTaskLimit"] = r.bigFloatToFloat64(data.AlertAiDailyFixTaskLimit.ValueBigFloat())
    }
    if !data.AlertInvestigationDedupeWindowMinutes.IsUnknown() && !state.AlertInvestigationDedupeWindowMinutes.IsUnknown() && !data.AlertInvestigationDedupeWindowMinutes.Equal(state.AlertInvestigationDedupeWindowMinutes) {
        requestDataMap["alertInvestigationDedupeWindowMinutes"] = r.bigFloatToFloat64(data.AlertInvestigationDedupeWindowMinutes.ValueBigFloat())
    }
    if !data.IncidentInvestigationMinimumSeverityId.IsUnknown() && !state.IncidentInvestigationMinimumSeverityId.IsUnknown() && !data.IncidentInvestigationMinimumSeverityId.Equal(state.IncidentInvestigationMinimumSeverityId) {
        requestDataMap["incidentInvestigationMinimumSeverityId"] = data.IncidentInvestigationMinimumSeverityId.ValueString()
    }
    if !data.IncidentInvestigationDedupeWindowMinutes.IsUnknown() && !state.IncidentInvestigationDedupeWindowMinutes.IsUnknown() && !data.IncidentInvestigationDedupeWindowMinutes.Equal(state.IncidentInvestigationDedupeWindowMinutes) {
        requestDataMap["incidentInvestigationDedupeWindowMinutes"] = r.bigFloatToFloat64(data.IncidentInvestigationDedupeWindowMinutes.ValueBigFloat())
    }
    if !data.AiMaxConcurrentInvestigations.IsUnknown() && !state.AiMaxConcurrentInvestigations.IsUnknown() && !data.AiMaxConcurrentInvestigations.Equal(state.AiMaxConcurrentInvestigations) {
        requestDataMap["aiMaxConcurrentInvestigations"] = r.bigFloatToFloat64(data.AiMaxConcurrentInvestigations.ValueBigFloat())
    }
    if !data.IncidentAiMaxConcurrentInvestigations.IsUnknown() && !state.IncidentAiMaxConcurrentInvestigations.IsUnknown() && !data.IncidentAiMaxConcurrentInvestigations.Equal(state.IncidentAiMaxConcurrentInvestigations) {
        requestDataMap["incidentAiMaxConcurrentInvestigations"] = r.bigFloatToFloat64(data.IncidentAiMaxConcurrentInvestigations.ValueBigFloat())
    }
    if !data.AlertAiMaxConcurrentInvestigations.IsUnknown() && !state.AlertAiMaxConcurrentInvestigations.IsUnknown() && !data.AlertAiMaxConcurrentInvestigations.Equal(state.AlertAiMaxConcurrentInvestigations) {
        requestDataMap["alertAiMaxConcurrentInvestigations"] = r.bigFloatToFloat64(data.AlertAiMaxConcurrentInvestigations.ValueBigFloat())
    }
    if !data.IncidentAiInvestigationTimeLimitInMinutes.IsUnknown() && !state.IncidentAiInvestigationTimeLimitInMinutes.IsUnknown() && !data.IncidentAiInvestigationTimeLimitInMinutes.Equal(state.IncidentAiInvestigationTimeLimitInMinutes) {
        requestDataMap["incidentAiInvestigationTimeLimitInMinutes"] = r.bigFloatToFloat64(data.IncidentAiInvestigationTimeLimitInMinutes.ValueBigFloat())
    }
    if !data.AlertAiInvestigationTimeLimitInMinutes.IsUnknown() && !state.AlertAiInvestigationTimeLimitInMinutes.IsUnknown() && !data.AlertAiInvestigationTimeLimitInMinutes.Equal(state.AlertAiInvestigationTimeLimitInMinutes) {
        requestDataMap["alertAiInvestigationTimeLimitInMinutes"] = r.bigFloatToFloat64(data.AlertAiInvestigationTimeLimitInMinutes.ValueBigFloat())
    }
    if !data.EnableAutoRechargeAiBalance.IsUnknown() && !state.EnableAutoRechargeAiBalance.IsUnknown() && !data.EnableAutoRechargeAiBalance.Equal(state.EnableAutoRechargeAiBalance) {
        requestDataMap["enableAutoRechargeAiBalance"] = data.EnableAutoRechargeAiBalance.ValueBool()
    }
    if !data.SendInvoicesByEmail.IsUnknown() && !state.SendInvoicesByEmail.IsUnknown() && !data.SendInvoicesByEmail.Equal(state.SendInvoicesByEmail) {
        requestDataMap["sendInvoicesByEmail"] = data.SendInvoicesByEmail.ValueBool()
    }
    if !data.DoNotAddGlobalProbesByDefaultOnNewMonitors.IsUnknown() && !state.DoNotAddGlobalProbesByDefaultOnNewMonitors.IsUnknown() && !data.DoNotAddGlobalProbesByDefaultOnNewMonitors.Equal(state.DoNotAddGlobalProbesByDefaultOnNewMonitors) {
        requestDataMap["doNotAddGlobalProbesByDefaultOnNewMonitors"] = data.DoNotAddGlobalProbesByDefaultOnNewMonitors.ValueBool()
    }
    if !data.DefaultMetricCardinalityBudget.IsUnknown() && !state.DefaultMetricCardinalityBudget.IsUnknown() && !data.DefaultMetricCardinalityBudget.Equal(state.DefaultMetricCardinalityBudget) {
        requestDataMap["defaultMetricCardinalityBudget"] = r.bigFloatToFloat64(data.DefaultMetricCardinalityBudget.ValueBigFloat())
    }
    if !data.DefaultTelemetryRetentionInDays.IsUnknown() && !state.DefaultTelemetryRetentionInDays.IsUnknown() && !data.DefaultTelemetryRetentionInDays.Equal(state.DefaultTelemetryRetentionInDays) {
        requestDataMap["defaultTelemetryRetentionInDays"] = r.bigFloatToFloat64(data.DefaultTelemetryRetentionInDays.ValueBigFloat())
    }
    if !data.TelemetryRetentionConfig.IsUnknown() && !state.TelemetryRetentionConfig.IsUnknown() && !data.TelemetryRetentionConfig.Equal(state.TelemetryRetentionConfig) {
        var telemetryretentionconfigData interface{}
        if err := json.Unmarshal([]byte(data.TelemetryRetentionConfig.ValueString()), &telemetryretentionconfigData); err == nil {
            requestDataMap["telemetryRetentionConfig"] = telemetryretentionconfigData
        } else {
            requestDataMap["telemetryRetentionConfig"] = data.TelemetryRetentionConfig.ValueString()
        }
    }
    if !data.DefaultMetricDownsamplingRetentionDays.IsUnknown() && !state.DefaultMetricDownsamplingRetentionDays.IsUnknown() && !data.DefaultMetricDownsamplingRetentionDays.Equal(state.DefaultMetricDownsamplingRetentionDays) {
        var defaultmetricdownsamplingretentiondaysData interface{}
        if err := json.Unmarshal([]byte(data.DefaultMetricDownsamplingRetentionDays.ValueString()), &defaultmetricdownsamplingretentiondaysData); err == nil {
            requestDataMap["defaultMetricDownsamplingRetentionDays"] = defaultmetricdownsamplingretentiondaysData
        } else {
            requestDataMap["defaultMetricDownsamplingRetentionDays"] = data.DefaultMetricDownsamplingRetentionDays.ValueString()
        }
    }
    if !data.EnableAuditLogs.IsUnknown() && !state.EnableAuditLogs.IsUnknown() && !data.EnableAuditLogs.Equal(state.EnableAuditLogs) {
        requestDataMap["enableAuditLogs"] = data.EnableAuditLogs.ValueBool()
    }
    if !data.IsSessionReplayAllowed.IsUnknown() && !state.IsSessionReplayAllowed.IsUnknown() && !data.IsSessionReplayAllowed.Equal(state.IsSessionReplayAllowed) {
        requestDataMap["isSessionReplayAllowed"] = data.IsSessionReplayAllowed.ValueBool()
    }
    if !data.AuditLogsRetentionInDays.IsUnknown() && !state.AuditLogsRetentionInDays.IsUnknown() && !data.AuditLogsRetentionInDays.Equal(state.AuditLogsRetentionInDays) {
        requestDataMap["auditLogsRetentionInDays"] = r.bigFloatToFloat64(data.AuditLogsRetentionInDays.ValueBigFloat())
    }
    if !data.StoreSystemEventsInAuditLogs.IsUnknown() && !state.StoreSystemEventsInAuditLogs.IsUnknown() && !data.StoreSystemEventsInAuditLogs.Equal(state.StoreSystemEventsInAuditLogs) {
        requestDataMap["storeSystemEventsInAuditLogs"] = data.StoreSystemEventsInAuditLogs.ValueBool()
    }

    // Only call the API when there are changed fields to send. An empty
    // update body is rejected by the API; state is still refreshed below so
    // this method never writes unverified plan values into state.
    if len(projectRequest["data"].(map[string]interface{})) > 0 {
        httpResp, err := r.client.Put(ctx, "/project/" + data.Id.ValueString() + "", projectRequest)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update project, got error: %s", err))
            return
        }

        // Parse the update response
        var projectResponse map[string]interface{}
        err = r.client.ParseResponse(httpResp, &projectResponse)
        if err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to update project: %s", err))
            return
        }
        _ = projectResponse
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "name": true,
        "paymentProviderPlanId": true,
        "businessDetails": true,
        "businessDetailsCountry": true,
        "financeAccountingEmail": true,
        "paymentProviderPromoCode": true,
        "isFeatureFlagMonitorGroupsEnabled": true,
        "incidentNumberPrefix": true,
        "alertNumberPrefix": true,
        "scheduledMaintenanceNumberPrefix": true,
        "incidentEpisodeNumberPrefix": true,
        "alertEpisodeNumberPrefix": true,
        "enableAutomaticIncidentInvestigation": true,
        "enableAutomaticAlertInvestigation": true,
        "enableAutomaticIncidentRemediation": true,
        "enableAutomaticAlertRemediation": true,
        "enableAutomaticPostmortemDraft": true,
        "enableIncidentInstrumentationFixTasks": true,
        "enableAlertInstrumentationFixTasks": true,
        "enableAutomaticIncidentCodeFixes": true,
        "enableAutomaticAlertCodeFixes": true,
        "enableAiInsights": true,
        "enableInsightFixTasks": true,
        "autoArchiveNonActionableExceptions": true,
        "sendInvoicesByEmail": true,
        "enableAuditLogs": true,
        "isSessionReplayAllowed": true,
        "auditLogsRetentionInDays": true,
        "storeSystemEventsInAuditLogs": true,
        "requireSsoForLogin": true,
        "requireSsoWithSsoProviderId": true,
        "autoRechargeSmsOrCallByBalanceInUSD": true,
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": true,
        "enableSmsNotifications": true,
        "enableWhatsAppNotifications": true,
        "enableTelegramNotifications": true,
        "enableCallNotifications": true,
        "disableOnCallNotificationFallback": true,
        "enableAutoRechargeSmsOrCallBalance": true,
        "autoAiRechargeByBalanceInUSD": true,
        "autoRechargeAiWhenCurrentBalanceFallsInUSD": true,
        "enableAi": true,
        "aiDailyTokenLimit": true,
        "aiDailySpendLimitInUSD": true,
        "acknowledgeLinkedAlertsWhenIncidentAcknowledged": true,
        "resolveLinkedAlertsWhenIncidentResolved": true,
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
        "doNotAddGlobalProbesByDefaultOnNewMonitors": true,
        "defaultMetricCardinalityBudget": true,
        "defaultTelemetryRetentionInDays": true,
        "telemetryRetentionConfig": true,
        "defaultMetricDownsamplingRetentionDays": true,
        "createdAt": true,
        "updatedAt": true,
        "slug": true,
        "paymentProviderSubscriptionId": true,
        "paymentProviderMeteredSubscriptionId": true,
        "paymentProviderSubscriptionSeats": true,
        "trialEndsAt": true,
        "paymentProviderCustomerId": true,
        "paymentProviderSubscriptionStatus": true,
        "paymentProviderMeteredSubscriptionStatus": true,
        "createdByUserId": true,
        "workflowRunsInLast30Days": true,
        "smsOrCallCurrentBalanceInUSDCents": true,
        "aiCurrentBalanceInUSDCents": true,
        "planName": true,
        "dataResidency": true,
        "resellerId": true,
        "resellerPlanId": true,
        "letCustomerSupportAccessProject": true,
        "gitHubAppInstallationId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect(ctx, "/project/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read project after update: %s", err))
        return
    }

    // Update the model with response data from the Read operation
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := readResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = readResponse
    }

    if obj, ok := dataMap["name"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Name = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Name = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Name = types.StringValue(string(jsonBytes))
            } else {
                data.Name = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Name = types.StringValue(string(jsonBytes))
        } else {
            data.Name = types.StringNull()
        }
    } else if val, ok := dataMap["name"].(string); ok {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderPlanId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPlanId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderPlanId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderPlanId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderPlanId"].(string); ok {
        data.PaymentProviderPlanId = types.StringValue(val)
    } else {
        data.PaymentProviderPlanId = types.StringNull()
    }
    if obj, ok := dataMap["businessDetails"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BusinessDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.BusinessDetails = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.BusinessDetails = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.BusinessDetails = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetails = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.BusinessDetails = types.StringValue(string(jsonBytes))
        } else {
            data.BusinessDetails = types.StringNull()
        }
    } else if val, ok := dataMap["businessDetails"].(string); ok {
        data.BusinessDetails = types.StringValue(val)
    } else {
        data.BusinessDetails = types.StringNull()
    }
    if obj, ok := dataMap["businessDetailsCountry"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.BusinessDetailsCountry = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.BusinessDetailsCountry = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
            } else {
                data.BusinessDetailsCountry = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.BusinessDetailsCountry = types.StringValue(string(jsonBytes))
        } else {
            data.BusinessDetailsCountry = types.StringNull()
        }
    } else if val, ok := dataMap["businessDetailsCountry"].(string); ok {
        data.BusinessDetailsCountry = types.StringValue(val)
    } else {
        data.BusinessDetailsCountry = types.StringNull()
    }
    if obj, ok := dataMap["financeAccountingEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.FinanceAccountingEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.FinanceAccountingEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
            } else {
                data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
            } else {
                data.FinanceAccountingEmail = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.FinanceAccountingEmail = types.StringValue(string(jsonBytes))
        } else {
            data.FinanceAccountingEmail = types.StringNull()
        }
    } else if val, ok := dataMap["financeAccountingEmail"].(string); ok {
        data.FinanceAccountingEmail = types.StringValue(val)
    } else {
        data.FinanceAccountingEmail = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderPromoCode"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderPromoCode = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderPromoCode = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderPromoCode = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderPromoCode = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderPromoCode = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderPromoCode"].(string); ok {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if val, ok := dataMap["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    }
    if obj, ok := dataMap["incidentNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["incidentNumberPrefix"].(string); ok {
        data.IncidentNumberPrefix = types.StringValue(val)
    } else {
        data.IncidentNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["alertNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.AlertNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["alertNumberPrefix"].(string); ok {
        data.AlertNumberPrefix = types.StringValue(val)
    } else {
        data.AlertNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["scheduledMaintenanceNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.ScheduledMaintenanceNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ScheduledMaintenanceNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.ScheduledMaintenanceNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["scheduledMaintenanceNumberPrefix"].(string); ok {
        data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
    } else {
        data.ScheduledMaintenanceNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["incidentEpisodeNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentEpisodeNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["incidentEpisodeNumberPrefix"].(string); ok {
        data.IncidentEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.IncidentEpisodeNumberPrefix = types.StringNull()
    }
    if obj, ok := dataMap["alertEpisodeNumberPrefix"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertEpisodeNumberPrefix = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
            } else {
                data.AlertEpisodeNumberPrefix = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertEpisodeNumberPrefix = types.StringValue(string(jsonBytes))
        } else {
            data.AlertEpisodeNumberPrefix = types.StringNull()
        }
    } else if val, ok := dataMap["alertEpisodeNumberPrefix"].(string); ok {
        data.AlertEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.AlertEpisodeNumberPrefix = types.StringNull()
    }
    if val, ok := dataMap["enableAutomaticIncidentInvestigation"].(bool); ok {
        data.EnableAutomaticIncidentInvestigation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticAlertInvestigation"].(bool); ok {
        data.EnableAutomaticAlertInvestigation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticIncidentRemediation"].(bool); ok {
        data.EnableAutomaticIncidentRemediation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticAlertRemediation"].(bool); ok {
        data.EnableAutomaticAlertRemediation = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticPostmortemDraft"].(bool); ok {
        data.EnableAutomaticPostmortemDraft = types.BoolValue(val)
    }
    if val, ok := dataMap["enableIncidentInstrumentationFixTasks"].(bool); ok {
        data.EnableIncidentInstrumentationFixTasks = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAlertInstrumentationFixTasks"].(bool); ok {
        data.EnableAlertInstrumentationFixTasks = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticIncidentCodeFixes"].(bool); ok {
        data.EnableAutomaticIncidentCodeFixes = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutomaticAlertCodeFixes"].(bool); ok {
        data.EnableAutomaticAlertCodeFixes = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAiInsights"].(bool); ok {
        data.EnableAiInsights = types.BoolValue(val)
    }
    if val, ok := dataMap["enableInsightFixTasks"].(bool); ok {
        data.EnableInsightFixTasks = types.BoolValue(val)
    }
    if val, ok := dataMap["autoArchiveNonActionableExceptions"].(bool); ok {
        data.AutoArchiveNonActionableExceptions = types.BoolValue(val)
    }
    if val, ok := dataMap["sendInvoicesByEmail"].(bool); ok {
        data.SendInvoicesByEmail = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAuditLogs"].(bool); ok {
        data.EnableAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["isSessionReplayAllowed"].(bool); ok {
        data.IsSessionReplayAllowed = types.BoolValue(val)
    }
    if val, ok := dataMap["auditLogsRetentionInDays"].(float64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["auditLogsRetentionInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.AuditLogsRetentionInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AuditLogsRetentionInDays = types.NumberNull()
    }
    if val, ok := dataMap["storeSystemEventsInAuditLogs"].(bool); ok {
        data.StoreSystemEventsInAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if obj, ok := dataMap["requireSsoWithSsoProviderId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.RequireSsoWithSsoProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.RequireSsoWithSsoProviderId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
            } else {
                data.RequireSsoWithSsoProviderId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.RequireSsoWithSsoProviderId = types.StringValue(string(jsonBytes))
        } else {
            data.RequireSsoWithSsoProviderId = types.StringNull()
        }
    } else if val, ok := dataMap["requireSsoWithSsoProviderId"].(string); ok {
        data.RequireSsoWithSsoProviderId = types.StringValue(val)
    } else {
        data.RequireSsoWithSsoProviderId = types.StringNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := dataMap["enableSmsNotifications"].(bool); ok {
        data.EnableSmsNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["enableWhatsAppNotifications"].(bool); ok {
        data.EnableWhatsAppNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["enableTelegramNotifications"].(bool); ok {
        data.EnableTelegramNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["enableCallNotifications"].(bool); ok {
        data.EnableCallNotifications = types.BoolValue(val)
    }
    if val, ok := dataMap["disableOnCallNotificationFallback"].(bool); ok {
        data.DisableOnCallNotificationFallback = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(float64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoAiRechargeByBalanceInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := dataMap["enableAi"].(bool); ok {
        data.EnableAi = types.BoolValue(val)
    }
    if val, ok := dataMap["aiDailyTokenLimit"].(float64); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailyTokenLimit"].(int); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailyTokenLimit"].(int64); ok {
        data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailyTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailyTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["aiDailySpendLimitInUSD"].(float64); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailySpendLimitInUSD"].(int); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailySpendLimitInUSD"].(int64); ok {
        data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailySpendLimitInUSD"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailySpendLimitInUsd = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailySpendLimitInUsd = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailySpendLimitInUsd = types.NumberNull()
    }
    if val, ok := dataMap["acknowledgeLinkedAlertsWhenIncidentAcknowledged"].(bool); ok {
        data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged = types.BoolValue(val)
    }
    if val, ok := dataMap["resolveLinkedAlertsWhenIncidentResolved"].(bool); ok {
        data.ResolveLinkedAlertsWhenIncidentResolved = types.BoolValue(val)
    }
    if obj, ok := dataMap["alertInvestigationMinimumSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.AlertInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.AlertInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.AlertInvestigationMinimumSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["alertInvestigationMinimumSeverityId"].(string); ok {
        data.AlertInvestigationMinimumSeverityId = types.StringValue(val)
    } else {
        data.AlertInvestigationMinimumSeverityId = types.StringNull()
    }
    if val, ok := dataMap["aiDailyAutonomousTokenLimit"].(float64); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailyAutonomousTokenLimit"].(int); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailyAutonomousTokenLimit"].(int64); ok {
        data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(float64); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(int); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(int64); ok {
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(float64); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(int); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(int64); ok {
        data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiDailyAutonomousTokenLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiDailyAutonomousTokenLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiDailyAutonomousTokenLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiDailyAutonomousTokenLimit = types.NumberNull()
    }
    if val, ok := dataMap["aiDailyFixTaskLimit"].(float64); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiDailyFixTaskLimit"].(int); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiDailyFixTaskLimit"].(int64); ok {
        data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiDailyFixTaskLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiDailyFixTaskLimit"].(float64); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiDailyFixTaskLimit"].(int); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiDailyFixTaskLimit"].(int64); ok {
        data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiDailyFixTaskLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := dataMap["alertAiDailyFixTaskLimit"].(float64); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiDailyFixTaskLimit"].(int); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiDailyFixTaskLimit"].(int64); ok {
        data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiDailyFixTaskLimit"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiDailyFixTaskLimit = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiDailyFixTaskLimit = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiDailyFixTaskLimit = types.NumberNull()
    }
    if val, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(float64); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(int); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(int64); ok {
        data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertInvestigationDedupeWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertInvestigationDedupeWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertInvestigationDedupeWindowMinutes = types.NumberNull()
    }
    if obj, ok := dataMap["incidentInvestigationMinimumSeverityId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
            } else {
                data.IncidentInvestigationMinimumSeverityId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.IncidentInvestigationMinimumSeverityId = types.StringValue(string(jsonBytes))
        } else {
            data.IncidentInvestigationMinimumSeverityId = types.StringNull()
        }
    } else if val, ok := dataMap["incidentInvestigationMinimumSeverityId"].(string); ok {
        data.IncidentInvestigationMinimumSeverityId = types.StringValue(val)
    } else {
        data.IncidentInvestigationMinimumSeverityId = types.StringNull()
    }
    if val, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(float64); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(int); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(int64); ok {
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentInvestigationDedupeWindowMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentInvestigationDedupeWindowMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentInvestigationDedupeWindowMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentInvestigationDedupeWindowMinutes = types.NumberNull()
    }
    if val, ok := dataMap["aiMaxConcurrentInvestigations"].(float64); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiMaxConcurrentInvestigations"].(int); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiMaxConcurrentInvestigations"].(int64); ok {
        data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(float64); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(int); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(int64); ok {
        data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := dataMap["alertAiMaxConcurrentInvestigations"].(float64); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiMaxConcurrentInvestigations"].(int); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiMaxConcurrentInvestigations"].(int64); ok {
        data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiMaxConcurrentInvestigations"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiMaxConcurrentInvestigations = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiMaxConcurrentInvestigations = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiMaxConcurrentInvestigations = types.NumberNull()
    }
    if val, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(float64); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(int); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(int64); ok {
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["incidentAiInvestigationTimeLimitInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.IncidentAiInvestigationTimeLimitInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(float64); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(int); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(int64); ok {
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["alertAiInvestigationTimeLimitInMinutes"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AlertAiInvestigationTimeLimitInMinutes = types.NumberValue(big.NewFloat(val))
        } else {
            data.AlertAiInvestigationTimeLimitInMinutes = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AlertAiInvestigationTimeLimitInMinutes = types.NumberNull()
    }
    if val, ok := dataMap["enableAutoRechargeAiBalance"].(bool); ok {
        data.EnableAutoRechargeAiBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["doNotAddGlobalProbesByDefaultOnNewMonitors"].(bool); ok {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = types.BoolValue(val)
    }
    if val, ok := dataMap["defaultMetricCardinalityBudget"].(float64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["defaultMetricCardinalityBudget"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
        } else {
            data.DefaultMetricCardinalityBudget = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.DefaultMetricCardinalityBudget = types.NumberNull()
    }
    if val, ok := dataMap["defaultTelemetryRetentionInDays"].(float64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["defaultTelemetryRetentionInDays"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
        } else {
            data.DefaultTelemetryRetentionInDays = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.DefaultTelemetryRetentionInDays = types.NumberNull()
    }
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TelemetryRetentionConfig = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TelemetryRetentionConfig = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TelemetryRetentionConfig = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(string); ok {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
    } else {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.CreatedAt = NewRFC3339Value(val)
        } else {
            data.CreatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewRFC3339Value(val)
    } else {
        data.CreatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.UpdatedAt = NewRFC3339Value(val)
        } else {
            data.UpdatedAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewRFC3339Value(val)
    } else {
        data.UpdatedAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["slug"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Slug = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Slug = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Slug = types.StringValue(string(jsonBytes))
            } else {
                data.Slug = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Slug = types.StringValue(string(jsonBytes))
        } else {
            data.Slug = types.StringNull()
        }
    } else if val, ok := dataMap["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderSubscriptionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderSubscriptionId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderSubscriptionId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderSubscriptionId"].(string); ok {
        data.PaymentProviderSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionId = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderMeteredSubscriptionId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderMeteredSubscriptionId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderMeteredSubscriptionId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionId"].(string); ok {
        data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionSeats"].(float64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["paymentProviderSubscriptionSeats"].(int); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["paymentProviderSubscriptionSeats"].(int64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["paymentProviderSubscriptionSeats"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
        } else {
            data.PaymentProviderSubscriptionSeats = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if obj, ok := dataMap["trialEndsAt"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok && val != "" {
            data.TrialEndsAt = NewRFC3339Value(val)
        } else {
            data.TrialEndsAt = NewRFC3339Null()
        }
    } else if val, ok := dataMap["trialEndsAt"].(string); ok && val != "" {
        data.TrialEndsAt = NewRFC3339Value(val)
    } else {
        data.TrialEndsAt = NewRFC3339Null()
    }
    if obj, ok := dataMap["paymentProviderCustomerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderCustomerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderCustomerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderCustomerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderCustomerId = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderCustomerId = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderCustomerId"].(string); ok {
        data.PaymentProviderCustomerId = types.StringValue(val)
    } else {
        data.PaymentProviderCustomerId = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderSubscriptionStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderSubscriptionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderSubscriptionStatus = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderSubscriptionStatus"].(string); ok {
        data.PaymentProviderSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionStatus = types.StringNull()
    }
    if obj, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
            } else {
                data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(string(jsonBytes))
        } else {
            data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
        }
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(string); ok {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
    }
    if obj, ok := dataMap["createdByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["createdByUserId"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["workflowRunsInLast30Days"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
        } else {
            data.WorkflowRunsInLast30Days = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(float64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if obj, ok := dataMap["aiCurrentBalanceInUSDCents"].(map[string]interface{}); ok {
        // Unwrap numeric wrapper objects (e.g. {_type: "Port", value: 443})
        if val, ok := obj["value"].(float64); ok {
            data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
        } else {
            data.AiCurrentBalanceInUsdCents = types.NumberNull()
        }
    } else {
        // Missing or unrecognized value: null, never unknown, so apply can complete.
        data.AiCurrentBalanceInUsdCents = types.NumberNull()
    }
    if obj, ok := dataMap["planName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.PlanName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.PlanName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.PlanName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.PlanName = types.StringValue(string(jsonBytes))
            } else {
                data.PlanName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.PlanName = types.StringValue(string(jsonBytes))
            } else {
                data.PlanName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.PlanName = types.StringValue(string(jsonBytes))
        } else {
            data.PlanName = types.StringNull()
        }
    } else if val, ok := dataMap["planName"].(string); ok {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
    }
    if obj, ok := dataMap["dataResidency"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DataResidency = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DataResidency = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DataResidency = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DataResidency = types.StringValue(string(jsonBytes))
            } else {
                data.DataResidency = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DataResidency = types.StringValue(string(jsonBytes))
            } else {
                data.DataResidency = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DataResidency = types.StringValue(string(jsonBytes))
        } else {
            data.DataResidency = types.StringNull()
        }
    } else if val, ok := dataMap["dataResidency"].(string); ok {
        data.DataResidency = types.StringValue(val)
    } else {
        data.DataResidency = types.StringNull()
    }
    if obj, ok := dataMap["resellerId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResellerId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ResellerId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ResellerId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ResellerId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ResellerId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ResellerId = types.StringValue(string(jsonBytes))
        } else {
            data.ResellerId = types.StringNull()
        }
    } else if val, ok := dataMap["resellerId"].(string); ok {
        data.ResellerId = types.StringValue(val)
    } else {
        data.ResellerId = types.StringNull()
    }
    if obj, ok := dataMap["resellerPlanId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ResellerPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.ResellerPlanId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.ResellerPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.ResellerPlanId = types.StringValue(string(jsonBytes))
            } else {
                data.ResellerPlanId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.ResellerPlanId = types.StringValue(string(jsonBytes))
        } else {
            data.ResellerPlanId = types.StringNull()
        }
    } else if val, ok := dataMap["resellerPlanId"].(string); ok {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := dataMap["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    }
    if obj, ok := dataMap["gitHubAppInstallationId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.GitHubAppInstallationId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
            } else {
                data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
            } else {
                data.GitHubAppInstallationId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.GitHubAppInstallationId = types.StringValue(string(jsonBytes))
        } else {
            data.GitHubAppInstallationId = types.StringNull()
        }
    } else if val, ok := dataMap["gitHubAppInstallationId"].(string); ok {
        data.GitHubAppInstallationId = types.StringValue(val)
    } else {
        data.GitHubAppInstallationId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    data.Id = state.Id

    // Unconfigured attributes keep their planned value; see keepPlannedValues.
    r.keepPlannedValues(&data, &plan, &config)

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data ProjectResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    httpResp, err := r.client.Delete(ctx, "/project/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete project, got error: %s", err))
        return
    }

    // A failed delete must keep the resource in state — silently dropping it
    // orphans real infrastructure. 404 means it is already gone.
    if httpResp.StatusCode >= 400 && httpResp.StatusCode != http.StatusNotFound {
        err = r.client.ParseResponse(httpResp, nil)
        resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to delete project: %s", err))
        return
    }
    if httpResp.Body != nil {
        httpResp.Body.Close()
    }
}


func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}


// keepPlannedValues puts back, after a create or an update, the planned value
// of each optional attribute the configuration leaves out. The server keeps
// some of these up to date on its own (when it last checked a heartbeat, the
// status a probe last reported...), so the value read back after the write can
// already differ from the plan, and Terraform would fail the apply with
// "Provider produced inconsistent result after apply". The next refresh reads
// the server's value, which is never a diff for an attribute nobody configured.
func (r *ProjectResource) keepPlannedValues(data *ProjectResourceModel, plan *ProjectResourceModel, config *ProjectResourceModel) {
    if config.PaymentProviderPlanId.IsNull() && !plan.PaymentProviderPlanId.IsUnknown() {
        data.PaymentProviderPlanId = plan.PaymentProviderPlanId
    }
    if config.BusinessDetails.IsNull() && !plan.BusinessDetails.IsUnknown() {
        data.BusinessDetails = plan.BusinessDetails
    }
    if config.BusinessDetailsCountry.IsNull() && !plan.BusinessDetailsCountry.IsUnknown() {
        data.BusinessDetailsCountry = plan.BusinessDetailsCountry
    }
    if config.FinanceAccountingEmail.IsNull() && !plan.FinanceAccountingEmail.IsUnknown() {
        data.FinanceAccountingEmail = plan.FinanceAccountingEmail
    }
    if config.PaymentProviderPromoCode.IsNull() && !plan.PaymentProviderPromoCode.IsUnknown() {
        data.PaymentProviderPromoCode = plan.PaymentProviderPromoCode
    }
    if config.IsFeatureFlagMonitorGroupsEnabled.IsNull() && !plan.IsFeatureFlagMonitorGroupsEnabled.IsUnknown() {
        data.IsFeatureFlagMonitorGroupsEnabled = plan.IsFeatureFlagMonitorGroupsEnabled
    }
    if config.IncidentNumberPrefix.IsNull() && !plan.IncidentNumberPrefix.IsUnknown() {
        data.IncidentNumberPrefix = plan.IncidentNumberPrefix
    }
    if config.AlertNumberPrefix.IsNull() && !plan.AlertNumberPrefix.IsUnknown() {
        data.AlertNumberPrefix = plan.AlertNumberPrefix
    }
    if config.ScheduledMaintenanceNumberPrefix.IsNull() && !plan.ScheduledMaintenanceNumberPrefix.IsUnknown() {
        data.ScheduledMaintenanceNumberPrefix = plan.ScheduledMaintenanceNumberPrefix
    }
    if config.IncidentEpisodeNumberPrefix.IsNull() && !plan.IncidentEpisodeNumberPrefix.IsUnknown() {
        data.IncidentEpisodeNumberPrefix = plan.IncidentEpisodeNumberPrefix
    }
    if config.AlertEpisodeNumberPrefix.IsNull() && !plan.AlertEpisodeNumberPrefix.IsUnknown() {
        data.AlertEpisodeNumberPrefix = plan.AlertEpisodeNumberPrefix
    }
    if config.EnableAutomaticIncidentInvestigation.IsNull() && !plan.EnableAutomaticIncidentInvestigation.IsUnknown() {
        data.EnableAutomaticIncidentInvestigation = plan.EnableAutomaticIncidentInvestigation
    }
    if config.EnableAutomaticAlertInvestigation.IsNull() && !plan.EnableAutomaticAlertInvestigation.IsUnknown() {
        data.EnableAutomaticAlertInvestigation = plan.EnableAutomaticAlertInvestigation
    }
    if config.EnableAutomaticIncidentRemediation.IsNull() && !plan.EnableAutomaticIncidentRemediation.IsUnknown() {
        data.EnableAutomaticIncidentRemediation = plan.EnableAutomaticIncidentRemediation
    }
    if config.EnableAutomaticAlertRemediation.IsNull() && !plan.EnableAutomaticAlertRemediation.IsUnknown() {
        data.EnableAutomaticAlertRemediation = plan.EnableAutomaticAlertRemediation
    }
    if config.EnableAutomaticPostmortemDraft.IsNull() && !plan.EnableAutomaticPostmortemDraft.IsUnknown() {
        data.EnableAutomaticPostmortemDraft = plan.EnableAutomaticPostmortemDraft
    }
    if config.EnableIncidentInstrumentationFixTasks.IsNull() && !plan.EnableIncidentInstrumentationFixTasks.IsUnknown() {
        data.EnableIncidentInstrumentationFixTasks = plan.EnableIncidentInstrumentationFixTasks
    }
    if config.EnableAlertInstrumentationFixTasks.IsNull() && !plan.EnableAlertInstrumentationFixTasks.IsUnknown() {
        data.EnableAlertInstrumentationFixTasks = plan.EnableAlertInstrumentationFixTasks
    }
    if config.EnableAutomaticIncidentCodeFixes.IsNull() && !plan.EnableAutomaticIncidentCodeFixes.IsUnknown() {
        data.EnableAutomaticIncidentCodeFixes = plan.EnableAutomaticIncidentCodeFixes
    }
    if config.EnableAutomaticAlertCodeFixes.IsNull() && !plan.EnableAutomaticAlertCodeFixes.IsUnknown() {
        data.EnableAutomaticAlertCodeFixes = plan.EnableAutomaticAlertCodeFixes
    }
    if config.EnableAiInsights.IsNull() && !plan.EnableAiInsights.IsUnknown() {
        data.EnableAiInsights = plan.EnableAiInsights
    }
    if config.EnableInsightFixTasks.IsNull() && !plan.EnableInsightFixTasks.IsUnknown() {
        data.EnableInsightFixTasks = plan.EnableInsightFixTasks
    }
    if config.AutoArchiveNonActionableExceptions.IsNull() && !plan.AutoArchiveNonActionableExceptions.IsUnknown() {
        data.AutoArchiveNonActionableExceptions = plan.AutoArchiveNonActionableExceptions
    }
    if config.SendInvoicesByEmail.IsNull() && !plan.SendInvoicesByEmail.IsUnknown() {
        data.SendInvoicesByEmail = plan.SendInvoicesByEmail
    }
    if config.EnableAuditLogs.IsNull() && !plan.EnableAuditLogs.IsUnknown() {
        data.EnableAuditLogs = plan.EnableAuditLogs
    }
    if config.IsSessionReplayAllowed.IsNull() && !plan.IsSessionReplayAllowed.IsUnknown() {
        data.IsSessionReplayAllowed = plan.IsSessionReplayAllowed
    }
    if config.AuditLogsRetentionInDays.IsNull() && !plan.AuditLogsRetentionInDays.IsUnknown() {
        data.AuditLogsRetentionInDays = plan.AuditLogsRetentionInDays
    }
    if config.StoreSystemEventsInAuditLogs.IsNull() && !plan.StoreSystemEventsInAuditLogs.IsUnknown() {
        data.StoreSystemEventsInAuditLogs = plan.StoreSystemEventsInAuditLogs
    }
    if config.RequireSsoForLogin.IsNull() && !plan.RequireSsoForLogin.IsUnknown() {
        data.RequireSsoForLogin = plan.RequireSsoForLogin
    }
    if config.RequireSsoWithSsoProviderId.IsNull() && !plan.RequireSsoWithSsoProviderId.IsUnknown() {
        data.RequireSsoWithSsoProviderId = plan.RequireSsoWithSsoProviderId
    }
    if config.AutoRechargeSmsOrCallByBalanceInUsd.IsNull() && !plan.AutoRechargeSmsOrCallByBalanceInUsd.IsUnknown() {
        data.AutoRechargeSmsOrCallByBalanceInUsd = plan.AutoRechargeSmsOrCallByBalanceInUsd
    }
    if config.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.IsNull() && !plan.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.IsUnknown() {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = plan.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd
    }
    if config.EnableSmsNotifications.IsNull() && !plan.EnableSmsNotifications.IsUnknown() {
        data.EnableSmsNotifications = plan.EnableSmsNotifications
    }
    if config.EnableWhatsAppNotifications.IsNull() && !plan.EnableWhatsAppNotifications.IsUnknown() {
        data.EnableWhatsAppNotifications = plan.EnableWhatsAppNotifications
    }
    if config.EnableTelegramNotifications.IsNull() && !plan.EnableTelegramNotifications.IsUnknown() {
        data.EnableTelegramNotifications = plan.EnableTelegramNotifications
    }
    if config.EnableCallNotifications.IsNull() && !plan.EnableCallNotifications.IsUnknown() {
        data.EnableCallNotifications = plan.EnableCallNotifications
    }
    if config.DisableOnCallNotificationFallback.IsNull() && !plan.DisableOnCallNotificationFallback.IsUnknown() {
        data.DisableOnCallNotificationFallback = plan.DisableOnCallNotificationFallback
    }
    if config.EnableAutoRechargeSmsOrCallBalance.IsNull() && !plan.EnableAutoRechargeSmsOrCallBalance.IsUnknown() {
        data.EnableAutoRechargeSmsOrCallBalance = plan.EnableAutoRechargeSmsOrCallBalance
    }
    if config.AutoAiRechargeByBalanceInUsd.IsNull() && !plan.AutoAiRechargeByBalanceInUsd.IsUnknown() {
        data.AutoAiRechargeByBalanceInUsd = plan.AutoAiRechargeByBalanceInUsd
    }
    if config.AutoRechargeAiWhenCurrentBalanceFallsInUsd.IsNull() && !plan.AutoRechargeAiWhenCurrentBalanceFallsInUsd.IsUnknown() {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = plan.AutoRechargeAiWhenCurrentBalanceFallsInUsd
    }
    if config.EnableAi.IsNull() && !plan.EnableAi.IsUnknown() {
        data.EnableAi = plan.EnableAi
    }
    if config.AiDailyTokenLimit.IsNull() && !plan.AiDailyTokenLimit.IsUnknown() {
        data.AiDailyTokenLimit = plan.AiDailyTokenLimit
    }
    if config.AiDailySpendLimitInUsd.IsNull() && !plan.AiDailySpendLimitInUsd.IsUnknown() {
        data.AiDailySpendLimitInUsd = plan.AiDailySpendLimitInUsd
    }
    if config.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.IsNull() && !plan.AcknowledgeLinkedAlertsWhenIncidentAcknowledged.IsUnknown() {
        data.AcknowledgeLinkedAlertsWhenIncidentAcknowledged = plan.AcknowledgeLinkedAlertsWhenIncidentAcknowledged
    }
    if config.ResolveLinkedAlertsWhenIncidentResolved.IsNull() && !plan.ResolveLinkedAlertsWhenIncidentResolved.IsUnknown() {
        data.ResolveLinkedAlertsWhenIncidentResolved = plan.ResolveLinkedAlertsWhenIncidentResolved
    }
    if config.AlertInvestigationMinimumSeverityId.IsNull() && !plan.AlertInvestigationMinimumSeverityId.IsUnknown() {
        data.AlertInvestigationMinimumSeverityId = plan.AlertInvestigationMinimumSeverityId
    }
    if config.AiDailyAutonomousTokenLimit.IsNull() && !plan.AiDailyAutonomousTokenLimit.IsUnknown() {
        data.AiDailyAutonomousTokenLimit = plan.AiDailyAutonomousTokenLimit
    }
    if config.IncidentAiDailyAutonomousTokenLimit.IsNull() && !plan.IncidentAiDailyAutonomousTokenLimit.IsUnknown() {
        data.IncidentAiDailyAutonomousTokenLimit = plan.IncidentAiDailyAutonomousTokenLimit
    }
    if config.AlertAiDailyAutonomousTokenLimit.IsNull() && !plan.AlertAiDailyAutonomousTokenLimit.IsUnknown() {
        data.AlertAiDailyAutonomousTokenLimit = plan.AlertAiDailyAutonomousTokenLimit
    }
    if config.AiDailyFixTaskLimit.IsNull() && !plan.AiDailyFixTaskLimit.IsUnknown() {
        data.AiDailyFixTaskLimit = plan.AiDailyFixTaskLimit
    }
    if config.IncidentAiDailyFixTaskLimit.IsNull() && !plan.IncidentAiDailyFixTaskLimit.IsUnknown() {
        data.IncidentAiDailyFixTaskLimit = plan.IncidentAiDailyFixTaskLimit
    }
    if config.AlertAiDailyFixTaskLimit.IsNull() && !plan.AlertAiDailyFixTaskLimit.IsUnknown() {
        data.AlertAiDailyFixTaskLimit = plan.AlertAiDailyFixTaskLimit
    }
    if config.AlertInvestigationDedupeWindowMinutes.IsNull() && !plan.AlertInvestigationDedupeWindowMinutes.IsUnknown() {
        data.AlertInvestigationDedupeWindowMinutes = plan.AlertInvestigationDedupeWindowMinutes
    }
    if config.IncidentInvestigationMinimumSeverityId.IsNull() && !plan.IncidentInvestigationMinimumSeverityId.IsUnknown() {
        data.IncidentInvestigationMinimumSeverityId = plan.IncidentInvestigationMinimumSeverityId
    }
    if config.IncidentInvestigationDedupeWindowMinutes.IsNull() && !plan.IncidentInvestigationDedupeWindowMinutes.IsUnknown() {
        data.IncidentInvestigationDedupeWindowMinutes = plan.IncidentInvestigationDedupeWindowMinutes
    }
    if config.AiMaxConcurrentInvestigations.IsNull() && !plan.AiMaxConcurrentInvestigations.IsUnknown() {
        data.AiMaxConcurrentInvestigations = plan.AiMaxConcurrentInvestigations
    }
    if config.IncidentAiMaxConcurrentInvestigations.IsNull() && !plan.IncidentAiMaxConcurrentInvestigations.IsUnknown() {
        data.IncidentAiMaxConcurrentInvestigations = plan.IncidentAiMaxConcurrentInvestigations
    }
    if config.AlertAiMaxConcurrentInvestigations.IsNull() && !plan.AlertAiMaxConcurrentInvestigations.IsUnknown() {
        data.AlertAiMaxConcurrentInvestigations = plan.AlertAiMaxConcurrentInvestigations
    }
    if config.IncidentAiInvestigationTimeLimitInMinutes.IsNull() && !plan.IncidentAiInvestigationTimeLimitInMinutes.IsUnknown() {
        data.IncidentAiInvestigationTimeLimitInMinutes = plan.IncidentAiInvestigationTimeLimitInMinutes
    }
    if config.AlertAiInvestigationTimeLimitInMinutes.IsNull() && !plan.AlertAiInvestigationTimeLimitInMinutes.IsUnknown() {
        data.AlertAiInvestigationTimeLimitInMinutes = plan.AlertAiInvestigationTimeLimitInMinutes
    }
    if config.EnableAutoRechargeAiBalance.IsNull() && !plan.EnableAutoRechargeAiBalance.IsUnknown() {
        data.EnableAutoRechargeAiBalance = plan.EnableAutoRechargeAiBalance
    }
    if config.DoNotAddGlobalProbesByDefaultOnNewMonitors.IsNull() && !plan.DoNotAddGlobalProbesByDefaultOnNewMonitors.IsUnknown() {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = plan.DoNotAddGlobalProbesByDefaultOnNewMonitors
    }
    if config.DefaultMetricCardinalityBudget.IsNull() && !plan.DefaultMetricCardinalityBudget.IsUnknown() {
        data.DefaultMetricCardinalityBudget = plan.DefaultMetricCardinalityBudget
    }
    if config.DefaultTelemetryRetentionInDays.IsNull() && !plan.DefaultTelemetryRetentionInDays.IsUnknown() {
        data.DefaultTelemetryRetentionInDays = plan.DefaultTelemetryRetentionInDays
    }
    if config.TelemetryRetentionConfig.IsNull() && !plan.TelemetryRetentionConfig.IsUnknown() {
        data.TelemetryRetentionConfig = plan.TelemetryRetentionConfig
    }
    if config.DefaultMetricDownsamplingRetentionDays.IsNull() && !plan.DefaultMetricDownsamplingRetentionDays.IsUnknown() {
        data.DefaultMetricDownsamplingRetentionDays = plan.DefaultMetricDownsamplingRetentionDays
    }
}

// Helper method to convert Terraform map to Go interface{}
func (r *ProjectResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
    if terraformMap.IsNull() || terraformMap.IsUnknown() {
        return nil
    }
    
    result := make(map[string]string)
    terraformMap.ElementsAs(context.Background(), &result, false)
    
    // Convert map[string]string to map[string]interface{}
    interfaceResult := make(map[string]interface{})
    for key, value := range result {
        interfaceResult[key] = value
    }
    
    return interfaceResult
}

// Helper method to convert Terraform list to Go interface{}
func (r *ProjectResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformList.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}

// Helper method to convert Terraform set to Go interface{}
func (r *ProjectResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
    if terraformSet.IsNull() || terraformSet.IsUnknown() {
        return nil
    }

    var stringList []string
    terraformSet.ElementsAs(context.Background(), &stringList, false)

    // Convert string array to OneUptime format with _id fields
    var result []interface{}
    for _, str := range stringList {
        if str != "" {
            result = append(result, map[string]interface{}{
                "_id": str,
            })
        }
    }
    return result
}


// Helper method to parse JSON field for complex objects
func (r *ProjectResource) parseJSONField(terraformString basetypes.StringValuable) interface{} {
    sv, _ := terraformString.ToStringValue(context.Background())
    if sv.IsNull() || sv.IsUnknown() || sv.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(sv.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return sv.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *ProjectResource) normalizeURLWrappers(value interface{}) interface{} {
    switch v := value.(type) {
    case map[string]interface{}:
        if typeStr, ok := v["_type"].(string); ok && typeStr == "URL" {
            if val, ok := v["value"].(string); ok {
                v["value"] = r.normalizeURLString(val)
            }
        }
        for key, child := range v {
            v[key] = r.normalizeURLWrappers(child)
        }
        return v
    case []interface{}:
        for i, child := range v {
            v[i] = r.normalizeURLWrappers(child)
        }
        return v
    default:
        return v
    }
}

func (r *ProjectResource) normalizeURLString(value string) string {
    parsed, err := url.Parse(value)
    if err != nil {
        return value
    }
    if parsed.Path == "/" && parsed.RawQuery == "" && parsed.Fragment == "" {
        return strings.TrimSuffix(value, "/")
    }
    return value
}

// Helper method to convert *big.Float to float64 for JSON serialization
func (r *ProjectResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType.
// The registry itself lives in objecttypes.go, shared across the package.
func (r *ProjectResource) isValidOneUptimeObjectType(typeStr string) bool {
    return validOneUptimeObjectTypes[typeStr]
}
