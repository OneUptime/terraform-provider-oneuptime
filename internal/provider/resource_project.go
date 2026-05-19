package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
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
    SendInvoicesByEmail types.Bool `tfsdk:"send_invoices_by_email"`
    UtmContent types.String `tfsdk:"utm_content"`
    EnableAuditLogs types.Bool `tfsdk:"enable_audit_logs"`
    AuditLogsRetentionInDays types.Number `tfsdk:"audit_logs_retention_in_days"`
    StoreSystemEventsInAuditLogs types.Bool `tfsdk:"store_system_events_in_audit_logs"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    AutoRechargeSmsOrCallByBalanceInUsd types.Number `tfsdk:"auto_recharge_sms_or_call_by_balance_in_usd"`
    AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd types.Number `tfsdk:"auto_recharge_sms_or_call_when_current_balance_falls_in_usd"`
    EnableSmsNotifications types.Bool `tfsdk:"enable_sms_notifications"`
    EnableWhatsAppNotifications types.Bool `tfsdk:"enable_whats_app_notifications"`
    EnableTelegramNotifications types.Bool `tfsdk:"enable_telegram_notifications"`
    EnableCallNotifications types.Bool `tfsdk:"enable_call_notifications"`
    EnableAutoRechargeSmsOrCallBalance types.Bool `tfsdk:"enable_auto_recharge_sms_or_call_balance"`
    AutoAiRechargeByBalanceInUsd types.Number `tfsdk:"auto_ai_recharge_by_balance_in_usd"`
    AutoRechargeAiWhenCurrentBalanceFallsInUsd types.Number `tfsdk:"auto_recharge_ai_when_current_balance_falls_in_usd"`
    EnableAi types.Bool `tfsdk:"enable_ai"`
    EnableAutoRechargeAiBalance types.Bool `tfsdk:"enable_auto_recharge_ai_balance"`
    DoNotAddGlobalProbesByDefaultOnNewMonitors types.Bool `tfsdk:"do_not_add_global_probes_by_default_on_new_monitors"`
    GitHubAppInstallationId types.String `tfsdk:"git_hub_app_installation_id"`
    DefaultMetricCardinalityBudget types.Number `tfsdk:"default_metric_cardinality_budget"`
    DefaultTelemetryRetentionInDays types.Number `tfsdk:"default_telemetry_retention_in_days"`
    TelemetryRetentionConfig JSONSubsetValue `tfsdk:"telemetry_retention_config"`
    DefaultMetricDownsamplingRetentionDays JSONSubsetValue `tfsdk:"default_metric_downsampling_retention_days"`
    CreatedAt JSONSubsetValue `tfsdk:"created_at"`
    UpdatedAt JSONSubsetValue `tfsdk:"updated_at"`
    DeletedAt JSONSubsetValue `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    PaymentProviderSubscriptionId types.String `tfsdk:"payment_provider_subscription_id"`
    PaymentProviderMeteredSubscriptionId types.String `tfsdk:"payment_provider_metered_subscription_id"`
    PaymentProviderSubscriptionSeats types.Number `tfsdk:"payment_provider_subscription_seats"`
    TrialEndsAt JSONSubsetValue `tfsdk:"trial_ends_at"`
    PaymentProviderCustomerId types.String `tfsdk:"payment_provider_customer_id"`
    PaymentProviderSubscriptionStatus types.String `tfsdk:"payment_provider_subscription_status"`
    PaymentProviderMeteredSubscriptionStatus types.String `tfsdk:"payment_provider_metered_subscription_status"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    WorkflowRunsInLast30Days types.Number `tfsdk:"workflow_runs_in_last30_days"`
    SmsOrCallCurrentBalanceInUsdCents types.Number `tfsdk:"sms_or_call_current_balance_in_usd_cents"`
    AiCurrentBalanceInUsdCents types.Number `tfsdk:"ai_current_balance_in_usd_cents"`
    PlanName types.String `tfsdk:"plan_name"`
    ResellerId types.String `tfsdk:"reseller_id"`
    ResellerPlanId types.String `tfsdk:"reseller_plan_id"`
    LetCustomerSupportAccessProject types.Bool `tfsdk:"let_customer_support_access_project"`
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "project resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Any friendly name of this object. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing, Edit Project]",
                Required: true,
            },
            "payment_provider_plan_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "business_details": schema.StringAttribute{
                MarkdownDescription: "Business legal name, address and any tax information to appear on invoices.. Permissions - Create: [Project Owner, Manage Billing], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "business_details_country": schema.StringAttribute{
                MarkdownDescription: "Two-letter ISO country code for billing address (e.g., US, GB, DE).. Permissions - Create: [Project Owner, Manage Billing], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "finance_accounting_email": schema.StringAttribute{
                MarkdownDescription: "Invoices, receipts and billing related notifications will be sent to these emails in addition to project owner. Separate multiple emails with a comma.. Permissions - Create: [Project Owner, Manage Billing], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "payment_provider_promo_code": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_feature_flag_monitor_groups_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Feature Flag Monitor Groups Enabled. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing, Edit Project]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "active_monitors_limit": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [No access - you don't have permission for this operation], Update: [No access - you don't have permission for this operation]",
                Optional: true,
            },
            "seat_limit": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [No access - you don't have permission for this operation], Update: [No access - you don't have permission for this operation]",
                Optional: true,
            },
            "incident_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident numbers (e.g., 'INC-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert numbers (e.g., 'ALT-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "scheduled_maintenance_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for scheduled maintenance numbers (e.g., 'SM-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "incident_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident episode numbers (e.g., 'IE-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "alert_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert episode numbers (e.g., 'AE-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "send_invoices_by_email": schema.BoolAttribute{
                MarkdownDescription: "When enabled, invoices will be automatically sent to the finance/accounting email when they are generated.. Permissions - Create: [Project Owner, Manage Billing], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "utm_content": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [No access - you don't have permission for this operation], Update: [No access - you don't have permission for this operation]",
                Optional: true,
            },
            "enable_audit_logs": schema.BoolAttribute{
                MarkdownDescription: "When enabled, changes to resources in this project are recorded as audit log entries.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "audit_logs_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "Number of days to retain audit log entries. Minimum 7, maximum 180.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(7)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "store_system_events_in_audit_logs": schema.BoolAttribute{
                MarkdownDescription: "When enabled, audit logs will also include events triggered by the system. By default, only events triggered by users are recorded.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_recharge_sms_or_call_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for SMS, Call, and WhatsApp. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(20)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_recharge_sms_or_call_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for SMS, Call, and WhatsApp. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(10)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_sms_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable SMS notifications for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_whats_app_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable WhatsApp notifications for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_telegram_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable Telegram notifications for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_call_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable call notifications for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_auto_recharge_sms_or_call_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for SMS, Call, and WhatsApp balance for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_ai_recharge_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for AI services. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(20)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "auto_recharge_ai_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for AI services. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: numberdefault.StaticBigFloat(big.NewFloat(10)),
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_ai": schema.BoolAttribute{
                MarkdownDescription: "Enable AI services for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "enable_auto_recharge_ai_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for AI balance for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "do_not_add_global_probes_by_default_on_new_monitors": schema.BoolAttribute{
                MarkdownDescription: "If enabled, global probes will NOT be automatically added to new monitors. Enable this only if you are using ONLY custom probes to monitor your resources.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "git_hub_app_installation_id": schema.StringAttribute{
                MarkdownDescription: "The GitHub App installation ID for this project. This is set when the GitHub App is installed on the organization.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "default_metric_cardinality_budget": schema.NumberAttribute{
                MarkdownDescription: "Project-wide default max distinct series per metric. Services without a per-service override use this value.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "default_telemetry_retention_in_days": schema.NumberAttribute{
                MarkdownDescription: "Project-wide default number of days to retain telemetry data (logs, traces, metrics). Services without a per-service override use this value.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.Number{
                    numberplanmodifier.UseStateForUnknown(),
                },
            },
            "telemetry_retention_config": schema.StringAttribute{
                MarkdownDescription: "Project-wide per-pillar retention overrides for telemetry data (logs by severity, traces by status, metrics, profiles). Falls back to defaultTelemetryRetentionInDays when a pillar or bucket is not set.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin]",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "default_metric_downsampling_retention_days": schema.StringAttribute{
                MarkdownDescription: "Project-wide default retention for each downsampling tier (raw, 1m, 5m, 1h, 1d) in days.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin]",
                CustomType: JSONSubsetType{},
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "updated_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "deleted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Object version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_metered_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_subscription_seats": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "trial_ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                CustomType: JSONSubsetType{},
                Computed: true,
            },
            "payment_provider_customer_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_metered_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
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
            "workflow_runs_in_last30_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Read Workflow], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "sms_or_call_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for SMS, Call, and WhatsApp. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "ai_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for AI services. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "plan_name": schema.StringAttribute{
                MarkdownDescription: "Name of the plan this project is subscribed to.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "OneUptime customer support can access this project. This is used for debugging purposes.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Viewer, Read Project], Update: [Project Owner, Project Admin]",
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



    // Create API request body
    projectRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "name": data.Name.ValueString(),
        "paymentProviderPlanId": data.PaymentProviderPlanId.ValueString(),
        "businessDetails": data.BusinessDetails.ValueString(),
        "businessDetailsCountry": data.BusinessDetailsCountry.ValueString(),
        "financeAccountingEmail": data.FinanceAccountingEmail.ValueString(),
        "paymentProviderPromoCode": data.PaymentProviderPromoCode.ValueString(),
        "isFeatureFlagMonitorGroupsEnabled": data.IsFeatureFlagMonitorGroupsEnabled.ValueBool(),
        "activeMonitorsLimit": r.bigFloatToFloat64(data.ActiveMonitorsLimit.ValueBigFloat()),
        "seatLimit": r.bigFloatToFloat64(data.SeatLimit.ValueBigFloat()),
        "incidentNumberPrefix": data.IncidentNumberPrefix.ValueString(),
        "alertNumberPrefix": data.AlertNumberPrefix.ValueString(),
        "scheduledMaintenanceNumberPrefix": data.ScheduledMaintenanceNumberPrefix.ValueString(),
        "incidentEpisodeNumberPrefix": data.IncidentEpisodeNumberPrefix.ValueString(),
        "alertEpisodeNumberPrefix": data.AlertEpisodeNumberPrefix.ValueString(),
        "sendInvoicesByEmail": data.SendInvoicesByEmail.ValueBool(),
        "utmContent": data.UtmContent.ValueString(),
        "enableAuditLogs": data.EnableAuditLogs.ValueBool(),
        "auditLogsRetentionInDays": r.bigFloatToFloat64(data.AuditLogsRetentionInDays.ValueBigFloat()),
        "storeSystemEventsInAuditLogs": data.StoreSystemEventsInAuditLogs.ValueBool(),
        "requireSsoForLogin": data.RequireSsoForLogin.ValueBool(),
        "autoRechargeSmsOrCallByBalanceInUSD": r.bigFloatToFloat64(data.AutoRechargeSmsOrCallByBalanceInUsd.ValueBigFloat()),
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": r.bigFloatToFloat64(data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd.ValueBigFloat()),
        "enableSmsNotifications": data.EnableSmsNotifications.ValueBool(),
        "enableWhatsAppNotifications": data.EnableWhatsAppNotifications.ValueBool(),
        "enableTelegramNotifications": data.EnableTelegramNotifications.ValueBool(),
        "enableCallNotifications": data.EnableCallNotifications.ValueBool(),
        "enableAutoRechargeSmsOrCallBalance": data.EnableAutoRechargeSmsOrCallBalance.ValueBool(),
        "autoAiRechargeByBalanceInUSD": r.bigFloatToFloat64(data.AutoAiRechargeByBalanceInUsd.ValueBigFloat()),
        "autoRechargeAiWhenCurrentBalanceFallsInUSD": r.bigFloatToFloat64(data.AutoRechargeAiWhenCurrentBalanceFallsInUsd.ValueBigFloat()),
        "enableAi": data.EnableAi.ValueBool(),
        "enableAutoRechargeAiBalance": data.EnableAutoRechargeAiBalance.ValueBool(),
        "doNotAddGlobalProbesByDefaultOnNewMonitors": data.DoNotAddGlobalProbesByDefaultOnNewMonitors.ValueBool(),
        "gitHubAppInstallationId": data.GitHubAppInstallationId.ValueString(),
        "defaultMetricCardinalityBudget": r.bigFloatToFloat64(data.DefaultMetricCardinalityBudget.ValueBigFloat()),
        "defaultTelemetryRetentionInDays": r.bigFloatToFloat64(data.DefaultTelemetryRetentionInDays.ValueBigFloat()),
        "telemetryRetentionConfig": r.parseJSONField(data.TelemetryRetentionConfig),
        "defaultMetricDownsamplingRetentionDays": r.parseJSONField(data.DefaultMetricDownsamplingRetentionDays),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/project", projectRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create project, got error: %s", err))
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

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderPlanId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["businessDetails"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["businessDetailsCountry"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["financeAccountingEmail"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderPromoCode"].(string); ok && val != "" {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if val, ok := dataMap["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["activeMonitorsLimit"].(float64); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["activeMonitorsLimit"].(int); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["activeMonitorsLimit"].(int64); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["activeMonitorsLimit"] == nil {
        data.ActiveMonitorsLimit = types.NumberNull()
    }
    if val, ok := dataMap["seatLimit"].(float64); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["seatLimit"].(int); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["seatLimit"].(int64); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["seatLimit"] == nil {
        data.SeatLimit = types.NumberNull()
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
    } else if val, ok := dataMap["incidentNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["alertNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["scheduledMaintenanceNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["incidentEpisodeNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["alertEpisodeNumberPrefix"].(string); ok && val != "" {
        data.AlertEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.AlertEpisodeNumberPrefix = types.StringNull()
    }
    if val, ok := dataMap["sendInvoicesByEmail"].(bool); ok {
        data.SendInvoicesByEmail = types.BoolValue(val)
    }
    if obj, ok := dataMap["utmContent"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UtmContent = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UtmContent = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UtmContent = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UtmContent = types.StringValue(string(jsonBytes))
            } else {
                data.UtmContent = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UtmContent = types.StringValue(string(jsonBytes))
            } else {
                data.UtmContent = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UtmContent = types.StringValue(string(jsonBytes))
        } else {
            data.UtmContent = types.StringNull()
        }
    } else if val, ok := dataMap["utmContent"].(string); ok && val != "" {
        data.UtmContent = types.StringValue(val)
    } else {
        data.UtmContent = types.StringNull()
    }
    if val, ok := dataMap["enableAuditLogs"].(bool); ok {
        data.EnableAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["auditLogsRetentionInDays"].(float64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["auditLogsRetentionInDays"] == nil {
        data.AuditLogsRetentionInDays = types.NumberNull()
    }
    if val, ok := dataMap["storeSystemEventsInAuditLogs"].(bool); ok {
        data.StoreSystemEventsInAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoRechargeSmsOrCallByBalanceInUSD"] == nil {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"] == nil {
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
    if val, ok := dataMap["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(float64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoAiRechargeByBalanceInUSD"] == nil {
        data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"] == nil {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := dataMap["enableAi"].(bool); ok {
        data.EnableAi = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutoRechargeAiBalance"].(bool); ok {
        data.EnableAutoRechargeAiBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["doNotAddGlobalProbesByDefaultOnNewMonitors"].(bool); ok {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = types.BoolValue(val)
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
    } else if val, ok := dataMap["gitHubAppInstallationId"].(string); ok && val != "" {
        data.GitHubAppInstallationId = types.StringValue(val)
    } else {
        data.GitHubAppInstallationId = types.StringNull()
    }
    if val, ok := dataMap["defaultMetricCardinalityBudget"].(float64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["defaultMetricCardinalityBudget"] == nil {
        data.DefaultMetricCardinalityBudget = types.NumberNull()
    }
    if val, ok := dataMap["defaultTelemetryRetentionInDays"].(float64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["defaultTelemetryRetentionInDays"] == nil {
        data.DefaultTelemetryRetentionInDays = types.NumberNull()
    }
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
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
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok && val != "" {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
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
    } else if val, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(string); ok && val != "" {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
    } else {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CreatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewJSONSubsetValue(val)
    } else {
        data.CreatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.UpdatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewJSONSubsetValue(val)
    } else {
        data.UpdatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DeletedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewJSONSubsetValue(val)
    } else {
        data.DeletedAt = NewJSONSubsetNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
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
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderSubscriptionId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionId"].(string); ok && val != "" {
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
    } else if dataMap["paymentProviderSubscriptionSeats"] == nil {
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if obj, ok := dataMap["trialEndsAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrialEndsAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TrialEndsAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TrialEndsAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TrialEndsAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrialEndsAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TrialEndsAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrialEndsAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TrialEndsAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TrialEndsAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["trialEndsAt"].(string); ok && val != "" {
        data.TrialEndsAt = NewJSONSubsetValue(val)
    } else {
        data.TrialEndsAt = NewJSONSubsetNull()
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
    } else if val, ok := dataMap["paymentProviderCustomerId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderSubscriptionStatus"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["workflowRunsInLast30Days"] == nil {
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["smsOrCallCurrentBalanceInUSDCents"] == nil {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(float64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["aiCurrentBalanceInUSDCents"] == nil {
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
    } else if val, ok := dataMap["planName"].(string); ok && val != "" {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
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
    } else if val, ok := dataMap["resellerId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["resellerPlanId"].(string); ok && val != "" {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := dataMap["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

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
        "activeMonitorsLimit": true,
        "seatLimit": true,
        "incidentNumberPrefix": true,
        "alertNumberPrefix": true,
        "scheduledMaintenanceNumberPrefix": true,
        "incidentEpisodeNumberPrefix": true,
        "alertEpisodeNumberPrefix": true,
        "sendInvoicesByEmail": true,
        "utmContent": true,
        "enableAuditLogs": true,
        "auditLogsRetentionInDays": true,
        "storeSystemEventsInAuditLogs": true,
        "requireSsoForLogin": true,
        "autoRechargeSmsOrCallByBalanceInUSD": true,
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": true,
        "enableSmsNotifications": true,
        "enableWhatsAppNotifications": true,
        "enableTelegramNotifications": true,
        "enableCallNotifications": true,
        "enableAutoRechargeSmsOrCallBalance": true,
        "autoAiRechargeByBalanceInUSD": true,
        "autoRechargeAiWhenCurrentBalanceFallsInUSD": true,
        "enableAi": true,
        "enableAutoRechargeAiBalance": true,
        "doNotAddGlobalProbesByDefaultOnNewMonitors": true,
        "gitHubAppInstallationId": true,
        "defaultMetricCardinalityBudget": true,
        "defaultTelemetryRetentionInDays": true,
        "telemetryRetentionConfig": true,
        "defaultMetricDownsamplingRetentionDays": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "paymentProviderSubscriptionId": true,
        "paymentProviderMeteredSubscriptionId": true,
        "paymentProviderSubscriptionSeats": true,
        "trialEndsAt": true,
        "paymentProviderCustomerId": true,
        "paymentProviderSubscriptionStatus": true,
        "paymentProviderMeteredSubscriptionStatus": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "workflowRunsInLast30Days": true,
        "smsOrCallCurrentBalanceInUSDCents": true,
        "aiCurrentBalanceInUSDCents": true,
        "planName": true,
        "resellerId": true,
        "resellerPlanId": true,
        "letCustomerSupportAccessProject": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/project/" + data.Id.ValueString() + "/get-item", selectParam)
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

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderPlanId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["businessDetails"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["businessDetailsCountry"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["financeAccountingEmail"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderPromoCode"].(string); ok && val != "" {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if val, ok := dataMap["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["activeMonitorsLimit"].(float64); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["activeMonitorsLimit"].(int); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["activeMonitorsLimit"].(int64); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["activeMonitorsLimit"] == nil {
        data.ActiveMonitorsLimit = types.NumberNull()
    }
    if val, ok := dataMap["seatLimit"].(float64); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["seatLimit"].(int); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["seatLimit"].(int64); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["seatLimit"] == nil {
        data.SeatLimit = types.NumberNull()
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
    } else if val, ok := dataMap["incidentNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["alertNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["scheduledMaintenanceNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["incidentEpisodeNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["alertEpisodeNumberPrefix"].(string); ok && val != "" {
        data.AlertEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.AlertEpisodeNumberPrefix = types.StringNull()
    }
    if val, ok := dataMap["sendInvoicesByEmail"].(bool); ok {
        data.SendInvoicesByEmail = types.BoolValue(val)
    }
    if obj, ok := dataMap["utmContent"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UtmContent = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UtmContent = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UtmContent = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UtmContent = types.StringValue(string(jsonBytes))
            } else {
                data.UtmContent = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UtmContent = types.StringValue(string(jsonBytes))
            } else {
                data.UtmContent = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UtmContent = types.StringValue(string(jsonBytes))
        } else {
            data.UtmContent = types.StringNull()
        }
    } else if val, ok := dataMap["utmContent"].(string); ok && val != "" {
        data.UtmContent = types.StringValue(val)
    } else {
        data.UtmContent = types.StringNull()
    }
    if val, ok := dataMap["enableAuditLogs"].(bool); ok {
        data.EnableAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["auditLogsRetentionInDays"].(float64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["auditLogsRetentionInDays"] == nil {
        data.AuditLogsRetentionInDays = types.NumberNull()
    }
    if val, ok := dataMap["storeSystemEventsInAuditLogs"].(bool); ok {
        data.StoreSystemEventsInAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoRechargeSmsOrCallByBalanceInUSD"] == nil {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"] == nil {
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
    if val, ok := dataMap["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(float64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoAiRechargeByBalanceInUSD"] == nil {
        data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"] == nil {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := dataMap["enableAi"].(bool); ok {
        data.EnableAi = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutoRechargeAiBalance"].(bool); ok {
        data.EnableAutoRechargeAiBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["doNotAddGlobalProbesByDefaultOnNewMonitors"].(bool); ok {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = types.BoolValue(val)
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
    } else if val, ok := dataMap["gitHubAppInstallationId"].(string); ok && val != "" {
        data.GitHubAppInstallationId = types.StringValue(val)
    } else {
        data.GitHubAppInstallationId = types.StringNull()
    }
    if val, ok := dataMap["defaultMetricCardinalityBudget"].(float64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["defaultMetricCardinalityBudget"] == nil {
        data.DefaultMetricCardinalityBudget = types.NumberNull()
    }
    if val, ok := dataMap["defaultTelemetryRetentionInDays"].(float64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["defaultTelemetryRetentionInDays"] == nil {
        data.DefaultTelemetryRetentionInDays = types.NumberNull()
    }
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
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
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok && val != "" {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
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
    } else if val, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(string); ok && val != "" {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
    } else {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CreatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewJSONSubsetValue(val)
    } else {
        data.CreatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.UpdatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewJSONSubsetValue(val)
    } else {
        data.UpdatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DeletedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewJSONSubsetValue(val)
    } else {
        data.DeletedAt = NewJSONSubsetNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
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
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderSubscriptionId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionId"].(string); ok && val != "" {
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
    } else if dataMap["paymentProviderSubscriptionSeats"] == nil {
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if obj, ok := dataMap["trialEndsAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrialEndsAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TrialEndsAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TrialEndsAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TrialEndsAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrialEndsAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TrialEndsAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrialEndsAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TrialEndsAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TrialEndsAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["trialEndsAt"].(string); ok && val != "" {
        data.TrialEndsAt = NewJSONSubsetValue(val)
    } else {
        data.TrialEndsAt = NewJSONSubsetNull()
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
    } else if val, ok := dataMap["paymentProviderCustomerId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderSubscriptionStatus"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["workflowRunsInLast30Days"] == nil {
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["smsOrCallCurrentBalanceInUSDCents"] == nil {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(float64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["aiCurrentBalanceInUSDCents"] == nil {
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
    } else if val, ok := dataMap["planName"].(string); ok && val != "" {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
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
    } else if val, ok := dataMap["resellerId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["resellerPlanId"].(string); ok && val != "" {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := dataMap["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
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
    if !data.EnableAutoRechargeAiBalance.IsUnknown() && !state.EnableAutoRechargeAiBalance.IsUnknown() && !data.EnableAutoRechargeAiBalance.Equal(state.EnableAutoRechargeAiBalance) {
        requestDataMap["enableAutoRechargeAiBalance"] = data.EnableAutoRechargeAiBalance.ValueBool()
    }
    if !data.SendInvoicesByEmail.IsUnknown() && !state.SendInvoicesByEmail.IsUnknown() && !data.SendInvoicesByEmail.Equal(state.SendInvoicesByEmail) {
        requestDataMap["sendInvoicesByEmail"] = data.SendInvoicesByEmail.ValueBool()
    }
    if !data.DoNotAddGlobalProbesByDefaultOnNewMonitors.IsUnknown() && !state.DoNotAddGlobalProbesByDefaultOnNewMonitors.IsUnknown() && !data.DoNotAddGlobalProbesByDefaultOnNewMonitors.Equal(state.DoNotAddGlobalProbesByDefaultOnNewMonitors) {
        requestDataMap["doNotAddGlobalProbesByDefaultOnNewMonitors"] = data.DoNotAddGlobalProbesByDefaultOnNewMonitors.ValueBool()
    }
    if !data.GitHubAppInstallationId.IsUnknown() && !state.GitHubAppInstallationId.IsUnknown() && !data.GitHubAppInstallationId.Equal(state.GitHubAppInstallationId) {
        requestDataMap["gitHubAppInstallationId"] = data.GitHubAppInstallationId.ValueString()
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
    if !data.AuditLogsRetentionInDays.IsUnknown() && !state.AuditLogsRetentionInDays.IsUnknown() && !data.AuditLogsRetentionInDays.Equal(state.AuditLogsRetentionInDays) {
        requestDataMap["auditLogsRetentionInDays"] = r.bigFloatToFloat64(data.AuditLogsRetentionInDays.ValueBigFloat())
    }
    if !data.StoreSystemEventsInAuditLogs.IsUnknown() && !state.StoreSystemEventsInAuditLogs.IsUnknown() && !data.StoreSystemEventsInAuditLogs.Equal(state.StoreSystemEventsInAuditLogs) {
        requestDataMap["storeSystemEventsInAuditLogs"] = data.StoreSystemEventsInAuditLogs.ValueBool()
    }

    // Make API call
    httpResp, err := r.client.Put("/project/" + data.Id.ValueString() + "", projectRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update project, got error: %s", err))
        return
    }

    // Parse the update response
    var projectResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &projectResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse project response, got error: %s", err))
        return
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
        "activeMonitorsLimit": true,
        "seatLimit": true,
        "incidentNumberPrefix": true,
        "alertNumberPrefix": true,
        "scheduledMaintenanceNumberPrefix": true,
        "incidentEpisodeNumberPrefix": true,
        "alertEpisodeNumberPrefix": true,
        "sendInvoicesByEmail": true,
        "utmContent": true,
        "enableAuditLogs": true,
        "auditLogsRetentionInDays": true,
        "storeSystemEventsInAuditLogs": true,
        "requireSsoForLogin": true,
        "autoRechargeSmsOrCallByBalanceInUSD": true,
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": true,
        "enableSmsNotifications": true,
        "enableWhatsAppNotifications": true,
        "enableTelegramNotifications": true,
        "enableCallNotifications": true,
        "enableAutoRechargeSmsOrCallBalance": true,
        "autoAiRechargeByBalanceInUSD": true,
        "autoRechargeAiWhenCurrentBalanceFallsInUSD": true,
        "enableAi": true,
        "enableAutoRechargeAiBalance": true,
        "doNotAddGlobalProbesByDefaultOnNewMonitors": true,
        "gitHubAppInstallationId": true,
        "defaultMetricCardinalityBudget": true,
        "defaultTelemetryRetentionInDays": true,
        "telemetryRetentionConfig": true,
        "defaultMetricDownsamplingRetentionDays": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "slug": true,
        "paymentProviderSubscriptionId": true,
        "paymentProviderMeteredSubscriptionId": true,
        "paymentProviderSubscriptionSeats": true,
        "trialEndsAt": true,
        "paymentProviderCustomerId": true,
        "paymentProviderSubscriptionStatus": true,
        "paymentProviderMeteredSubscriptionStatus": true,
        "createdByUserId": true,
        "deletedByUserId": true,
        "workflowRunsInLast30Days": true,
        "smsOrCallCurrentBalanceInUSDCents": true,
        "aiCurrentBalanceInUSDCents": true,
        "planName": true,
        "resellerId": true,
        "resellerPlanId": true,
        "letCustomerSupportAccessProject": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/project/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse project read response, got error: %s", err))
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

    if obj, ok := dataMap["id"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.Id = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.Id = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.Id = types.StringValue(string(jsonBytes))
            } else {
                data.Id = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.Id = types.StringValue(string(jsonBytes))
        } else {
            data.Id = types.StringNull()
        }
    } else if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
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
    } else if val, ok := dataMap["name"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderPlanId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["businessDetails"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["businessDetailsCountry"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["financeAccountingEmail"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderPromoCode"].(string); ok && val != "" {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if val, ok := dataMap["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    }
    if val, ok := dataMap["activeMonitorsLimit"].(float64); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["activeMonitorsLimit"].(int); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["activeMonitorsLimit"].(int64); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["activeMonitorsLimit"] == nil {
        data.ActiveMonitorsLimit = types.NumberNull()
    }
    if val, ok := dataMap["seatLimit"].(float64); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["seatLimit"].(int); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["seatLimit"].(int64); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["seatLimit"] == nil {
        data.SeatLimit = types.NumberNull()
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
    } else if val, ok := dataMap["incidentNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["alertNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["scheduledMaintenanceNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["incidentEpisodeNumberPrefix"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["alertEpisodeNumberPrefix"].(string); ok && val != "" {
        data.AlertEpisodeNumberPrefix = types.StringValue(val)
    } else {
        data.AlertEpisodeNumberPrefix = types.StringNull()
    }
    if val, ok := dataMap["sendInvoicesByEmail"].(bool); ok {
        data.SendInvoicesByEmail = types.BoolValue(val)
    }
    if obj, ok := dataMap["utmContent"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UtmContent = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UtmContent = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UtmContent = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UtmContent = types.StringValue(string(jsonBytes))
            } else {
                data.UtmContent = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UtmContent = types.StringValue(string(jsonBytes))
            } else {
                data.UtmContent = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UtmContent = types.StringValue(string(jsonBytes))
        } else {
            data.UtmContent = types.StringNull()
        }
    } else if val, ok := dataMap["utmContent"].(string); ok && val != "" {
        data.UtmContent = types.StringValue(val)
    } else {
        data.UtmContent = types.StringNull()
    }
    if val, ok := dataMap["enableAuditLogs"].(bool); ok {
        data.EnableAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["auditLogsRetentionInDays"].(float64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["auditLogsRetentionInDays"].(int64); ok {
        data.AuditLogsRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["auditLogsRetentionInDays"] == nil {
        data.AuditLogsRetentionInDays = types.NumberNull()
    }
    if val, ok := dataMap["storeSystemEventsInAuditLogs"].(bool); ok {
        data.StoreSystemEventsInAuditLogs = types.BoolValue(val)
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoRechargeSmsOrCallByBalanceInUSD"] == nil {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"] == nil {
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
    if val, ok := dataMap["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(float64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoAiRechargeByBalanceInUSD"].(int64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoAiRechargeByBalanceInUSD"] == nil {
        data.AutoAiRechargeByBalanceInUsd = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"].(int64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["autoRechargeAiWhenCurrentBalanceFallsInUSD"] == nil {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberNull()
    }
    if val, ok := dataMap["enableAi"].(bool); ok {
        data.EnableAi = types.BoolValue(val)
    }
    if val, ok := dataMap["enableAutoRechargeAiBalance"].(bool); ok {
        data.EnableAutoRechargeAiBalance = types.BoolValue(val)
    }
    if val, ok := dataMap["doNotAddGlobalProbesByDefaultOnNewMonitors"].(bool); ok {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = types.BoolValue(val)
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
    } else if val, ok := dataMap["gitHubAppInstallationId"].(string); ok && val != "" {
        data.GitHubAppInstallationId = types.StringValue(val)
    } else {
        data.GitHubAppInstallationId = types.StringNull()
    }
    if val, ok := dataMap["defaultMetricCardinalityBudget"].(float64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultMetricCardinalityBudget"].(int64); ok {
        data.DefaultMetricCardinalityBudget = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["defaultMetricCardinalityBudget"] == nil {
        data.DefaultMetricCardinalityBudget = types.NumberNull()
    }
    if val, ok := dataMap["defaultTelemetryRetentionInDays"].(float64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["defaultTelemetryRetentionInDays"].(int64); ok {
        data.DefaultTelemetryRetentionInDays = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["defaultTelemetryRetentionInDays"] == nil {
        data.DefaultTelemetryRetentionInDays = types.NumberNull()
    }
    if obj, ok := dataMap["telemetryRetentionConfig"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
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
    } else if val, ok := dataMap["telemetryRetentionConfig"].(string); ok && val != "" {
        data.TelemetryRetentionConfig = NewJSONSubsetValue(val)
    } else {
        data.TelemetryRetentionConfig = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
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
    } else if val, ok := dataMap["defaultMetricDownsamplingRetentionDays"].(string); ok && val != "" {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetValue(val)
    } else {
        data.DefaultMetricDownsamplingRetentionDays = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.CreatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.CreatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = NewJSONSubsetValue(val)
    } else {
        data.CreatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.UpdatedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.UpdatedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = NewJSONSubsetValue(val)
    } else {
        data.UpdatedAt = NewJSONSubsetNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.DeletedAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.DeletedAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = NewJSONSubsetValue(val)
    } else {
        data.DeletedAt = NewJSONSubsetNull()
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["version"].(int); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["version"].(int64); ok {
        data.Version = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
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
    } else if val, ok := dataMap["slug"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderSubscriptionId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionId"].(string); ok && val != "" {
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
    } else if dataMap["paymentProviderSubscriptionSeats"] == nil {
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if obj, ok := dataMap["trialEndsAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.TrialEndsAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.TrialEndsAt = NewJSONSubsetValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.TrialEndsAt = NewJSONSubsetValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.TrialEndsAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrialEndsAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.TrialEndsAt = NewJSONSubsetValue(string(jsonBytes))
            } else {
                data.TrialEndsAt = NewJSONSubsetValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.TrialEndsAt = NewJSONSubsetValue(string(jsonBytes))
        } else {
            data.TrialEndsAt = NewJSONSubsetNull()
        }
    } else if val, ok := dataMap["trialEndsAt"].(string); ok && val != "" {
        data.TrialEndsAt = NewJSONSubsetValue(val)
    } else {
        data.TrialEndsAt = NewJSONSubsetNull()
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
    } else if val, ok := dataMap["paymentProviderCustomerId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderSubscriptionStatus"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if obj, ok := dataMap["deletedByUserId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedByUserId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedByUserId = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedByUserId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedByUserId = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedByUserId = types.StringNull()
        }
    } else if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["workflowRunsInLast30Days"].(int64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["workflowRunsInLast30Days"] == nil {
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(int64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["smsOrCallCurrentBalanceInUSDCents"] == nil {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberNull()
    }
    if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(float64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if val, ok := dataMap["aiCurrentBalanceInUSDCents"].(int64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(float64(val)))
    } else if dataMap["aiCurrentBalanceInUSDCents"] == nil {
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
    } else if val, ok := dataMap["planName"].(string); ok && val != "" {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
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
    } else if val, ok := dataMap["resellerId"].(string); ok && val != "" {
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
    } else if val, ok := dataMap["resellerPlanId"].(string); ok && val != "" {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := dataMap["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

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
    _, err := r.client.Delete("/project/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete project, got error: %s", err))
        return
    }
}


func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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

// Helper method to check if a type string is a valid OneUptime ObjectType
// Only these types should be marshalled/unmarshalled as typed wrapper objects
// This list is dynamically generated from Common/Types/JSON.ts ObjectType enum
func (r *ProjectResource) isValidOneUptimeObjectType(typeStr string) bool {
    validTypes := map[string]bool{
        "ObjectID": true,
        "Decimal": true,
        "Name": true,
        "EqualTo": true,
        "EqualToOrNull": true,
        "MonitorSteps": true,
        "MonitorStep": true,
        "Recurring": true,
        "RestrictionTimes": true,
        "MonitorCriteria": true,
        "PositiveNumber": true,
        "MonitorCriteriaInstance": true,
        "NotEqual": true,
        "Email": true,
        "Phone": true,
        "Color": true,
        "Domain": true,
        "Version": true,
        "IP": true,
        "Route": true,
        "URL": true,
        "Permission": true,
        "Search": true,
        "MultiSearch": true,
        "GreaterThan": true,
        "GreaterThanOrEqual": true,
        "GreaterThanOrNull": true,
        "LessThanOrNull": true,
        "LessThan": true,
        "LessThanOrEqual": true,
        "Port": true,
        "Hostname": true,
        "HashedString": true,
        "DateTime": true,
        "Buffer": true,
        "InBetween": true,
        "NotNull": true,
        "IsNull": true,
        "Includes": true,
        "IncludesAll": true,
        "IncludesNone": true,
        "StartsWith": true,
        "EndsWith": true,
        "NotContains": true,
        "DashboardComponent": true,
        "DashboardViewConfig": true,
    }
    return validTypes[typeStr]
}
