package provider

import (
    "context"
    "fmt"
    "math/big"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProjectDataDataSource{}

func NewProjectDataDataSource() datasource.DataSource {
    return &ProjectDataDataSource{}
}

// ProjectDataDataSource defines the data source implementation.
type ProjectDataDataSource struct {
    client *Client
}

// ProjectDataDataSourceModel describes the data source data model.
type ProjectDataDataSourceModel struct {
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
    EnableCallNotifications types.Bool `tfsdk:"enable_call_notifications"`
    EnableAutoRechargeSmsOrCallBalance types.Bool `tfsdk:"enable_auto_recharge_sms_or_call_balance"`
    AiCurrentBalanceInUsdCents types.Number `tfsdk:"ai_current_balance_in_usd_cents"`
    AutoAiRechargeByBalanceInUsd types.Number `tfsdk:"auto_ai_recharge_by_balance_in_usd"`
    AutoRechargeAiWhenCurrentBalanceFallsInUsd types.Number `tfsdk:"auto_recharge_ai_when_current_balance_falls_in_usd"`
    EnableAi types.Bool `tfsdk:"enable_ai"`
    EnableAutoRechargeAiBalance types.Bool `tfsdk:"enable_auto_recharge_ai_balance"`
    SendInvoicesByEmail types.Bool `tfsdk:"send_invoices_by_email"`
    PlanName types.String `tfsdk:"plan_name"`
    ResellerId types.String `tfsdk:"reseller_id"`
    ResellerPlanId types.String `tfsdk:"reseller_plan_id"`
    LetCustomerSupportAccessProject types.Bool `tfsdk:"let_customer_support_access_project"`
    DoNotAddGlobalProbesByDefaultOnNewMonitors types.Bool `tfsdk:"do_not_add_global_probes_by_default_on_new_monitors"`
    GitHubAppInstallationId types.String `tfsdk:"git_hub_app_installation_id"`
}

func (d *ProjectDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_project_data"
}

func (d *ProjectDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "project_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
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
                MarkdownDescription: "Friendly globally unique name for your object. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_plan_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "payment_provider_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_metered_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_subscription_seats": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "trial_ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "payment_provider_customer_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "business_details": schema.StringAttribute{
                MarkdownDescription: "Business legal name, address and any tax information to appear on invoices.. Permissions - Create: [Project Owner, Manage Billing], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "business_details_country": schema.StringAttribute{
                MarkdownDescription: "Two-letter ISO country code for billing address (e.g., US, GB, DE).. Permissions - Create: [Project Owner, Manage Billing], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "finance_accounting_email": schema.StringAttribute{
                MarkdownDescription: "Email object",
                Computed: true,
            },
            "payment_provider_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_metered_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_promo_code": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "Is Feature Flag Monitor Groups Enabled. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing, Edit Project]",
                Computed: true,
            },
            "workflow_runs_in_last30_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read Workflow, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Project]",
                Computed: true,
            },
            "incident_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident numbers (e.g., 'INC-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Project]",
                Computed: true,
            },
            "alert_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert numbers (e.g., 'ALT-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Project]",
                Computed: true,
            },
            "scheduled_maintenance_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for scheduled maintenance numbers (e.g., 'SM-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Project]",
                Computed: true,
            },
            "incident_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for incident episode numbers (e.g., 'IE-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Project]",
                Computed: true,
            },
            "alert_episode_number_prefix": schema.StringAttribute{
                MarkdownDescription: "Custom prefix for alert episode numbers (e.g., 'AE-'). If empty, '#' is used.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Project]",
                Computed: true,
            },
            "sms_or_call_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for SMS, Call, and WhatsApp. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "auto_recharge_sms_or_call_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for SMS, Call, and WhatsApp. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "auto_recharge_sms_or_call_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for SMS, Call, and WhatsApp. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "enable_sms_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable SMS notifications for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "enable_whats_app_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable WhatsApp notifications for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "enable_call_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable call notifications for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "enable_auto_recharge_sms_or_call_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for SMS, Call, and WhatsApp balance for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "ai_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Balance in USD for AI services. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "auto_ai_recharge_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge amount in USD for AI services. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "auto_recharge_ai_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Auto recharge is triggered when current balance falls to this amount in USD for AI services. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "enable_ai": schema.BoolAttribute{
                MarkdownDescription: "Enable AI services for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "enable_auto_recharge_ai_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge for AI balance for this project.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "send_invoices_by_email": schema.BoolAttribute{
                MarkdownDescription: "When enabled, invoices will be automatically sent to the finance/accounting email when they are generated.. Permissions - Create: [Project Owner, Manage Billing], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "plan_name": schema.StringAttribute{
                MarkdownDescription: "Name of the plan this project is subscribed to.. Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "OneUptime customer support can access this project. This is used for debugging purposes.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Project Admin]",
                Computed: true,
            },
            "do_not_add_global_probes_by_default_on_new_monitors": schema.BoolAttribute{
                MarkdownDescription: "If enabled, global probes will NOT be automatically added to new monitors. Enable this only if you are using ONLY custom probes to monitor your resources.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User, Read All Project Resources], Update: [Project Owner, Project Admin, Edit Project]",
                Computed: true,
            },
            "git_hub_app_installation_id": schema.StringAttribute{
                MarkdownDescription: "The GitHub App installation ID for this project. This is set when the GitHub App is installed on the organization.. Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read All Project Resources], Update: [Project Owner, Project Admin]",
                Computed: true,
            },
        },
    }
}

func (d *ProjectDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProjectDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ProjectDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "project" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read project_data, got error: %s", err))
        return
    }

    var projectDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &projectDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse project_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := projectDataResponse["data"].(map[string]interface{}); ok {
        projectDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := projectDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := projectDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := projectDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := projectDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := projectDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := projectDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := projectDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := projectDataResponse["payment_provider_plan_id"].(string); ok {
        data.PaymentProviderPlanId = types.StringValue(val)
    }
    if val, ok := projectDataResponse["payment_provider_subscription_id"].(string); ok {
        data.PaymentProviderSubscriptionId = types.StringValue(val)
    }
    if val, ok := projectDataResponse["payment_provider_metered_subscription_id"].(string); ok {
        data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
    }
    if val, ok := projectDataResponse["payment_provider_subscription_seats"].(float64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := projectDataResponse["trial_ends_at"].(string); ok {
        data.TrialEndsAt = types.StringValue(val)
    }
    if val, ok := projectDataResponse["payment_provider_customer_id"].(string); ok {
        data.PaymentProviderCustomerId = types.StringValue(val)
    }
    if val, ok := projectDataResponse["business_details"].(string); ok {
        data.BusinessDetails = types.StringValue(val)
    }
    if val, ok := projectDataResponse["business_details_country"].(string); ok {
        data.BusinessDetailsCountry = types.StringValue(val)
    }
    if val, ok := projectDataResponse["finance_accounting_email"].(string); ok {
        data.FinanceAccountingEmail = types.StringValue(val)
    }
    if val, ok := projectDataResponse["payment_provider_subscription_status"].(string); ok {
        data.PaymentProviderSubscriptionStatus = types.StringValue(val)
    }
    if val, ok := projectDataResponse["payment_provider_metered_subscription_status"].(string); ok {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
    }
    if val, ok := projectDataResponse["payment_provider_promo_code"].(string); ok {
        data.PaymentProviderPromoCode = types.StringValue(val)
    }
    if val, ok := projectDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := projectDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := projectDataResponse["is_feature_flag_monitor_groups_enabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["workflow_runs_in_last30_days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := projectDataResponse["require_sso_for_login"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["incident_number_prefix"].(string); ok {
        data.IncidentNumberPrefix = types.StringValue(val)
    }
    if val, ok := projectDataResponse["alert_number_prefix"].(string); ok {
        data.AlertNumberPrefix = types.StringValue(val)
    }
    if val, ok := projectDataResponse["scheduled_maintenance_number_prefix"].(string); ok {
        data.ScheduledMaintenanceNumberPrefix = types.StringValue(val)
    }
    if val, ok := projectDataResponse["incident_episode_number_prefix"].(string); ok {
        data.IncidentEpisodeNumberPrefix = types.StringValue(val)
    }
    if val, ok := projectDataResponse["alert_episode_number_prefix"].(string); ok {
        data.AlertEpisodeNumberPrefix = types.StringValue(val)
    }
    if val, ok := projectDataResponse["sms_or_call_current_balance_in_usd_cents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := projectDataResponse["auto_recharge_sms_or_call_by_balance_in_usd"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := projectDataResponse["auto_recharge_sms_or_call_when_current_balance_falls_in_usd"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := projectDataResponse["enable_sms_notifications"].(bool); ok {
        data.EnableSmsNotifications = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["enable_whats_app_notifications"].(bool); ok {
        data.EnableWhatsAppNotifications = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["enable_call_notifications"].(bool); ok {
        data.EnableCallNotifications = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["enable_auto_recharge_sms_or_call_balance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["ai_current_balance_in_usd_cents"].(float64); ok {
        data.AiCurrentBalanceInUsdCents = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := projectDataResponse["auto_ai_recharge_by_balance_in_usd"].(float64); ok {
        data.AutoAiRechargeByBalanceInUsd = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := projectDataResponse["auto_recharge_ai_when_current_balance_falls_in_usd"].(float64); ok {
        data.AutoRechargeAiWhenCurrentBalanceFallsInUsd = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := projectDataResponse["enable_ai"].(bool); ok {
        data.EnableAi = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["enable_auto_recharge_ai_balance"].(bool); ok {
        data.EnableAutoRechargeAiBalance = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["send_invoices_by_email"].(bool); ok {
        data.SendInvoicesByEmail = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["plan_name"].(string); ok {
        data.PlanName = types.StringValue(val)
    }
    if val, ok := projectDataResponse["reseller_id"].(string); ok {
        data.ResellerId = types.StringValue(val)
    }
    if val, ok := projectDataResponse["reseller_plan_id"].(string); ok {
        data.ResellerPlanId = types.StringValue(val)
    }
    if val, ok := projectDataResponse["let_customer_support_access_project"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["do_not_add_global_probes_by_default_on_new_monitors"].(bool); ok {
        data.DoNotAddGlobalProbesByDefaultOnNewMonitors = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["git_hub_app_installation_id"].(string); ok {
        data.GitHubAppInstallationId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
