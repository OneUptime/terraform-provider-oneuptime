package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "sort"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StatusPageSubscriberDataSource{}

func NewStatusPageSubscriberDataSource() datasource.DataSource {
    return &StatusPageSubscriberDataSource{}
}

// StatusPageSubscriberDataSource defines the data source implementation.
type StatusPageSubscriberDataSource struct {
    client *Client
}

// StatusPageSubscriberDataSourceModel describes the data source data model.
type StatusPageSubscriberDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    SubscriberEmail types.String `tfsdk:"subscriber_email"`
    SubscriberPhone types.String `tfsdk:"subscriber_phone"`
    SubscriberWebhook types.String `tfsdk:"subscriber_webhook"`
    SlackWorkspaceName types.String `tfsdk:"slack_workspace_name"`
    MicrosoftTeamsWorkspaceName types.String `tfsdk:"microsoft_teams_workspace_name"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsSubscriptionConfirmed types.Bool `tfsdk:"is_subscription_confirmed"`
    IsUnsubscribed types.Bool `tfsdk:"is_unsubscribed"`
    UnsubscribedAt types.String `tfsdk:"unsubscribed_at"`
    IsAddedByTeam types.Bool `tfsdk:"is_added_by_team"`
    SendYouHaveSubscribedMessage types.Bool `tfsdk:"send_you_have_subscribed_message"`
    IsSubscribedToAllResources types.Bool `tfsdk:"is_subscribed_to_all_resources"`
    IsSubscribedToAllEventTypes types.Bool `tfsdk:"is_subscribed_to_all_event_types"`
    StatusPageResources types.Set `tfsdk:"status_page_resources"`
    StatusPageEventTypes types.String `tfsdk:"status_page_event_types"`
    InternalNote types.String `tfsdk:"internal_note"`
}

func (d *StatusPageSubscriberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_subscriber"
}

func (d *StatusPageSubscriberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Subscriber that subscribed to your status page Look up an existing status page subscriber by `id`, or by any of its other arguments (`created_by_user_id`, `internal_note`, `is_added_by_team`, ...): each one set must match, and exactly one status page subscriber may match them all.",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.",
                Computed: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.",
                Optional: true,
                Computed: true,
            },
            "subscriber_email": schema.StringAttribute{
                MarkdownDescription: "Email address of the subscriber.",
                Computed: true,
            },
            "subscriber_phone": schema.StringAttribute{
                MarkdownDescription: "Phone number of subscriber.",
                Computed: true,
            },
            "subscriber_webhook": schema.StringAttribute{
                MarkdownDescription: "Webhook to ping when events happen on Status Page.",
                Optional: true,
                Computed: true,
            },
            "slack_workspace_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Slack workspace for validation and identification.",
                Optional: true,
                Computed: true,
            },
            "microsoft_teams_workspace_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Microsoft Teams workspace for validation and identification.",
                Optional: true,
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).",
                Optional: true,
                Computed: true,
            },
            "is_subscription_confirmed": schema.BoolAttribute{
                MarkdownDescription: "Has subscriber confirmed their subscription? (for example, by clicking on a confirmation link in an email).",
                Optional: true,
                Computed: true,
            },
            "is_unsubscribed": schema.BoolAttribute{
                MarkdownDescription: "Is Subscriber Unsubscribed?",
                Optional: true,
                Computed: true,
            },
            "unsubscribed_at": schema.StringAttribute{
                MarkdownDescription: "When this subscriber unsubscribed. Set by OneUptime when Is Unsubscribed is turned on, and cleared when it is turned off; any value sent for it is ignored.",
                Computed: true,
            },
            "is_added_by_team": schema.BoolAttribute{
                MarkdownDescription: "Whether your team added this subscriber (from the dashboard, with an API key or by a workflow) rather than the subscriber signing up on the status page. Set by OneUptime when the subscriber is created; any value sent for it is ignored.",
                Optional: true,
                Computed: true,
            },
            "send_you_have_subscribed_message": schema.BoolAttribute{
                MarkdownDescription: "Send You Have Subscribed Message when subscriber is created?",
                Optional: true,
                Computed: true,
            },
            "is_subscribed_to_all_resources": schema.BoolAttribute{
                MarkdownDescription: "Is Subscriber Subscribed to All Resources on this status page?",
                Optional: true,
                Computed: true,
            },
            "is_subscribed_to_all_event_types": schema.BoolAttribute{
                MarkdownDescription: "Is Subscriber Subscribed to All Event Types (like Incidents, Scheduled Events, Announcements) on this status page?",
                Optional: true,
                Computed: true,
            },
            "status_page_resources": schema.SetAttribute{
                MarkdownDescription: "Relation to Status Page Resources where this subscriber is subscribed to. IDs of `oneuptime_status_page_resource` resources.",
                Computed: true,
                ElementType: types.StringType,
            },
            "status_page_event_types": schema.StringAttribute{
                MarkdownDescription: "Which event types is the subscriber subscribed to (like Incidents, Scheduled Events, Announcements). A JSON value: write it with `jsonencode()`.",
                Computed: true,
            },
            "internal_note": schema.StringAttribute{
                MarkdownDescription: "Any notes or text you would like to add to this subscriber object. This is for internal use only.",
                Optional: true,
                Computed: true,
            },
        },
    }
}

func (d *StatusPageSubscriberDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageSubscriberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageSubscriberDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    hasId := !data.Id.IsNull() && !data.Id.IsUnknown() && data.Id.ValueString() != ""

    // Every other argument set in configuration narrows the lookup.
    filters := map[string]interface{}{}
    filterNames := []string{}
    if !data.StatusPageId.IsNull() && !data.StatusPageId.IsUnknown() {
        filters["statusPageId"] = data.StatusPageId.ValueString()
        filterNames = append(filterNames, "status_page_id = "+fmt.Sprintf("%q", data.StatusPageId.ValueString()))
    }
    if !data.SubscriberWebhook.IsNull() && !data.SubscriberWebhook.IsUnknown() {
        filters["subscriberWebhook"] = data.SubscriberWebhook.ValueString()
        filterNames = append(filterNames, "subscriber_webhook = "+fmt.Sprintf("%q", data.SubscriberWebhook.ValueString()))
    }
    if !data.SlackWorkspaceName.IsNull() && !data.SlackWorkspaceName.IsUnknown() {
        filters["slackWorkspaceName"] = data.SlackWorkspaceName.ValueString()
        filterNames = append(filterNames, "slack_workspace_name = "+fmt.Sprintf("%q", data.SlackWorkspaceName.ValueString()))
    }
    if !data.MicrosoftTeamsWorkspaceName.IsNull() && !data.MicrosoftTeamsWorkspaceName.IsUnknown() {
        filters["microsoftTeamsWorkspaceName"] = data.MicrosoftTeamsWorkspaceName.ValueString()
        filterNames = append(filterNames, "microsoft_teams_workspace_name = "+fmt.Sprintf("%q", data.MicrosoftTeamsWorkspaceName.ValueString()))
    }
    if !data.CreatedByUserId.IsNull() && !data.CreatedByUserId.IsUnknown() {
        filters["createdByUserId"] = data.CreatedByUserId.ValueString()
        filterNames = append(filterNames, "created_by_user_id = "+fmt.Sprintf("%q", data.CreatedByUserId.ValueString()))
    }
    if !data.IsSubscriptionConfirmed.IsNull() && !data.IsSubscriptionConfirmed.IsUnknown() {
        filters["isSubscriptionConfirmed"] = data.IsSubscriptionConfirmed.ValueBool()
        filterNames = append(filterNames, "is_subscription_confirmed = "+fmt.Sprintf("%t", data.IsSubscriptionConfirmed.ValueBool()))
    }
    if !data.IsUnsubscribed.IsNull() && !data.IsUnsubscribed.IsUnknown() {
        filters["isUnsubscribed"] = data.IsUnsubscribed.ValueBool()
        filterNames = append(filterNames, "is_unsubscribed = "+fmt.Sprintf("%t", data.IsUnsubscribed.ValueBool()))
    }
    if !data.IsAddedByTeam.IsNull() && !data.IsAddedByTeam.IsUnknown() {
        filters["isAddedByTeam"] = data.IsAddedByTeam.ValueBool()
        filterNames = append(filterNames, "is_added_by_team = "+fmt.Sprintf("%t", data.IsAddedByTeam.ValueBool()))
    }
    if !data.SendYouHaveSubscribedMessage.IsNull() && !data.SendYouHaveSubscribedMessage.IsUnknown() {
        filters["sendYouHaveSubscribedMessage"] = data.SendYouHaveSubscribedMessage.ValueBool()
        filterNames = append(filterNames, "send_you_have_subscribed_message = "+fmt.Sprintf("%t", data.SendYouHaveSubscribedMessage.ValueBool()))
    }
    if !data.IsSubscribedToAllResources.IsNull() && !data.IsSubscribedToAllResources.IsUnknown() {
        filters["isSubscribedToAllResources"] = data.IsSubscribedToAllResources.ValueBool()
        filterNames = append(filterNames, "is_subscribed_to_all_resources = "+fmt.Sprintf("%t", data.IsSubscribedToAllResources.ValueBool()))
    }
    if !data.IsSubscribedToAllEventTypes.IsNull() && !data.IsSubscribedToAllEventTypes.IsUnknown() {
        filters["isSubscribedToAllEventTypes"] = data.IsSubscribedToAllEventTypes.ValueBool()
        filterNames = append(filterNames, "is_subscribed_to_all_event_types = "+fmt.Sprintf("%t", data.IsSubscribedToAllEventTypes.ValueBool()))
    }
    if !data.InternalNote.IsNull() && !data.InternalNote.IsUnknown() {
        filters["internalNote"] = data.InternalNote.ValueString()
        filterNames = append(filterNames, "internal_note = "+fmt.Sprintf("%q", data.InternalNote.ValueString()))
    }

    if hasId && len(filters) > 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Look the status page subscriber up either by `id` or by its other arguments, not both.",
        )
        return
    }
    if !hasId && len(filters) == 0 {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Set `id`, or at least one other argument to look the status page subscriber up by.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "createdAt": true,
        "updatedAt": true,
        "projectId": true,
        "statusPageId": true,
        "subscriberEmail": true,
        "subscriberPhone": true,
        "subscriberWebhook": true,
        "slackWorkspaceName": true,
        "microsoftTeamsWorkspaceName": true,
        "createdByUserId": true,
        "isSubscriptionConfirmed": true,
        "isUnsubscribed": true,
        "unsubscribedAt": true,
        "isAddedByTeam": true,
        "sendYouHaveSubscribedMessage": true,
        "isSubscribedToAllResources": true,
        "isSubscribedToAllEventTypes": true,
        "statusPageResources": true,
        "statusPageEventTypes": true,
        "internalNote": true,
        "_id": true,
    }

    var item map[string]interface{}
    if hasId {
        readPath := "/status-page-subscriber/" + data.Id.ValueString() + "/get-item"
        httpResp, err := d.client.PostWithSelect(ctx, readPath, selectParam)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_subscriber, got error: %s", err))
            return
        }
        if httpResp.StatusCode == http.StatusNotFound {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page subscriber found with id %q.", data.Id.ValueString()))
            return
        }
        var itemResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &itemResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to read status_page_subscriber: %s", err))
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
        httpResp, err := d.client.PostBodyWithSelect(ctx, "/status-page-subscriber/get-list", listBody)
        if err != nil {
            resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list status_page_subscriber, got error: %s", err))
            return
        }
        var listResponse map[string]interface{}
        if err := d.client.ParseResponse(httpResp, &listResponse); err != nil {
            resp.Diagnostics.AddError("OneUptime API Error", fmt.Sprintf("Unable to list status_page_subscriber: %s", err))
            return
        }
        items, _ := listResponse["data"].([]interface{})
        if len(items) == 0 {
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status page subscriber matches %s.", describeLookup(filterNames)))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status page subscriber matches %s. Set more arguments to narrow the lookup down to one, or look it up by id.", describeLookup(filterNames)))
            return
        }
        first, ok := items[0].(map[string]interface{})
        if !ok {
            resp.Diagnostics.AddError("OneUptime API Error", "Unexpected list response shape for status_page_subscriber.")
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
    if obj, ok := item["projectId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.ProjectId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.ProjectId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.ProjectId = types.StringValue(string(jsonBytes))
        } else {
            data.ProjectId = types.StringNull()
        }
    } else if val, ok := item["projectId"].(string); ok {
        data.ProjectId = types.StringValue(val)
    } else {
        data.ProjectId = types.StringNull()
    }
    if obj, ok := item["statusPageId"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageId = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageId = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageId = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageId = types.StringNull()
        }
    } else if val, ok := item["statusPageId"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    } else {
        data.StatusPageId = types.StringNull()
    }
    if obj, ok := item["subscriberEmail"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberEmail = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberEmail = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberEmail = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberEmail = types.StringNull()
        }
    } else if val, ok := item["subscriberEmail"].(string); ok {
        data.SubscriberEmail = types.StringValue(val)
    } else {
        data.SubscriberEmail = types.StringNull()
    }
    if obj, ok := item["subscriberPhone"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberPhone = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberPhone = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberPhone = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberPhone = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberPhone = types.StringNull()
        }
    } else if val, ok := item["subscriberPhone"].(string); ok {
        data.SubscriberPhone = types.StringValue(val)
    } else {
        data.SubscriberPhone = types.StringNull()
    }
    if obj, ok := item["subscriberWebhook"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SubscriberWebhook = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SubscriberWebhook = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SubscriberWebhook = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SubscriberWebhook = types.StringValue(string(jsonBytes))
        } else {
            data.SubscriberWebhook = types.StringNull()
        }
    } else if val, ok := item["subscriberWebhook"].(string); ok {
        data.SubscriberWebhook = types.StringValue(val)
    } else {
        data.SubscriberWebhook = types.StringNull()
    }
    if obj, ok := item["slackWorkspaceName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.SlackWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.SlackWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.SlackWorkspaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.SlackWorkspaceName = types.StringValue(string(jsonBytes))
        } else {
            data.SlackWorkspaceName = types.StringNull()
        }
    } else if val, ok := item["slackWorkspaceName"].(string); ok {
        data.SlackWorkspaceName = types.StringValue(val)
    } else {
        data.SlackWorkspaceName = types.StringNull()
    }
    if obj, ok := item["microsoftTeamsWorkspaceName"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.MicrosoftTeamsWorkspaceName = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.MicrosoftTeamsWorkspaceName = types.StringValue(string(jsonBytes))
        } else {
            data.MicrosoftTeamsWorkspaceName = types.StringNull()
        }
    } else if val, ok := item["microsoftTeamsWorkspaceName"].(string); ok {
        data.MicrosoftTeamsWorkspaceName = types.StringValue(val)
    } else {
        data.MicrosoftTeamsWorkspaceName = types.StringNull()
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
    if val, ok := item["isSubscriptionConfirmed"].(bool); ok {
        data.IsSubscriptionConfirmed = types.BoolValue(val)
    } else {
        data.IsSubscriptionConfirmed = types.BoolNull()
    }
    if val, ok := item["isUnsubscribed"].(bool); ok {
        data.IsUnsubscribed = types.BoolValue(val)
    } else {
        data.IsUnsubscribed = types.BoolNull()
    }
    if obj, ok := item["unsubscribedAt"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.UnsubscribedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.UnsubscribedAt = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.UnsubscribedAt = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.UnsubscribedAt = types.StringValue(string(jsonBytes))
        } else {
            data.UnsubscribedAt = types.StringNull()
        }
    } else if val, ok := item["unsubscribedAt"].(string); ok {
        data.UnsubscribedAt = types.StringValue(val)
    } else {
        data.UnsubscribedAt = types.StringNull()
    }
    if val, ok := item["isAddedByTeam"].(bool); ok {
        data.IsAddedByTeam = types.BoolValue(val)
    } else {
        data.IsAddedByTeam = types.BoolNull()
    }
    if val, ok := item["sendYouHaveSubscribedMessage"].(bool); ok {
        data.SendYouHaveSubscribedMessage = types.BoolValue(val)
    } else {
        data.SendYouHaveSubscribedMessage = types.BoolNull()
    }
    if val, ok := item["isSubscribedToAllResources"].(bool); ok {
        data.IsSubscribedToAllResources = types.BoolValue(val)
    } else {
        data.IsSubscribedToAllResources = types.BoolNull()
    }
    if val, ok := item["isSubscribedToAllEventTypes"].(bool); ok {
        data.IsSubscribedToAllEventTypes = types.BoolValue(val)
    } else {
        data.IsSubscribedToAllEventTypes = types.BoolNull()
    }
    if val, ok := item["statusPageResources"].([]interface{}); ok {
        var setItems []attr.Value
        for _, item := range val {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if id, ok := itemMap["_id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if id, ok := itemMap["id"].(string); ok {
                    setItems = append(setItems, types.StringValue(id))
                } else if jsonBytes, err := json.Marshal(itemMap); err == nil {
                    setItems = append(setItems, types.StringValue(string(jsonBytes)))
                }
            } else if str, ok := item.(string); ok {
                setItems = append(setItems, types.StringValue(str))
            } else {
                setItems = append(setItems, types.StringValue(fmt.Sprintf("%v", item)))
            }
        }
        sort.Slice(setItems, func(i, j int) bool {
            return setItems[i].(types.String).ValueString() < setItems[j].(types.String).ValueString()
        })
        data.StatusPageResources = types.SetValueMust(types.StringType, setItems)
    } else {
        data.StatusPageResources = types.SetNull(types.StringType)
    }
    if obj, ok := item["statusPageEventTypes"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.StatusPageEventTypes = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.StatusPageEventTypes = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.StatusPageEventTypes = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.StatusPageEventTypes = types.StringValue(string(jsonBytes))
        } else {
            data.StatusPageEventTypes = types.StringNull()
        }
    } else if val, ok := item["statusPageEventTypes"].(string); ok {
        data.StatusPageEventTypes = types.StringValue(val)
    } else {
        data.StatusPageEventTypes = types.StringNull()
    }
    if obj, ok := item["internalNote"].(map[string]interface{}); ok {
        if val, ok := obj["_id"].(string); ok && val != "" {
            data.InternalNote = types.StringValue(val)
        } else if val, ok := obj["value"].(string); ok {
            data.InternalNote = types.StringValue(val)
        } else if val, ok := obj["value"].(float64); ok {
            data.InternalNote = types.StringValue(fmt.Sprintf("%v", val))
        } else if jsonBytes, err := json.Marshal(obj); err == nil {
            data.InternalNote = types.StringValue(string(jsonBytes))
        } else {
            data.InternalNote = types.StringNull()
        }
    } else if val, ok := item["internalNote"].(string); ok {
        data.InternalNote = types.StringValue(val)
    } else {
        data.InternalNote = types.StringNull()
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
