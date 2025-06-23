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
var _ resource.Resource = &StatusPageSubscriberResource{}
var _ resource.ResourceWithImportState = &StatusPageSubscriberResource{}

func NewStatusPageSubscriberResource() resource.Resource {
    return &StatusPageSubscriberResource{}
}

// StatusPageSubscriberResource defines the resource implementation.
type StatusPageSubscriberResource struct {
    client *Client
}

// StatusPageSubscriberResourceModel describes the resource data model.
type StatusPageSubscriberResourceModel struct {
    Id types.String `tfsdk:"id"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    SubscriberEmail types.Map `tfsdk:"subscriber_email"`
    SubscriberPhone types.Map `tfsdk:"subscriber_phone"`
    SubscriberWebhook types.String `tfsdk:"subscriber_webhook"`
    SlackIncomingWebhookUrl types.String `tfsdk:"slack_incoming_webhook_url"`
    SlackWorkspaceName types.String `tfsdk:"slack_workspace_name"`
    IsSubscriptionConfirmed types.Bool `tfsdk:"is_subscription_confirmed"`
    SubscriptionConfirmationToken types.String `tfsdk:"subscription_confirmation_token"`
    IsUnsubscribed types.Bool `tfsdk:"is_unsubscribed"`
    SendYouHaveSubscribedMessage types.Bool `tfsdk:"send_you_have_subscribed_message"`
    IsSubscribedToAllResources types.Bool `tfsdk:"is_subscribed_to_all_resources"`
    IsSubscribedToAllEventTypes types.Bool `tfsdk:"is_subscribed_to_all_event_types"`
    StatusPageResources types.List `tfsdk:"status_page_resources"`
    StatusPageEventTypes types.Map `tfsdk:"status_page_event_types"`
    InternalNote types.String `tfsdk:"internal_note"`
    CreatedAt types.Map `tfsdk:"created_at"`
    UpdatedAt types.Map `tfsdk:"updated_at"`
    DeletedAt types.Map `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
}

func (r *StatusPageSubscriberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_subscriber"
}

func (r *StatusPageSubscriberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_subscriber resource",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique identifier for the resource",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Optional: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "subscriber_email": schema.MapAttribute{
                MarkdownDescription: "Email object",
                Optional: true,
                ElementType: types.StringType,
            },
            "subscriber_phone": schema.MapAttribute{
                MarkdownDescription: "Phone object",
                Optional: true,
                ElementType: types.StringType,
            },
            "subscriber_webhook": schema.StringAttribute{
                MarkdownDescription: "Webhook",
                Optional: true,
            },
            "slack_incoming_webhook_url": schema.StringAttribute{
                MarkdownDescription: "Slack Incoming Webhook URL",
                Optional: true,
            },
            "slack_workspace_name": schema.StringAttribute{
                MarkdownDescription: "Slack Workspace Name",
                Optional: true,
            },
            "is_subscription_confirmed": schema.BoolAttribute{
                MarkdownDescription: "Is Subscription Confirmed",
                Optional: true,
            },
            "subscription_confirmation_token": schema.StringAttribute{
                MarkdownDescription: "Subscription Confirmation Token",
                Optional: true,
            },
            "is_unsubscribed": schema.BoolAttribute{
                MarkdownDescription: "Is Unsubscribed",
                Optional: true,
            },
            "send_you_have_subscribed_message": schema.BoolAttribute{
                MarkdownDescription: "Send You Have Subscribed Message",
                Optional: true,
            },
            "is_subscribed_to_all_resources": schema.BoolAttribute{
                MarkdownDescription: "Is Subscribed to All Resources",
                Optional: true,
            },
            "is_subscribed_to_all_event_types": schema.BoolAttribute{
                MarkdownDescription: "Is Subscribed to All Event Types",
                Optional: true,
            },
            "status_page_resources": schema.ListAttribute{
                MarkdownDescription: "Subscribed to Resources",
                Optional: true,
                ElementType: types.StringType,
            },
            "status_page_event_types": schema.MapAttribute{
                MarkdownDescription: "Subscribed to Event Types",
                Optional: true,
                ElementType: types.StringType,
            },
            "internal_note": schema.StringAttribute{
                MarkdownDescription: "Notes",
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
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
        },
    }
}

func (r *StatusPageSubscriberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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


func (r *StatusPageSubscriberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data StatusPageSubscriberResourceModel

    // Read Terraform plan data into the model
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create API request body
    statusPageSubscriberRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "projectId": data.ProjectId.ValueString(),
        "statusPageId": data.StatusPageId.ValueString(),
        "subscriberEmail": r.convertTerraformMapToInterface(data.SubscriberEmail),
        "subscriberPhone": r.convertTerraformMapToInterface(data.SubscriberPhone),
        "subscriberWebhook": data.SubscriberWebhook.ValueString(),
        "slackIncomingWebhookUrl": data.SlackIncomingWebhookUrl.ValueString(),
        "slackWorkspaceName": data.SlackWorkspaceName.ValueString(),
        "isSubscriptionConfirmed": data.IsSubscriptionConfirmed.ValueBool(),
        "subscriptionConfirmationToken": data.SubscriptionConfirmationToken.ValueString(),
        "isUnsubscribed": data.IsUnsubscribed.ValueBool(),
        "sendYouHaveSubscribedMessage": data.SendYouHaveSubscribedMessage.ValueBool(),
        "isSubscribedToAllResources": data.IsSubscribedToAllResources.ValueBool(),
        "isSubscribedToAllEventTypes": data.IsSubscribedToAllEventTypes.ValueBool(),
        "statusPageResources": r.convertTerraformListToInterface(data.StatusPageResources),
        "statusPageEventTypes": r.convertTerraformMapToInterface(data.StatusPageEventTypes),
        "internalNote": data.InternalNote.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Post("/status-page-subscriber", statusPageSubscriberRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create status_page_subscriber, got error: %s", err))
        return
    }

    var statusPageSubscriberResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageSubscriberResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_subscriber response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageSubscriberResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageSubscriberResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if val, ok := dataMap["subscriberEmail"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SubscriberEmail = mapValue
    } else if dataMap["subscriberEmail"] == nil {
        data.SubscriberEmail = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["subscriberPhone"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SubscriberPhone = mapValue
    } else if dataMap["subscriberPhone"] == nil {
        data.SubscriberPhone = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["subscriberWebhook"].(string); ok && val != "" {
        data.SubscriberWebhook = types.StringValue(val)
    } else {
        data.SubscriberWebhook = types.StringNull()
    }
    if val, ok := dataMap["slackIncomingWebhookUrl"].(string); ok && val != "" {
        data.SlackIncomingWebhookUrl = types.StringValue(val)
    } else {
        data.SlackIncomingWebhookUrl = types.StringNull()
    }
    if val, ok := dataMap["slackWorkspaceName"].(string); ok && val != "" {
        data.SlackWorkspaceName = types.StringValue(val)
    } else {
        data.SlackWorkspaceName = types.StringNull()
    }
    if val, ok := dataMap["isSubscriptionConfirmed"].(bool); ok {
        data.IsSubscriptionConfirmed = types.BoolValue(val)
    } else if dataMap["isSubscriptionConfirmed"] == nil {
        data.IsSubscriptionConfirmed = types.BoolNull()
    }
    if val, ok := dataMap["subscriptionConfirmationToken"].(string); ok && val != "" {
        data.SubscriptionConfirmationToken = types.StringValue(val)
    } else {
        data.SubscriptionConfirmationToken = types.StringNull()
    }
    if val, ok := dataMap["isUnsubscribed"].(bool); ok {
        data.IsUnsubscribed = types.BoolValue(val)
    } else if dataMap["isUnsubscribed"] == nil {
        data.IsUnsubscribed = types.BoolNull()
    }
    if val, ok := dataMap["sendYouHaveSubscribedMessage"].(bool); ok {
        data.SendYouHaveSubscribedMessage = types.BoolValue(val)
    } else if dataMap["sendYouHaveSubscribedMessage"] == nil {
        data.SendYouHaveSubscribedMessage = types.BoolNull()
    }
    if val, ok := dataMap["isSubscribedToAllResources"].(bool); ok {
        data.IsSubscribedToAllResources = types.BoolValue(val)
    } else if dataMap["isSubscribedToAllResources"] == nil {
        data.IsSubscribedToAllResources = types.BoolNull()
    }
    if val, ok := dataMap["isSubscribedToAllEventTypes"].(bool); ok {
        data.IsSubscribedToAllEventTypes = types.BoolValue(val)
    } else if dataMap["isSubscribedToAllEventTypes"] == nil {
        data.IsSubscribedToAllEventTypes = types.BoolNull()
    }
    if val, ok := dataMap["statusPageResources"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.StatusPageResources = listValue
    } else if dataMap["statusPageResources"] == nil {
        data.StatusPageResources = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["statusPageEventTypes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StatusPageEventTypes = mapValue
    } else if dataMap["statusPageEventTypes"] == nil {
        data.StatusPageEventTypes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["internalNote"].(string); ok && val != "" {
        data.InternalNote = types.StringValue(val)
    } else {
        data.InternalNote = types.StringNull()
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
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
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

func (r *StatusPageSubscriberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data StatusPageSubscriberResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Create select parameter to get full object
    selectParam := map[string]interface{}{
        "projectId": true,
        "statusPageId": true,
        "subscriberEmail": true,
        "subscriberPhone": true,
        "subscriberWebhook": true,
        "slackIncomingWebhookUrl": true,
        "slackWorkspaceName": true,
        "isSubscriptionConfirmed": true,
        "subscriptionConfirmationToken": true,
        "isUnsubscribed": true,
        "sendYouHaveSubscribedMessage": true,
        "isSubscribedToAllResources": true,
        "isSubscribedToAllEventTypes": true,
        "statusPageResources": true,
        "statusPageEventTypes": true,
        "internalNote": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    // Make API call with select parameter
    httpResp, err := r.client.PostWithSelect("/status-page-subscriber/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_subscriber, got error: %s", err))
        return
    }

    if httpResp.StatusCode == http.StatusNotFound {
        resp.State.RemoveResource(ctx)
        return
    }

    var statusPageSubscriberResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageSubscriberResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_subscriber response, got error: %s", err))
        return
    }

    // Update the model with response data
    // Extract data from response wrapper
    var dataMap map[string]interface{}
    if wrapper, ok := statusPageSubscriberResponse["data"].(map[string]interface{}); ok {
        // Response is wrapped in a data field
        dataMap = wrapper
    } else {
        // Response is the direct object
        dataMap = statusPageSubscriberResponse
    }

    if val, ok := dataMap["id"].(string); ok && val != "" {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if val, ok := dataMap["subscriberEmail"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SubscriberEmail = mapValue
    } else if dataMap["subscriberEmail"] == nil {
        data.SubscriberEmail = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["subscriberPhone"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SubscriberPhone = mapValue
    } else if dataMap["subscriberPhone"] == nil {
        data.SubscriberPhone = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["subscriberWebhook"].(string); ok && val != "" {
        data.SubscriberWebhook = types.StringValue(val)
    } else {
        data.SubscriberWebhook = types.StringNull()
    }
    if val, ok := dataMap["slackIncomingWebhookUrl"].(string); ok && val != "" {
        data.SlackIncomingWebhookUrl = types.StringValue(val)
    } else {
        data.SlackIncomingWebhookUrl = types.StringNull()
    }
    if val, ok := dataMap["slackWorkspaceName"].(string); ok && val != "" {
        data.SlackWorkspaceName = types.StringValue(val)
    } else {
        data.SlackWorkspaceName = types.StringNull()
    }
    if val, ok := dataMap["isSubscriptionConfirmed"].(bool); ok {
        data.IsSubscriptionConfirmed = types.BoolValue(val)
    } else if dataMap["isSubscriptionConfirmed"] == nil {
        data.IsSubscriptionConfirmed = types.BoolNull()
    }
    if val, ok := dataMap["subscriptionConfirmationToken"].(string); ok && val != "" {
        data.SubscriptionConfirmationToken = types.StringValue(val)
    } else {
        data.SubscriptionConfirmationToken = types.StringNull()
    }
    if val, ok := dataMap["isUnsubscribed"].(bool); ok {
        data.IsUnsubscribed = types.BoolValue(val)
    } else if dataMap["isUnsubscribed"] == nil {
        data.IsUnsubscribed = types.BoolNull()
    }
    if val, ok := dataMap["sendYouHaveSubscribedMessage"].(bool); ok {
        data.SendYouHaveSubscribedMessage = types.BoolValue(val)
    } else if dataMap["sendYouHaveSubscribedMessage"] == nil {
        data.SendYouHaveSubscribedMessage = types.BoolNull()
    }
    if val, ok := dataMap["isSubscribedToAllResources"].(bool); ok {
        data.IsSubscribedToAllResources = types.BoolValue(val)
    } else if dataMap["isSubscribedToAllResources"] == nil {
        data.IsSubscribedToAllResources = types.BoolNull()
    }
    if val, ok := dataMap["isSubscribedToAllEventTypes"].(bool); ok {
        data.IsSubscribedToAllEventTypes = types.BoolValue(val)
    } else if dataMap["isSubscribedToAllEventTypes"] == nil {
        data.IsSubscribedToAllEventTypes = types.BoolNull()
    }
    if val, ok := dataMap["statusPageResources"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.StatusPageResources = listValue
    } else if dataMap["statusPageResources"] == nil {
        data.StatusPageResources = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["statusPageEventTypes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StatusPageEventTypes = mapValue
    } else if dataMap["statusPageEventTypes"] == nil {
        data.StatusPageEventTypes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["internalNote"].(string); ok && val != "" {
        data.InternalNote = types.StringValue(val)
    } else {
        data.InternalNote = types.StringNull()
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
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageSubscriberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var data StatusPageSubscriberResourceModel
    var state StatusPageSubscriberResourceModel

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
    statusPageSubscriberRequest := map[string]interface{}{
        "data": map[string]interface{}{
        "statusPageId": data.StatusPageId.ValueString(),
        "subscriberEmail": r.convertTerraformMapToInterface(data.SubscriberEmail),
        "subscriberPhone": r.convertTerraformMapToInterface(data.SubscriberPhone),
        "subscriberWebhook": data.SubscriberWebhook.ValueString(),
        "slackIncomingWebhookUrl": data.SlackIncomingWebhookUrl.ValueString(),
        "slackWorkspaceName": data.SlackWorkspaceName.ValueString(),
        "isSubscriptionConfirmed": data.IsSubscriptionConfirmed.ValueBool(),
        "subscriptionConfirmationToken": data.SubscriptionConfirmationToken.ValueString(),
        "isUnsubscribed": data.IsUnsubscribed.ValueBool(),
        "sendYouHaveSubscribedMessage": data.SendYouHaveSubscribedMessage.ValueBool(),
        "isSubscribedToAllResources": data.IsSubscribedToAllResources.ValueBool(),
        "isSubscribedToAllEventTypes": data.IsSubscribedToAllEventTypes.ValueBool(),
        "statusPageResources": r.convertTerraformListToInterface(data.StatusPageResources),
        "statusPageEventTypes": r.convertTerraformMapToInterface(data.StatusPageEventTypes),
        "internalNote": data.InternalNote.ValueString(),
        },
    }

    // Make API call
    httpResp, err := r.client.Put("/status-page-subscriber/" + data.Id.ValueString() + "", statusPageSubscriberRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update status_page_subscriber, got error: %s", err))
        return
    }

    // Parse the update response
    var statusPageSubscriberResponse map[string]interface{}
    err = r.client.ParseResponse(httpResp, &statusPageSubscriberResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_subscriber response, got error: %s", err))
        return
    }

    // After successful update, fetch the current state by calling Read with select parameter
    selectParam := map[string]interface{}{
        "projectId": true,
        "statusPageId": true,
        "subscriberEmail": true,
        "subscriberPhone": true,
        "subscriberWebhook": true,
        "slackIncomingWebhookUrl": true,
        "slackWorkspaceName": true,
        "isSubscriptionConfirmed": true,
        "subscriptionConfirmationToken": true,
        "isUnsubscribed": true,
        "sendYouHaveSubscribedMessage": true,
        "isSubscribedToAllResources": true,
        "isSubscribedToAllEventTypes": true,
        "statusPageResources": true,
        "statusPageEventTypes": true,
        "internalNote": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
        "createdByUserId": true,
        "_id": true,
    }

    readResp, err := r.client.PostWithSelect("/status-page-subscriber/" + data.Id.ValueString() + "/get-item", selectParam)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_subscriber after update, got error: %s", err))
        return
    }

    var readResponse map[string]interface{}
    err = r.client.ParseResponse(readResp, &readResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_subscriber read response, got error: %s", err))
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
    if obj, ok := dataMap["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := dataMap["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if val, ok := dataMap["subscriberEmail"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SubscriberEmail = mapValue
    } else if dataMap["subscriberEmail"] == nil {
        data.SubscriberEmail = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["subscriberPhone"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.SubscriberPhone = mapValue
    } else if dataMap["subscriberPhone"] == nil {
        data.SubscriberPhone = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["subscriberWebhook"].(string); ok && val != "" {
        data.SubscriberWebhook = types.StringValue(val)
    } else {
        data.SubscriberWebhook = types.StringNull()
    }
    if val, ok := dataMap["slackIncomingWebhookUrl"].(string); ok && val != "" {
        data.SlackIncomingWebhookUrl = types.StringValue(val)
    } else {
        data.SlackIncomingWebhookUrl = types.StringNull()
    }
    if val, ok := dataMap["slackWorkspaceName"].(string); ok && val != "" {
        data.SlackWorkspaceName = types.StringValue(val)
    } else {
        data.SlackWorkspaceName = types.StringNull()
    }
    if val, ok := dataMap["isSubscriptionConfirmed"].(bool); ok {
        data.IsSubscriptionConfirmed = types.BoolValue(val)
    } else if dataMap["isSubscriptionConfirmed"] == nil {
        data.IsSubscriptionConfirmed = types.BoolNull()
    }
    if val, ok := dataMap["subscriptionConfirmationToken"].(string); ok && val != "" {
        data.SubscriptionConfirmationToken = types.StringValue(val)
    } else {
        data.SubscriptionConfirmationToken = types.StringNull()
    }
    if val, ok := dataMap["isUnsubscribed"].(bool); ok {
        data.IsUnsubscribed = types.BoolValue(val)
    } else if dataMap["isUnsubscribed"] == nil {
        data.IsUnsubscribed = types.BoolNull()
    }
    if val, ok := dataMap["sendYouHaveSubscribedMessage"].(bool); ok {
        data.SendYouHaveSubscribedMessage = types.BoolValue(val)
    } else if dataMap["sendYouHaveSubscribedMessage"] == nil {
        data.SendYouHaveSubscribedMessage = types.BoolNull()
    }
    if val, ok := dataMap["isSubscribedToAllResources"].(bool); ok {
        data.IsSubscribedToAllResources = types.BoolValue(val)
    } else if dataMap["isSubscribedToAllResources"] == nil {
        data.IsSubscribedToAllResources = types.BoolNull()
    }
    if val, ok := dataMap["isSubscribedToAllEventTypes"].(bool); ok {
        data.IsSubscribedToAllEventTypes = types.BoolValue(val)
    } else if dataMap["isSubscribedToAllEventTypes"] == nil {
        data.IsSubscribedToAllEventTypes = types.BoolNull()
    }
    if val, ok := dataMap["statusPageResources"].([]interface{}); ok {
        // Convert API response list to Terraform list
        listValue, _ := types.ListValueFrom(ctx, types.StringType, val)
        data.StatusPageResources = listValue
    } else if dataMap["statusPageResources"] == nil {
        data.StatusPageResources = types.ListNull(types.StringType)
    }
    if val, ok := dataMap["statusPageEventTypes"].(map[string]interface{}); ok {
        // Convert API response map to Terraform map
        mapValue, _ := types.MapValueFrom(ctx, types.StringType, val)
        data.StatusPageEventTypes = mapValue
    } else if dataMap["statusPageEventTypes"] == nil {
        data.StatusPageEventTypes = types.MapNull(types.StringType)
    }
    if val, ok := dataMap["internalNote"].(string); ok && val != "" {
        data.InternalNote = types.StringValue(val)
    } else {
        data.InternalNote = types.StringNull()
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
    if val, ok := dataMap["createdByUserId"].(string); ok && val != "" {
        data.CreatedByUserId = types.StringValue(val)
    } else {
        data.CreatedByUserId = types.StringNull()
    }
    if val, ok := dataMap["_id"].(string); ok {
        data.Id = types.StringValue(val)
    } else {
        data.Id = types.StringNull()
    }

    // Save updated data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StatusPageSubscriberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var data StatusPageSubscriberResourceModel

    // Read Terraform prior state data into the model
    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Make API call
    _, err := r.client.Delete("/status-page-subscriber/" + data.Id.ValueString() + "")
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete status_page_subscriber, got error: %s", err))
        return
    }
}


func (r *StatusPageSubscriberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper method to convert Terraform map to Go interface{}
func (r *StatusPageSubscriberResource) convertTerraformMapToInterface(terraformMap types.Map) interface{} {
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
func (r *StatusPageSubscriberResource) convertTerraformListToInterface(terraformList types.List) interface{} {
    if terraformList.IsNull() || terraformList.IsUnknown() {
        return nil
    }
    
    var result []interface{}
    terraformList.ElementsAs(context.Background(), &result, false)
    return result
}
