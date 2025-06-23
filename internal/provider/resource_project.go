package provider

import (
    "context"
    "fmt"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
    "math/big"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
    PaymentProviderPromoCode types.String `tfsdk:"payment_provider_promo_code"`
    IsFeatureFlagMonitorGroupsEnabled types.Bool `tfsdk:"is_feature_flag_monitor_groups_enabled"`
    ActiveMonitorsLimit types.Number `tfsdk:"active_monitors_limit"`
    SeatLimit types.Number `tfsdk:"seat_limit"`
    UtmContent types.String `tfsdk:"utm_content"`
    RequireSsoForLogin types.Bool `tfsdk:"require_sso_for_login"`
    AutoRechargeSmsOrCallByBalanceInUSD types.Number `tfsdk:"auto_recharge_sms_or_call_by_balance_in_u_s_d"`
    AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD types.Number `tfsdk:"auto_recharge_sms_or_call_when_current_balance_falls_in_u_s_d"`
    EnableSmsNotifications types.Bool `tfsdk:"enable_sms_notifications"`
    EnableCallNotifications types.Bool `tfsdk:"enable_call_notifications"`
    EnableAutoRechargeSmsOrCallBalance types.Bool `tfsdk:"enable_auto_recharge_sms_or_call_balance"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    Slug types.String `tfsdk:"slug"`
    PaymentProviderSubscriptionId types.String `tfsdk:"payment_provider_subscription_id"`
    PaymentProviderMeteredSubscriptionId types.String `tfsdk:"payment_provider_metered_subscription_id"`
    PaymentProviderSubscriptionSeats types.Number `tfsdk:"payment_provider_subscription_seats"`
    TrialEndsAt types.Map `tfsdk:"trial_ends_at"`
    PaymentProviderCustomerId types.String `tfsdk:"payment_provider_customer_id"`
    PaymentProviderSubscriptionStatus types.String `tfsdk:"payment_provider_subscription_status"`
    PaymentProviderMeteredSubscriptionStatus types.String `tfsdk:"payment_provider_metered_subscription_status"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    WorkflowRunsInLast30Days types.Number `tfsdk:"workflow_runs_in_last30_days"`
    SmsOrCallCurrentBalanceInUSDCents types.Number `tfsdk:"sms_or_call_current_balance_in_u_s_d_cents"`
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
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name",
                Required: true,
            },
            "payment_provider_plan_id": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Logged in User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [Project Owner]",
                Optional: true,
            },
            "payment_provider_promo_code": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
                Optional: true,
            },
            "is_feature_flag_monitor_groups_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is Feature Flag Monitor Groups Enabled",
                Optional: true,
            },
            "active_monitors_limit": schema.NumberAttribute{
                Optional: true,
            },
            "seat_limit": schema.NumberAttribute{
                Optional: true,
            },
            "utm_content": schema.StringAttribute{
                Optional: true,
            },
            "require_sso_for_login": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [Project Owner, Project Admin, Edit Project]",
                Optional: true,
            },
            "auto_recharge_sms_or_call_by_balance_in_u_s_d": schema.NumberAttribute{
                MarkdownDescription: "Auto Recharge Amount",
                Optional: true,
            },
            "auto_recharge_sms_or_call_when_current_balance_falls_in_u_s_d": schema.NumberAttribute{
                MarkdownDescription: "Auto Recharge when current balance falls to",
                Optional: true,
            },
            "enable_sms_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable SMS Notifications",
                Optional: true,
            },
            "enable_call_notifications": schema.BoolAttribute{
                MarkdownDescription: "Enable Call Notifications",
                Optional: true,
            },
            "enable_auto_recharge_sms_or_call_balance": schema.BoolAttribute{
                MarkdownDescription: "Enable auto recharge SMS or Call balance",
                Optional: true,
            },
            "created_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "updated_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "deleted_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
            },
            "version": schema.NumberAttribute{
                MarkdownDescription: "Version",
                Computed: true,
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [User], Read: [Project Owner, Project Admin, Project Member, Read Project, Project User], Update: [No access - you don't have permission for this operation]",
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
            "trial_ends_at": schema.MapAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
                ElementType: types.StringType,
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
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "deleted_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "workflow_runs_in_last30_days": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project, Read Workflow], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "sms_or_call_current_balance_in_u_s_d_cents": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Project], Update: [No access - you don't have permission for this operation]",
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
        "paymentProviderPromoCode": data.PaymentProviderPromoCode.ValueString(),
        "isFeatureFlagMonitorGroupsEnabled": data.IsFeatureFlagMonitorGroupsEnabled.ValueBool(),
        "activeMonitorsLimit": data.ActiveMonitorsLimit.ValueBigFloat(),
        "seatLimit": data.SeatLimit.ValueBigFloat(),
        "utmContent": data.UtmContent.ValueString(),
        "requireSsoForLogin": data.RequireSsoForLogin.ValueBool(),
        "autoRechargeSmsOrCallByBalanceInUSD": data.AutoRechargeSmsOrCallByBalanceInUSD.ValueBigFloat(),
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD.ValueBigFloat(),
        "enableSmsNotifications": data.EnableSmsNotifications.ValueBool(),
        "enableCallNotifications": data.EnableCallNotifications.ValueBool(),
        "enableAutoRechargeSmsOrCallBalance": data.EnableAutoRechargeSmsOrCallBalance.ValueBool(),
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

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderPlanId"].(string); ok && val != "" {
        data.PaymentProviderPlanId = types.StringValue(val)
    } else {
        data.PaymentProviderPlanId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderPromoCode"].(string); ok && val != "" {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if val, ok := dataMap["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    } else if dataMap["isFeatureFlagMonitorGroupsEnabled"] == nil {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolNull()
    }
    if val, ok := dataMap["activeMonitorsLimit"].(float64); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(val))
    } else if dataMap["activeMonitorsLimit"] == nil {
        data.ActiveMonitorsLimit = types.NumberNull()
    }
    if val, ok := dataMap["seatLimit"].(float64); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(val))
    } else if dataMap["seatLimit"] == nil {
        data.SeatLimit = types.NumberNull()
    }
    if val, ok := dataMap["utmContent"].(string); ok && val != "" {
        data.UtmContent = types.StringValue(val)
    } else {
        data.UtmContent = types.StringNull()
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    } else if dataMap["requireSsoForLogin"] == nil {
        data.RequireSsoForLogin = types.BoolNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUSD = types.NumberValue(big.NewFloat(val))
    } else if dataMap["autoRechargeSmsOrCallByBalanceInUSD"] == nil {
        data.AutoRechargeSmsOrCallByBalanceInUSD = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD = types.NumberValue(big.NewFloat(val))
    } else if dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"] == nil {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD = types.NumberNull()
    }
    if val, ok := dataMap["enableSmsNotifications"].(bool); ok {
        data.EnableSmsNotifications = types.BoolValue(val)
    } else if dataMap["enableSmsNotifications"] == nil {
        data.EnableSmsNotifications = types.BoolNull()
    }
    if val, ok := dataMap["enableCallNotifications"].(bool); ok {
        data.EnableCallNotifications = types.BoolValue(val)
    } else if dataMap["enableCallNotifications"] == nil {
        data.EnableCallNotifications = types.BoolNull()
    }
    if val, ok := dataMap["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    } else if dataMap["enableAutoRechargeSmsOrCallBalance"] == nil {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionId"].(string); ok && val != "" {
        data.PaymentProviderSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderMeteredSubscriptionId"].(string); ok && val != "" {
        data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionSeats"].(float64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
    } else if dataMap["paymentProviderSubscriptionSeats"] == nil {
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if val, ok := dataMap["trialEndsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TrialEndsAt = mapValue
    } else if dataMap["trialEndsAt"] == nil {
        data.TrialEndsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["paymentProviderCustomerId"].(string); ok && val != "" {
        data.PaymentProviderCustomerId = types.StringValue(val)
    } else {
        data.PaymentProviderCustomerId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionStatus"].(string); ok && val != "" {
        data.PaymentProviderSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionStatus = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(string); ok && val != "" {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if dataMap["workflowRunsInLast30Days"] == nil {
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUSDCents = types.NumberValue(big.NewFloat(val))
    } else if dataMap["smsOrCallCurrentBalanceInUSDCents"] == nil {
        data.SmsOrCallCurrentBalanceInUSDCents = types.NumberNull()
    }
    if val, ok := dataMap["planName"].(string); ok && val != "" {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
    }
    if val, ok := dataMap["resellerId"].(string); ok && val != "" {
        data.ResellerId = types.StringValue(val)
    } else {
        data.ResellerId = types.StringNull()
    }
    if val, ok := dataMap["resellerPlanId"].(string); ok && val != "" {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := dataMap["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    } else if dataMap["letCustomerSupportAccessProject"] == nil {
        data.LetCustomerSupportAccessProject = types.BoolNull()
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
        "paymentProviderPromoCode": true,
        "isFeatureFlagMonitorGroupsEnabled": true,
        "activeMonitorsLimit": true,
        "seatLimit": true,
        "utmContent": true,
        "requireSsoForLogin": true,
        "autoRechargeSmsOrCallByBalanceInUSD": true,
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": true,
        "enableSmsNotifications": true,
        "enableCallNotifications": true,
        "enableAutoRechargeSmsOrCallBalance": true,
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

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderPlanId"].(string); ok && val != "" {
        data.PaymentProviderPlanId = types.StringValue(val)
    } else {
        data.PaymentProviderPlanId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderPromoCode"].(string); ok && val != "" {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if val, ok := dataMap["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    } else if dataMap["isFeatureFlagMonitorGroupsEnabled"] == nil {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolNull()
    }
    if val, ok := dataMap["activeMonitorsLimit"].(float64); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(val))
    } else if dataMap["activeMonitorsLimit"] == nil {
        data.ActiveMonitorsLimit = types.NumberNull()
    }
    if val, ok := dataMap["seatLimit"].(float64); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(val))
    } else if dataMap["seatLimit"] == nil {
        data.SeatLimit = types.NumberNull()
    }
    if val, ok := dataMap["utmContent"].(string); ok && val != "" {
        data.UtmContent = types.StringValue(val)
    } else {
        data.UtmContent = types.StringNull()
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    } else if dataMap["requireSsoForLogin"] == nil {
        data.RequireSsoForLogin = types.BoolNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUSD = types.NumberValue(big.NewFloat(val))
    } else if dataMap["autoRechargeSmsOrCallByBalanceInUSD"] == nil {
        data.AutoRechargeSmsOrCallByBalanceInUSD = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD = types.NumberValue(big.NewFloat(val))
    } else if dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"] == nil {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD = types.NumberNull()
    }
    if val, ok := dataMap["enableSmsNotifications"].(bool); ok {
        data.EnableSmsNotifications = types.BoolValue(val)
    } else if dataMap["enableSmsNotifications"] == nil {
        data.EnableSmsNotifications = types.BoolNull()
    }
    if val, ok := dataMap["enableCallNotifications"].(bool); ok {
        data.EnableCallNotifications = types.BoolValue(val)
    } else if dataMap["enableCallNotifications"] == nil {
        data.EnableCallNotifications = types.BoolNull()
    }
    if val, ok := dataMap["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    } else if dataMap["enableAutoRechargeSmsOrCallBalance"] == nil {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionId"].(string); ok && val != "" {
        data.PaymentProviderSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderMeteredSubscriptionId"].(string); ok && val != "" {
        data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionSeats"].(float64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
    } else if dataMap["paymentProviderSubscriptionSeats"] == nil {
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if val, ok := dataMap["trialEndsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TrialEndsAt = mapValue
    } else if dataMap["trialEndsAt"] == nil {
        data.TrialEndsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["paymentProviderCustomerId"].(string); ok && val != "" {
        data.PaymentProviderCustomerId = types.StringValue(val)
    } else {
        data.PaymentProviderCustomerId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionStatus"].(string); ok && val != "" {
        data.PaymentProviderSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionStatus = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(string); ok && val != "" {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if dataMap["workflowRunsInLast30Days"] == nil {
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUSDCents = types.NumberValue(big.NewFloat(val))
    } else if dataMap["smsOrCallCurrentBalanceInUSDCents"] == nil {
        data.SmsOrCallCurrentBalanceInUSDCents = types.NumberNull()
    }
    if val, ok := dataMap["planName"].(string); ok && val != "" {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
    }
    if val, ok := dataMap["resellerId"].(string); ok && val != "" {
        data.ResellerId = types.StringValue(val)
    } else {
        data.ResellerId = types.StringNull()
    }
    if val, ok := dataMap["resellerPlanId"].(string); ok && val != "" {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := dataMap["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    } else if dataMap["letCustomerSupportAccessProject"] == nil {
        data.LetCustomerSupportAccessProject = types.BoolNull()
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
        "data": map[string]interface{}{
        "name": data.Name.ValueString(),
        "paymentProviderPlanId": data.PaymentProviderPlanId.ValueString(),
        "paymentProviderPromoCode": data.PaymentProviderPromoCode.ValueString(),
        "isFeatureFlagMonitorGroupsEnabled": data.IsFeatureFlagMonitorGroupsEnabled.ValueBool(),
        "activeMonitorsLimit": data.ActiveMonitorsLimit.ValueBigFloat(),
        "seatLimit": data.SeatLimit.ValueBigFloat(),
        "utmContent": data.UtmContent.ValueString(),
        "requireSsoForLogin": data.RequireSsoForLogin.ValueBool(),
        "autoRechargeSmsOrCallByBalanceInUSD": data.AutoRechargeSmsOrCallByBalanceInUSD.ValueBigFloat(),
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD.ValueBigFloat(),
        "enableSmsNotifications": data.EnableSmsNotifications.ValueBool(),
        "enableCallNotifications": data.EnableCallNotifications.ValueBool(),
        "enableAutoRechargeSmsOrCallBalance": data.EnableAutoRechargeSmsOrCallBalance.ValueBool(),
        },
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
        "paymentProviderPromoCode": true,
        "isFeatureFlagMonitorGroupsEnabled": true,
        "activeMonitorsLimit": true,
        "seatLimit": true,
        "utmContent": true,
        "requireSsoForLogin": true,
        "autoRechargeSmsOrCallByBalanceInUSD": true,
        "autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD": true,
        "enableSmsNotifications": true,
        "enableCallNotifications": true,
        "enableAutoRechargeSmsOrCallBalance": true,
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

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if val, ok := dataMap["name"].(string); ok && val != "" {
        data.Name = types.StringValue(val)
    } else {
        data.Name = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderPlanId"].(string); ok && val != "" {
        data.PaymentProviderPlanId = types.StringValue(val)
    } else {
        data.PaymentProviderPlanId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderPromoCode"].(string); ok && val != "" {
        data.PaymentProviderPromoCode = types.StringValue(val)
    } else {
        data.PaymentProviderPromoCode = types.StringNull()
    }
    if val, ok := dataMap["isFeatureFlagMonitorGroupsEnabled"].(bool); ok {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolValue(val)
    } else if dataMap["isFeatureFlagMonitorGroupsEnabled"] == nil {
        data.IsFeatureFlagMonitorGroupsEnabled = types.BoolNull()
    }
    if val, ok := dataMap["activeMonitorsLimit"].(float64); ok {
        data.ActiveMonitorsLimit = types.NumberValue(big.NewFloat(val))
    } else if dataMap["activeMonitorsLimit"] == nil {
        data.ActiveMonitorsLimit = types.NumberNull()
    }
    if val, ok := dataMap["seatLimit"].(float64); ok {
        data.SeatLimit = types.NumberValue(big.NewFloat(val))
    } else if dataMap["seatLimit"] == nil {
        data.SeatLimit = types.NumberNull()
    }
    if val, ok := dataMap["utmContent"].(string); ok && val != "" {
        data.UtmContent = types.StringValue(val)
    } else {
        data.UtmContent = types.StringNull()
    }
    if val, ok := dataMap["requireSsoForLogin"].(bool); ok {
        data.RequireSsoForLogin = types.BoolValue(val)
    } else if dataMap["requireSsoForLogin"] == nil {
        data.RequireSsoForLogin = types.BoolNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallByBalanceInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallByBalanceInUSD = types.NumberValue(big.NewFloat(val))
    } else if dataMap["autoRechargeSmsOrCallByBalanceInUSD"] == nil {
        data.AutoRechargeSmsOrCallByBalanceInUSD = types.NumberNull()
    }
    if val, ok := dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"].(float64); ok {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD = types.NumberValue(big.NewFloat(val))
    } else if dataMap["autoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD"] == nil {
        data.AutoRechargeSmsOrCallWhenCurrentBalanceFallsInUSD = types.NumberNull()
    }
    if val, ok := dataMap["enableSmsNotifications"].(bool); ok {
        data.EnableSmsNotifications = types.BoolValue(val)
    } else if dataMap["enableSmsNotifications"] == nil {
        data.EnableSmsNotifications = types.BoolNull()
    }
    if val, ok := dataMap["enableCallNotifications"].(bool); ok {
        data.EnableCallNotifications = types.BoolValue(val)
    } else if dataMap["enableCallNotifications"] == nil {
        data.EnableCallNotifications = types.BoolNull()
    }
    if val, ok := dataMap["enableAutoRechargeSmsOrCallBalance"].(bool); ok {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolValue(val)
    } else if dataMap["enableAutoRechargeSmsOrCallBalance"] == nil {
        data.EnableAutoRechargeSmsOrCallBalance = types.BoolNull()
    }
    if val, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.CreatedAt = mapValue
    } else if dataMap["createdAt"] == nil {
        data.CreatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.UpdatedAt = mapValue
    } else if dataMap["updatedAt"] == nil {
        data.UpdatedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.DeletedAt = mapValue
    } else if dataMap["deletedAt"] == nil {
        data.DeletedAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    } else if dataMap["version"] == nil {
        data.Version = types.NumberNull()
    }
    if val, ok := dataMap["slug"].(string); ok && val != "" {
        data.Slug = types.StringValue(val)
    } else {
        data.Slug = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionId"].(string); ok && val != "" {
        data.PaymentProviderSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderMeteredSubscriptionId"].(string); ok && val != "" {
        data.PaymentProviderMeteredSubscriptionId = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionSeats"].(float64); ok {
        data.PaymentProviderSubscriptionSeats = types.NumberValue(big.NewFloat(val))
    } else if dataMap["paymentProviderSubscriptionSeats"] == nil {
        data.PaymentProviderSubscriptionSeats = types.NumberNull()
    }
    if val, ok := dataMap["trialEndsAt"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.TrialEndsAt = mapValue
    } else if dataMap["trialEndsAt"] == nil {
        data.TrialEndsAt = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["paymentProviderCustomerId"].(string); ok && val != "" {
        data.PaymentProviderCustomerId = types.StringValue(val)
    } else {
        data.PaymentProviderCustomerId = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderSubscriptionStatus"].(string); ok && val != "" {
        data.PaymentProviderSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderSubscriptionStatus = types.StringNull()
    }
    if val, ok := dataMap["paymentProviderMeteredSubscriptionStatus"].(string); ok && val != "" {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringValue(val)
    } else {
        data.PaymentProviderMeteredSubscriptionStatus = types.StringNull()
    }
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["deletedByUserId"].(string); ok && val != "" {
        data.DeletedByUserId = types.StringValue(val)
    } else {
        data.DeletedByUserId = types.StringNull()
    }
    if val, ok := dataMap["workflowRunsInLast30Days"].(float64); ok {
        data.WorkflowRunsInLast30Days = types.NumberValue(big.NewFloat(val))
    } else if dataMap["workflowRunsInLast30Days"] == nil {
        data.WorkflowRunsInLast30Days = types.NumberNull()
    }
    if val, ok := dataMap["smsOrCallCurrentBalanceInUSDCents"].(float64); ok {
        data.SmsOrCallCurrentBalanceInUSDCents = types.NumberValue(big.NewFloat(val))
    } else if dataMap["smsOrCallCurrentBalanceInUSDCents"] == nil {
        data.SmsOrCallCurrentBalanceInUSDCents = types.NumberNull()
    }
    if val, ok := dataMap["planName"].(string); ok && val != "" {
        data.PlanName = types.StringValue(val)
    } else {
        data.PlanName = types.StringNull()
    }
    if val, ok := dataMap["resellerId"].(string); ok && val != "" {
        data.ResellerId = types.StringValue(val)
    } else {
        data.ResellerId = types.StringNull()
    }
    if val, ok := dataMap["resellerPlanId"].(string); ok && val != "" {
        data.ResellerPlanId = types.StringValue(val)
    } else {
        data.ResellerPlanId = types.StringNull()
    }
    if val, ok := dataMap["letCustomerSupportAccessProject"].(bool); ok {
        data.LetCustomerSupportAccessProject = types.BoolValue(val)
    } else if dataMap["letCustomerSupportAccessProject"] == nil {
        data.LetCustomerSupportAccessProject = types.BoolNull()
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
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
