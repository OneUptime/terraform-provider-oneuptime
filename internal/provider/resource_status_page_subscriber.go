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
    "encoding/json"
    "net/url"
    "strings"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
    SubscriberEmail types.String `tfsdk:"subscriber_email"`
    SubscriberPhone types.String `tfsdk:"subscriber_phone"`
    SubscriberWebhook types.String `tfsdk:"subscriber_webhook"`
    SlackIncomingWebhookUrl types.String `tfsdk:"slack_incoming_webhook_url"`
    SlackWorkspaceName types.String `tfsdk:"slack_workspace_name"`
    MicrosoftTeamsIncomingWebhookUrl types.String `tfsdk:"microsoft_teams_incoming_webhook_url"`
    MicrosoftTeamsWorkspaceName types.String `tfsdk:"microsoft_teams_workspace_name"`
    IsSubscriptionConfirmed types.Bool `tfsdk:"is_subscription_confirmed"`
    SubscriptionConfirmationToken types.String `tfsdk:"subscription_confirmation_token"`
    IsUnsubscribed types.Bool `tfsdk:"is_unsubscribed"`
    SendYouHaveSubscribedMessage types.Bool `tfsdk:"send_you_have_subscribed_message"`
    IsSubscribedToAllResources types.Bool `tfsdk:"is_subscribed_to_all_resources"`
    IsSubscribedToAllEventTypes types.Bool `tfsdk:"is_subscribed_to_all_event_types"`
    StatusPageResources types.Set `tfsdk:"status_page_resources"`
    StatusPageEventTypes types.String `tfsdk:"status_page_event_types"`
    InternalNote types.String `tfsdk:"internal_note"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
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
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Required: true,
            },
            "subscriber_email": schema.StringAttribute{
                MarkdownDescription: "Email object",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "subscriber_phone": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "subscriber_webhook": schema.StringAttribute{
                MarkdownDescription: "Webhook to ping when events happen on Status Page. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "slack_incoming_webhook_url": schema.StringAttribute{
                MarkdownDescription: "Slack incoming webhook URL to send notifications to Slack channel. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [No access - you don't have permission for this operation], Update: [No access - you don't have permission for this operation]",
                Optional: true,
            },
            "slack_workspace_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Slack workspace for validation and identification. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "microsoft_teams_incoming_webhook_url": schema.StringAttribute{
                MarkdownDescription: "Microsoft Teams incoming webhook URL to send notifications to Teams channel. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [No access - you don't have permission for this operation], Update: [No access - you don't have permission for this operation]",
                Optional: true,
            },
            "microsoft_teams_workspace_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Microsoft Teams workspace for validation and identification. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "is_subscription_confirmed": schema.BoolAttribute{
                MarkdownDescription: "Has subscriber confirmed their subscription? (for example, by clicking on a confirmation link in an email). Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "subscription_confirmation_token": schema.StringAttribute{
                MarkdownDescription: "Token used to confirm subscription. This is a random token that is sent to the subscriber's email address to confirm their subscription.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber], Read: [No access - you don't have permission for this operation], Update: [No access - you don't have permission for this operation]",
                Optional: true,
            },
            "is_unsubscribed": schema.BoolAttribute{
                MarkdownDescription: "Is Subscriber Unsubscribed?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(false),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "send_you_have_subscribed_message": schema.BoolAttribute{
                MarkdownDescription: "Send You Have Subscribed Message when subscriber is created?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [No access - you don't have permission for this operation]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "is_subscribed_to_all_resources": schema.BoolAttribute{
                MarkdownDescription: "Is Subscriber Subscribed to All Resources on this status page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "is_subscribed_to_all_event_types": schema.BoolAttribute{
                MarkdownDescription: "Is Subscriber Subscribed to All Event Types (like Incidents, Scheduled Events, Announcements) on this status page?. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                Default: booldefault.StaticBool(true),
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.UseStateForUnknown(),
                },
            },
            "status_page_resources": schema.SetAttribute{
                MarkdownDescription: "Relation to Status Page Resources where this subscriber is subscribed to. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                ElementType: types.StringType,
                PlanModifiers: []planmodifier.Set{
                    setplanmodifier.UseStateForUnknown(),
                },
            },
            "status_page_event_types": schema.StringAttribute{
                MarkdownDescription: "Which event types is the subscriber subscribed to (like Incidents, Scheduled Events, Announcements). Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "internal_note": schema.StringAttribute{
                MarkdownDescription: "Any notes or text you would like to add to this subscriber object. This is for internal use only.. Permissions - Create: [Project Owner, Project Admin, Project Member, Status Page Manager, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Viewer, Status Page Manager, Read Status Page Subscriber, Read All Project Resources], Update: [Project Owner, Project Admin, Project Member, Status Page Manager, Edit Status Page Subscriber]",
                Optional: true,
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
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
        "statusPageId": data.StatusPageId.ValueString(),
        "subscriberEmail": r.parseJSONField(data.SubscriberEmail),
        "subscriberPhone": r.parseJSONField(data.SubscriberPhone),
        "subscriberWebhook": data.SubscriberWebhook.ValueString(),
        "slackIncomingWebhookUrl": data.SlackIncomingWebhookUrl.ValueString(),
        "slackWorkspaceName": data.SlackWorkspaceName.ValueString(),
        "microsoftTeamsIncomingWebhookUrl": data.MicrosoftTeamsIncomingWebhookUrl.ValueString(),
        "microsoftTeamsWorkspaceName": data.MicrosoftTeamsWorkspaceName.ValueString(),
        "isSubscriptionConfirmed": data.IsSubscriptionConfirmed.ValueBool(),
        "subscriptionConfirmationToken": data.SubscriptionConfirmationToken.ValueString(),
        "isUnsubscribed": data.IsUnsubscribed.ValueBool(),
        "sendYouHaveSubscribedMessage": data.SendYouHaveSubscribedMessage.ValueBool(),
        "isSubscribedToAllResources": data.IsSubscribedToAllResources.ValueBool(),
        "isSubscribedToAllEventTypes": data.IsSubscribedToAllEventTypes.ValueBool(),
        "statusPageResources": r.convertTerraformSetToInterface(data.StatusPageResources),
        "statusPageEventTypes": r.parseJSONField(data.StatusPageEventTypes),
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
    if obj, ok := dataMap["statusPageId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberEmail = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberEmail = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberEmail = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberEmail = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberEmail"].(string); ok && val != "" {
        data.SubscriberEmail = types.StringValue(val)
    } else {
        data.SubscriberEmail = types.StringNull()
    }
    if obj, ok := dataMap["subscriberPhone"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberPhone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberPhone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberPhone = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberPhone = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberPhone = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberPhone = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberPhone"].(string); ok && val != "" {
        data.SubscriberPhone = types.StringValue(val)
    } else {
        data.SubscriberPhone = types.StringNull()
    }
    if obj, ok := dataMap["subscriberWebhook"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberWebhook = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberWebhook = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberWebhook = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberWebhook = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberWebhook = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberWebhook = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberWebhook"].(string); ok && val != "" {
        data.SubscriberWebhook = types.StringValue(val)
    } else {
        data.SubscriberWebhook = types.StringNull()
    }
    if obj, ok := dataMap["slackIncomingWebhookUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SlackIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SlackIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SlackIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SlackIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.SlackIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SlackIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.SlackIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SlackIncomingWebhookUrl = types.StringValue(string(jsonBytes))
        } else {
            data.SlackIncomingWebhookUrl = types.StringNull()
        }
    } else if val, ok := dataMap["slackIncomingWebhookUrl"].(string); ok && val != "" {
        data.SlackIncomingWebhookUrl = types.StringValue(val)
    } else {
        data.SlackIncomingWebhookUrl = types.StringNull()
    }
    if obj, ok := dataMap["slackWorkspaceName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SlackWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SlackWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
        } else {
            data.SlackWorkspaceName = types.StringNull()
        }
    } else if val, ok := dataMap["slackWorkspaceName"].(string); ok && val != "" {
        data.SlackWorkspaceName = types.StringValue(val)
    } else {
        data.SlackWorkspaceName = types.StringNull()
    }
    if obj, ok := dataMap["microsoftTeamsIncomingWebhookUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(string(jsonBytes))
        } else {
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringNull()
        }
    } else if val, ok := dataMap["microsoftTeamsIncomingWebhookUrl"].(string); ok && val != "" {
        data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(val)
    } else {
        data.MicrosoftTeamsIncomingWebhookUrl = types.StringNull()
    }
    if obj, ok := dataMap["microsoftTeamsWorkspaceName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
        } else {
            data.MicrosoftTeamsWorkspaceName = types.StringNull()
        }
    } else if val, ok := dataMap["microsoftTeamsWorkspaceName"].(string); ok && val != "" {
        data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
    } else {
        data.MicrosoftTeamsWorkspaceName = types.StringNull()
    }
    if val, ok := dataMap["isSubscriptionConfirmed"].(bool); ok {
        data.IsSubscriptionConfirmed = types.BoolValue(val)
    }
    if obj, ok := dataMap["subscriptionConfirmationToken"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriptionConfirmationToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriptionConfirmationToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriptionConfirmationToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriptionConfirmationToken = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriptionConfirmationToken = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriptionConfirmationToken = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriptionConfirmationToken = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriptionConfirmationToken = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriptionConfirmationToken = types.StringNull()
        }
    } else if val, ok := dataMap["subscriptionConfirmationToken"].(string); ok && val != "" {
        data.SubscriptionConfirmationToken = types.StringValue(val)
    } else {
        data.SubscriptionConfirmationToken = types.StringNull()
    }
    if val, ok := dataMap["isUnsubscribed"].(bool); ok {
        data.IsUnsubscribed = types.BoolValue(val)
    }
    if val, ok := dataMap["sendYouHaveSubscribedMessage"].(bool); ok {
        data.SendYouHaveSubscribedMessage = types.BoolValue(val)
    }
    if val, ok := dataMap["isSubscribedToAllResources"].(bool); ok {
        data.IsSubscribedToAllResources = types.BoolValue(val)
    }
    if val, ok := dataMap["isSubscribedToAllEventTypes"].(bool); ok {
        data.IsSubscribedToAllEventTypes = types.BoolValue(val)
    }
    if val, ok := dataMap["statusPageResources"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.StatusPageResources = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.StatusPageResources = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["statusPageEventTypes"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageEventTypes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageEventTypes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageEventTypes = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageEventTypes"].(string); ok && val != "" {
        data.StatusPageEventTypes = types.StringValue(val)
    } else {
        data.StatusPageEventTypes = types.StringNull()
    }
    if obj, ok := dataMap["internalNote"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InternalNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.InternalNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.InternalNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.InternalNote = types.StringValue(string(jsonBytes))
            } else {
                data.InternalNote = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.InternalNote = types.StringValue(string(jsonBytes))
            } else {
                data.InternalNote = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.InternalNote = types.StringValue(string(jsonBytes))
        } else {
            data.InternalNote = types.StringNull()
        }
    } else if val, ok := dataMap["internalNote"].(string); ok && val != "" {
        data.InternalNote = types.StringValue(val)
    } else {
        data.InternalNote = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
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
        "microsoftTeamsIncomingWebhookUrl": true,
        "microsoftTeamsWorkspaceName": true,
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
    if obj, ok := dataMap["statusPageId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberEmail = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberEmail = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberEmail = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberEmail = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberEmail"].(string); ok && val != "" {
        data.SubscriberEmail = types.StringValue(val)
    } else {
        data.SubscriberEmail = types.StringNull()
    }
    if obj, ok := dataMap["subscriberPhone"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberPhone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberPhone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberPhone = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberPhone = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberPhone = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberPhone = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberPhone"].(string); ok && val != "" {
        data.SubscriberPhone = types.StringValue(val)
    } else {
        data.SubscriberPhone = types.StringNull()
    }
    if obj, ok := dataMap["subscriberWebhook"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberWebhook = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberWebhook = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberWebhook = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberWebhook = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberWebhook = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberWebhook = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberWebhook"].(string); ok && val != "" {
        data.SubscriberWebhook = types.StringValue(val)
    } else {
        data.SubscriberWebhook = types.StringNull()
    }
    if obj, ok := dataMap["slackIncomingWebhookUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SlackIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SlackIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SlackIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SlackIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.SlackIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SlackIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.SlackIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SlackIncomingWebhookUrl = types.StringValue(string(jsonBytes))
        } else {
            data.SlackIncomingWebhookUrl = types.StringNull()
        }
    } else if val, ok := dataMap["slackIncomingWebhookUrl"].(string); ok && val != "" {
        data.SlackIncomingWebhookUrl = types.StringValue(val)
    } else {
        data.SlackIncomingWebhookUrl = types.StringNull()
    }
    if obj, ok := dataMap["slackWorkspaceName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SlackWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SlackWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
        } else {
            data.SlackWorkspaceName = types.StringNull()
        }
    } else if val, ok := dataMap["slackWorkspaceName"].(string); ok && val != "" {
        data.SlackWorkspaceName = types.StringValue(val)
    } else {
        data.SlackWorkspaceName = types.StringNull()
    }
    if obj, ok := dataMap["microsoftTeamsIncomingWebhookUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(string(jsonBytes))
        } else {
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringNull()
        }
    } else if val, ok := dataMap["microsoftTeamsIncomingWebhookUrl"].(string); ok && val != "" {
        data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(val)
    } else {
        data.MicrosoftTeamsIncomingWebhookUrl = types.StringNull()
    }
    if obj, ok := dataMap["microsoftTeamsWorkspaceName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
        } else {
            data.MicrosoftTeamsWorkspaceName = types.StringNull()
        }
    } else if val, ok := dataMap["microsoftTeamsWorkspaceName"].(string); ok && val != "" {
        data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
    } else {
        data.MicrosoftTeamsWorkspaceName = types.StringNull()
    }
    if val, ok := dataMap["isSubscriptionConfirmed"].(bool); ok {
        data.IsSubscriptionConfirmed = types.BoolValue(val)
    }
    if obj, ok := dataMap["subscriptionConfirmationToken"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriptionConfirmationToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriptionConfirmationToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriptionConfirmationToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriptionConfirmationToken = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriptionConfirmationToken = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriptionConfirmationToken = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriptionConfirmationToken = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriptionConfirmationToken = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriptionConfirmationToken = types.StringNull()
        }
    } else if val, ok := dataMap["subscriptionConfirmationToken"].(string); ok && val != "" {
        data.SubscriptionConfirmationToken = types.StringValue(val)
    } else {
        data.SubscriptionConfirmationToken = types.StringNull()
    }
    if val, ok := dataMap["isUnsubscribed"].(bool); ok {
        data.IsUnsubscribed = types.BoolValue(val)
    }
    if val, ok := dataMap["sendYouHaveSubscribedMessage"].(bool); ok {
        data.SendYouHaveSubscribedMessage = types.BoolValue(val)
    }
    if val, ok := dataMap["isSubscribedToAllResources"].(bool); ok {
        data.IsSubscribedToAllResources = types.BoolValue(val)
    }
    if val, ok := dataMap["isSubscribedToAllEventTypes"].(bool); ok {
        data.IsSubscribedToAllEventTypes = types.BoolValue(val)
    }
    if val, ok := dataMap["statusPageResources"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.StatusPageResources = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.StatusPageResources = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["statusPageEventTypes"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageEventTypes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageEventTypes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageEventTypes = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageEventTypes"].(string); ok && val != "" {
        data.StatusPageEventTypes = types.StringValue(val)
    } else {
        data.StatusPageEventTypes = types.StringNull()
    }
    if obj, ok := dataMap["internalNote"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InternalNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.InternalNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.InternalNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.InternalNote = types.StringValue(string(jsonBytes))
            } else {
                data.InternalNote = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.InternalNote = types.StringValue(string(jsonBytes))
            } else {
                data.InternalNote = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.InternalNote = types.StringValue(string(jsonBytes))
        } else {
            data.InternalNote = types.StringNull()
        }
    } else if val, ok := dataMap["internalNote"].(string); ok && val != "" {
        data.InternalNote = types.StringValue(val)
    } else {
        data.InternalNote = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
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
        "data": map[string]interface{}{},
    }
    requestDataMap := statusPageSubscriberRequest["data"].(map[string]interface{})

    if !data.SubscriberEmail.IsUnknown() && !state.SubscriberEmail.IsUnknown() && !data.SubscriberEmail.Equal(state.SubscriberEmail) {
        var subscriberemailData interface{}
        if err := json.Unmarshal([]byte(data.SubscriberEmail.ValueString()), &subscriberemailData); err == nil {
            requestDataMap["subscriberEmail"] = subscriberemailData
        } else {
            requestDataMap["subscriberEmail"] = data.SubscriberEmail.ValueString()
        }
    }
    if !data.SubscriberPhone.IsUnknown() && !state.SubscriberPhone.IsUnknown() && !data.SubscriberPhone.Equal(state.SubscriberPhone) {
        var subscriberphoneData interface{}
        if err := json.Unmarshal([]byte(data.SubscriberPhone.ValueString()), &subscriberphoneData); err == nil {
            requestDataMap["subscriberPhone"] = subscriberphoneData
        } else {
            requestDataMap["subscriberPhone"] = data.SubscriberPhone.ValueString()
        }
    }
    if !data.SubscriberWebhook.IsUnknown() && !state.SubscriberWebhook.IsUnknown() && !data.SubscriberWebhook.Equal(state.SubscriberWebhook) {
        requestDataMap["subscriberWebhook"] = data.SubscriberWebhook.ValueString()
    }
    if !data.SlackWorkspaceName.IsUnknown() && !state.SlackWorkspaceName.IsUnknown() && !data.SlackWorkspaceName.Equal(state.SlackWorkspaceName) {
        requestDataMap["slackWorkspaceName"] = data.SlackWorkspaceName.ValueString()
    }
    if !data.MicrosoftTeamsWorkspaceName.IsUnknown() && !state.MicrosoftTeamsWorkspaceName.IsUnknown() && !data.MicrosoftTeamsWorkspaceName.Equal(state.MicrosoftTeamsWorkspaceName) {
        requestDataMap["microsoftTeamsWorkspaceName"] = data.MicrosoftTeamsWorkspaceName.ValueString()
    }
    if !data.IsSubscriptionConfirmed.IsUnknown() && !state.IsSubscriptionConfirmed.IsUnknown() && !data.IsSubscriptionConfirmed.Equal(state.IsSubscriptionConfirmed) {
        requestDataMap["isSubscriptionConfirmed"] = data.IsSubscriptionConfirmed.ValueBool()
    }
    if !data.IsUnsubscribed.IsUnknown() && !state.IsUnsubscribed.IsUnknown() && !data.IsUnsubscribed.Equal(state.IsUnsubscribed) {
        requestDataMap["isUnsubscribed"] = data.IsUnsubscribed.ValueBool()
    }
    if !data.IsSubscribedToAllResources.IsUnknown() && !state.IsSubscribedToAllResources.IsUnknown() && !data.IsSubscribedToAllResources.Equal(state.IsSubscribedToAllResources) {
        requestDataMap["isSubscribedToAllResources"] = data.IsSubscribedToAllResources.ValueBool()
    }
    if !data.IsSubscribedToAllEventTypes.IsUnknown() && !state.IsSubscribedToAllEventTypes.IsUnknown() && !data.IsSubscribedToAllEventTypes.Equal(state.IsSubscribedToAllEventTypes) {
        requestDataMap["isSubscribedToAllEventTypes"] = data.IsSubscribedToAllEventTypes.ValueBool()
    }
    if !data.StatusPageResources.IsUnknown() && !state.StatusPageResources.IsUnknown() && !data.StatusPageResources.Equal(state.StatusPageResources) {
        requestDataMap["statusPageResources"] = r.convertTerraformSetToInterface(data.StatusPageResources)
    }
    if !data.StatusPageEventTypes.IsUnknown() && !state.StatusPageEventTypes.IsUnknown() && !data.StatusPageEventTypes.Equal(state.StatusPageEventTypes) {
        var statuspageeventtypesData interface{}
        if err := json.Unmarshal([]byte(data.StatusPageEventTypes.ValueString()), &statuspageeventtypesData); err == nil {
            requestDataMap["statusPageEventTypes"] = statuspageeventtypesData
        } else {
            requestDataMap["statusPageEventTypes"] = data.StatusPageEventTypes.ValueString()
        }
    }
    if !data.InternalNote.IsUnknown() && !state.InternalNote.IsUnknown() && !data.InternalNote.Equal(state.InternalNote) {
        requestDataMap["internalNote"] = data.InternalNote.ValueString()
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
        "microsoftTeamsIncomingWebhookUrl": true,
        "microsoftTeamsWorkspaceName": true,
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
    if obj, ok := dataMap["statusPageId"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusPageId = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageId = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageId"].(string); ok && val != "" {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := dataMap["subscriberEmail"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberEmail = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberEmail = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberEmail = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberEmail = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberEmail"].(string); ok && val != "" {
        data.SubscriberEmail = types.StringValue(val)
    } else {
        data.SubscriberEmail = types.StringNull()
    }
    if obj, ok := dataMap["subscriberPhone"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberPhone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberPhone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberPhone = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberPhone = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberPhone = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberPhone = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberPhone"].(string); ok && val != "" {
        data.SubscriberPhone = types.StringValue(val)
    } else {
        data.SubscriberPhone = types.StringNull()
    }
    if obj, ok := dataMap["subscriberWebhook"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberWebhook = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriberWebhook = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriberWebhook = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriberWebhook = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriberWebhook = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberWebhook = types.StringNull()
        }
    } else if val, ok := dataMap["subscriberWebhook"].(string); ok && val != "" {
        data.SubscriberWebhook = types.StringValue(val)
    } else {
        data.SubscriberWebhook = types.StringNull()
    }
    if obj, ok := dataMap["slackIncomingWebhookUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SlackIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SlackIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SlackIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SlackIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.SlackIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SlackIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.SlackIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SlackIncomingWebhookUrl = types.StringValue(string(jsonBytes))
        } else {
            data.SlackIncomingWebhookUrl = types.StringNull()
        }
    } else if val, ok := dataMap["slackIncomingWebhookUrl"].(string); ok && val != "" {
        data.SlackIncomingWebhookUrl = types.StringValue(val)
    } else {
        data.SlackIncomingWebhookUrl = types.StringNull()
    }
    if obj, ok := dataMap["slackWorkspaceName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SlackWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SlackWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
        } else {
            data.SlackWorkspaceName = types.StringNull()
        }
    } else if val, ok := dataMap["slackWorkspaceName"].(string); ok && val != "" {
        data.SlackWorkspaceName = types.StringValue(val)
    } else {
        data.SlackWorkspaceName = types.StringNull()
    }
    if obj, ok := dataMap["microsoftTeamsIncomingWebhookUrl"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(string(jsonBytes))
        } else {
            data.MicrosoftTeamsIncomingWebhookUrl = types.StringNull()
        }
    } else if val, ok := dataMap["microsoftTeamsIncomingWebhookUrl"].(string); ok && val != "" {
        data.MicrosoftTeamsIncomingWebhookUrl = types.StringValue(val)
    } else {
        data.MicrosoftTeamsIncomingWebhookUrl = types.StringNull()
    }
    if obj, ok := dataMap["microsoftTeamsWorkspaceName"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
            } else {
                data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
        } else {
            data.MicrosoftTeamsWorkspaceName = types.StringNull()
        }
    } else if val, ok := dataMap["microsoftTeamsWorkspaceName"].(string); ok && val != "" {
        data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
    } else {
        data.MicrosoftTeamsWorkspaceName = types.StringNull()
    }
    if val, ok := dataMap["isSubscriptionConfirmed"].(bool); ok {
        data.IsSubscriptionConfirmed = types.BoolValue(val)
    }
    if obj, ok := dataMap["subscriptionConfirmationToken"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriptionConfirmationToken = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.SubscriptionConfirmationToken = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.SubscriptionConfirmationToken = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.SubscriptionConfirmationToken = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriptionConfirmationToken = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.SubscriptionConfirmationToken = types.StringValue(string(jsonBytes))
            } else {
                data.SubscriptionConfirmationToken = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.SubscriptionConfirmationToken = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriptionConfirmationToken = types.StringNull()
        }
    } else if val, ok := dataMap["subscriptionConfirmationToken"].(string); ok && val != "" {
        data.SubscriptionConfirmationToken = types.StringValue(val)
    } else {
        data.SubscriptionConfirmationToken = types.StringNull()
    }
    if val, ok := dataMap["isUnsubscribed"].(bool); ok {
        data.IsUnsubscribed = types.BoolValue(val)
    }
    if val, ok := dataMap["sendYouHaveSubscribedMessage"].(bool); ok {
        data.SendYouHaveSubscribedMessage = types.BoolValue(val)
    }
    if val, ok := dataMap["isSubscribedToAllResources"].(bool); ok {
        data.IsSubscribedToAllResources = types.BoolValue(val)
    }
    if val, ok := dataMap["isSubscribedToAllEventTypes"].(bool); ok {
        data.IsSubscribedToAllEventTypes = types.BoolValue(val)
    }
    if val, ok := dataMap["statusPageResources"].([]interface{}); ok {
        // Convert API response list to Terraform set
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                // Handle objects with _id field (OneUptime format)
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else {
                    // Convert entire object to JSON string if no id field
                    if jsonBytes, err := json.Marshal(itemMap); err == nil {
                        setItems = append(setItems, types.StringValue(string(jsonBytes)))
                    }
                }
            } else if str, ok := item.(string); ok {
                // Handle direct string values
                setItems = append(setItems, types.StringValue(str))
            }
        }
        // Sort set items for deterministic state representation
        sort.Slice(setItems, func(i, j int) bool {
            iStr := setItems[i].(types.String).ValueString()
            jStr := setItems[j].(types.String).ValueString()
            return iStr < jStr
        })
        data.StatusPageResources = types.SetValueMust(types.StringType, setItems)
    } else {
        // For sets, always use empty set instead of null to match default values
        data.StatusPageResources = types.SetValueMust(types.StringType, []attr.Value{})
    }
    if obj, ok := dataMap["statusPageEventTypes"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageEventTypes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.StatusPageEventTypes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
            } else {
                data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageEventTypes = types.StringNull()
        }
    } else if val, ok := dataMap["statusPageEventTypes"].(string); ok && val != "" {
        data.StatusPageEventTypes = types.StringValue(val)
    } else {
        data.StatusPageEventTypes = types.StringNull()
    }
    if obj, ok := dataMap["internalNote"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InternalNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.InternalNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.InternalNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.InternalNote = types.StringValue(string(jsonBytes))
            } else {
                data.InternalNote = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.InternalNote = types.StringValue(string(jsonBytes))
            } else {
                data.InternalNote = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.InternalNote = types.StringValue(string(jsonBytes))
        } else {
            data.InternalNote = types.StringNull()
        }
    } else if val, ok := dataMap["internalNote"].(string); ok && val != "" {
        data.InternalNote = types.StringValue(val)
    } else {
        data.InternalNote = types.StringNull()
    }
    if obj, ok := dataMap["createdAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.CreatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.CreatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.CreatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.CreatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.CreatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.CreatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["createdAt"].(string); ok && val != "" {
        data.CreatedAt = types.StringValue(val)
    } else {
        data.CreatedAt = types.StringNull()
    }
    if obj, ok := dataMap["updatedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.UpdatedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.UpdatedAt = types.StringValue(string(jsonBytes))
            } else {
                data.UpdatedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.UpdatedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UpdatedAt = types.StringNull()
        }
    } else if val, ok := dataMap["updatedAt"].(string); ok && val != "" {
        data.UpdatedAt = types.StringValue(val)
    } else {
        data.UpdatedAt = types.StringNull()
    }
    if obj, ok := dataMap["deletedAt"].(map[string]interface{}); ok {
        // Handle ObjectID type responses and wrapper objects (e.g., Version, DateTime, Name types)
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            // Unwrap wrapper objects - extract the inner value regardless of whether it's empty
            data.DeletedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            // Handle numeric values that might be returned as float64
            data.DeletedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if typeStr, typeOk := obj["_type"].(string); typeOk && r.isValidOneUptimeObjectType(typeStr) && obj["value"] != nil {
            // For typed wrapper objects (only valid OneUptime ObjectTypes), preserve the full structure including _type
            normalizedObj := r.normalizeURLWrappers(obj)
            if jsonBytes, err := json.Marshal(normalizedObj); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedObj))
            }
        } else if obj["value"] != nil {
            // Handle complex value types (maps, arrays) by marshaling to JSON
            normalizedValue := r.normalizeURLWrappers(obj["value"])
            if jsonBytes, err := json.Marshal(normalizedValue); err == nil {
                data.DeletedAt = types.StringValue(string(jsonBytes))
            } else {
                data.DeletedAt = types.StringValue(fmt.Sprintf("%v", normalizedValue))
            }
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            // Fallback to JSON marshaling for other complex objects
            data.DeletedAt = types.StringValue(string(jsonBytes))
        } else {
            data.DeletedAt = types.StringNull()
        }
    } else if val, ok := dataMap["deletedAt"].(string); ok && val != "" {
        data.DeletedAt = types.StringValue(val)
    } else {
        data.DeletedAt = types.StringNull()
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
func (r *StatusPageSubscriberResource) convertTerraformSetToInterface(terraformSet types.Set) interface{} {
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
func (r *StatusPageSubscriberResource) parseJSONField(terraformString types.String) interface{} {
    if terraformString.IsNull() || terraformString.IsUnknown() || terraformString.ValueString() == "" {
        return nil
    }

    var result interface{}
    if err := json.Unmarshal([]byte(terraformString.ValueString()), &result); err != nil {
        // If JSON parsing fails, return the raw string
        return terraformString.ValueString()
    }

    return result
}

// Normalize URL wrapper objects to avoid drift (e.g., trailing slash differences).
func (r *StatusPageSubscriberResource) normalizeURLWrappers(value interface{}) interface{} {
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

func (r *StatusPageSubscriberResource) normalizeURLString(value string) string {
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
func (r *StatusPageSubscriberResource) bigFloatToFloat64(bf *big.Float) interface{} {
    if bf == nil {
        return nil
    }
    f, _ := bf.Float64()
    return f
}

// Helper method to check if a type string is a valid OneUptime ObjectType
// Only these types should be marshalled/unmarshalled as typed wrapper objects
// This list is dynamically generated from Common/Types/JSON.ts ObjectType enum
func (r *StatusPageSubscriberResource) isValidOneUptimeObjectType(typeStr string) bool {
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
