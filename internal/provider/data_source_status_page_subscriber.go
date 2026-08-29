package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "math/big"
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
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
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
        MarkdownDescription: "Subscriber that subscribed to your status page Look up an existing status_page_subscriber by `id` or by `name`.",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "subscriber_email": schema.StringAttribute{
                MarkdownDescription: "Email object",
                Computed: true,
            },
            "subscriber_phone": schema.StringAttribute{
                MarkdownDescription: "Phone object",
                Computed: true,
            },
            "subscriber_webhook": schema.StringAttribute{
                MarkdownDescription: "Webhook to ping when events happen on Status Page.",
                Computed: true,
            },
            "slack_workspace_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Slack workspace for validation and identification.",
                Computed: true,
            },
            "microsoft_teams_workspace_name": schema.StringAttribute{
                MarkdownDescription: "Name of the Microsoft Teams workspace for validation and identification.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_subscription_confirmed": schema.BoolAttribute{
                MarkdownDescription: "Has subscriber confirmed their subscription? (for example, by clicking on a confirmation link in an email).",
                Computed: true,
            },
            "is_unsubscribed": schema.BoolAttribute{
                MarkdownDescription: "Is Subscriber Unsubscribed?.",
                Computed: true,
            },
            "send_you_have_subscribed_message": schema.BoolAttribute{
                MarkdownDescription: "Send You Have Subscribed Message when subscriber is created?.",
                Computed: true,
            },
            "is_subscribed_to_all_resources": schema.BoolAttribute{
                MarkdownDescription: "Is Subscriber Subscribed to All Resources on this status page?.",
                Computed: true,
            },
            "is_subscribed_to_all_event_types": schema.BoolAttribute{
                MarkdownDescription: "Is Subscriber Subscribed to All Event Types (like Incidents, Scheduled Events, Announcements) on this status page?.",
                Computed: true,
            },
            "status_page_resources": schema.SetAttribute{
                MarkdownDescription: "Relation to Status Page Resources where this subscriber is subscribed to.",
                Computed: true,
                ElementType: types.StringType,
            },
            "status_page_event_types": schema.StringAttribute{
                MarkdownDescription: "Which event types is the subscriber subscribed to (like Incidents, Scheduled Events, Announcements).",
                Computed: true,
            },
            "internal_note": schema.StringAttribute{
                MarkdownDescription: "Any notes or text you would like to add to this subscriber object. This is for internal use only..",
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

    hasId := !data.Id.IsNull() && data.Id.ValueString() != ""
    hasName := !data.Name.IsNull() && data.Name.ValueString() != ""
    if hasId == hasName {
        resp.Diagnostics.AddError(
            "Invalid Lookup",
            "Exactly one of `id` or `name` must be set to look up a status_page_subscriber.",
        )
        return
    }

    selectParam := map[string]interface{}{
        "name": true,
        "createdAt": true,
        "updatedAt": true,
        "deletedAt": true,
        "version": true,
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status_page_subscriber found with id %q.", data.Id.ValueString()))
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
    } else {
        listBody := map[string]interface{}{
            "query": map[string]interface{}{
                "name": data.Name.ValueString(),
            },
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
            resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No status_page_subscriber found with name %q.", data.Name.ValueString()))
            return
        }
        if len(items) > 1 {
            resp.Diagnostics.AddError("Ambiguous Match", fmt.Sprintf("More than one status_page_subscriber matches name %q. Use the id attribute to disambiguate.", data.Name.ValueString()))
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
