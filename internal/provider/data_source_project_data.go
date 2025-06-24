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
    PaymentProviderSubscriptionStatus types.String `tfsdk:"payment_provider_subscription_status"`
    PaymentProviderMeteredSubscriptionStatus types.String `tfsdk:"payment_provider_metered_subscription_status"`
    PaymentProviderPromoCode types.String `tfsdk:"payment_provider_promo_code"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    IsFeatureFlagMonitorGroupsEnabled types.Bool `tfsdk:"is_feature_flag_monitor_groups_enabled"`
    WorkflowRunsInLast30Days types.Number `tfsdk:"workflow_runs_in_last30_days"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    SmsOrCallCurrentBalanceInUsdCents types.Number `tfsdk:"sms_or_call_current_balance_in_usd_cents"`
    AutoRechargeSmsOrCallByBalanceInUsd types.Number `tfsdk:"auto_recharge_sms_or_call_by_balance_in_usd"`
    AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUsd types.Number `tfsdk:"auto_recharge_sms_or_call_when_current_balance_falls_in_usd"`
    EnableSmsNotifications types.Bool `tfsdk:"enable_sms_notifications"`
    EnableCallNotifications types.Bool `tfsdk:"enable_call_notifications"`
    EnableAutoRechargeSmsOrCallBalance types.Bool `tfsdk:"enable_auto_recharge_sms_or_call_balance"`
    PlanName types.String `tfsdk:"plan_name"`
    ResellerId types.String `tfsdk:"reseller_id"`
    ResellerPlanId types.String `tfsdk:"reseller_plan_id"`
    LetCustomerSupportAccessProject types.Bool `tfsdk:"let_customer_support_access_project"`
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
                MarkdownDescription: "Version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_plan_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [Project Owner]",
                Computed: true,
            },
            "payment_provider_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_metered_subscription_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_subscription_seats": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "trial_ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "payment_provider_customer_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_metered_subscription_status": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "payment_provider_promo_code": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [Project Owner, Manage Billing, Edit Project]",
                Computed: true,
            },
            "workflow_runs_in_last30_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read Workflow], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [Project Owner, Project Admin, Edit Project]",
                Computed: true,
            },
            "sms_or_call_current_balance_in_usd_cents": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "auto_recharge_sms_or_call_by_balance_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "auto_recharge_sms_or_call_when_current_balance_falls_in_usd": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "enable_sms_notifications": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "enable_call_notifications": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "enable_auto_recharge_sms_or_call_balance": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [Project Owner, Manage Billing]",
                Computed: true,
            },
            "plan_name": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
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
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project], Update: [Project Owner, Project Admin]",
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
    if val, ok := projectDataResponse["enable_call_notifications"].(bool); ok {
        data.EnableCallNotifications = types.BoolValue(val)
    }
    if val, ok := projectDataResponse["enable_auto_recharge_sms_or_call_balance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
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

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
