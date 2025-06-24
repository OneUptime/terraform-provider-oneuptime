package provider

import (
    "context"
    "fmt"
    "math/big"
    "github.com/hashicorp/terraform-plugin-framework/attr"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StatusPageSubscriberDataDataSource{}

func NewStatusPageSubscriberDataDataSource() datasource.DataSource {
    return &StatusPageSubscriberDataDataSource{}
}

// StatusPageSubscriberDataDataSource defines the data source implementation.
type StatusPageSubscriberDataDataSource struct {
    client *Client
}

// StatusPageSubscriberDataDataSourceModel describes the data source data model.
type StatusPageSubscriberDataDataSourceModel struct {
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
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    IsSubscriptionConfirmed types.Bool `tfsdk:"is_subscription_confirmed"`
    IsUnsubscribed types.Bool `tfsdk:"is_unsubscribed"`
    SendYouHaveSubscribedMessage types.Bool `tfsdk:"send_you_have_subscribed_message"`
    IsSubscribedToAllResources types.Bool `tfsdk:"is_subscribed_to_all_resources"`
    IsSubscribedToAllEventTypes types.Bool `tfsdk:"is_subscribed_to_all_event_types"`
    StatusPageResources types.List `tfsdk:"status_page_resources"`
    StatusPageEventTypes types.String `tfsdk:"status_page_event_types"`
    InternalNote types.String `tfsdk:"internal_note"`
}

func (d *StatusPageSubscriberDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_subscriber_data"
}

func (d *StatusPageSubscriberDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_subscriber_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
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
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Subscriber]",
                Computed: true,
            },
            "slack_workspace_name": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Subscriber]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_subscription_confirmed": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Subscriber]",
                Computed: true,
            },
            "is_unsubscribed": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Subscriber]",
                Computed: true,
            },
            "send_you_have_subscribed_message": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "is_subscribed_to_all_resources": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Subscriber]",
                Computed: true,
            },
            "is_subscribed_to_all_event_types": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Subscriber]",
                Computed: true,
            },
            "status_page_resources": schema.ListAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Subscriber]",
                Computed: true,
                ElementType: types.StringType,
            },
            "status_page_event_types": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Subscriber]",
                Computed: true,
            },
            "internal_note": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Subscriber, Public], Read: [Project Owner, Project Admin, Project Member, Read Status Page Subscriber], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Subscriber]",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageSubscriberDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageSubscriberDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageSubscriberDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-subscriber" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_subscriber_data, got error: %s", err))
        return
    }

    var statusPageSubscriberDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageSubscriberDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_subscriber_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageSubscriberDataResponse["data"].(map[string]interface{}); ok {
        statusPageSubscriberDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageSubscriberDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageSubscriberDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["subscriber_email"].(string); ok {
        data.SubscriberEmail = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["subscriber_phone"].(string); ok {
        data.SubscriberPhone = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["subscriber_webhook"].(string); ok {
        data.SubscriberWebhook = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["slack_workspace_name"].(string); ok {
        data.SlackWorkspaceName = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["is_subscription_confirmed"].(bool); ok {
        data.IsSubscriptionConfirmed = types.BoolValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["is_unsubscribed"].(bool); ok {
        data.IsUnsubscribed = types.BoolValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["send_you_have_subscribed_message"].(bool); ok {
        data.SendYouHaveSubscribedMessage = types.BoolValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["is_subscribed_to_all_resources"].(bool); ok {
        data.IsSubscribedToAllResources = types.BoolValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["is_subscribed_to_all_event_types"].(bool); ok {
        data.IsSubscribedToAllEventTypes = types.BoolValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["status_page_resources"].([]interface{}); ok {
        elements := make([]attr.Value, len(val))
        for i, item := range val {
            if strItem, ok := item.(string); ok {
                elements[i] = types.StringValue(strItem)
            } else {
                elements[i] = types.StringValue("")
            }
        }
        listValue, _ := types.ListValue(types.StringType, elements)
        data.StatusPageResources = listValue
    }
    if val, ok := statusPageSubscriberDataResponse["status_page_event_types"].(string); ok {
        data.StatusPageEventTypes = types.StringValue(val)
    }
    if val, ok := statusPageSubscriberDataResponse["internal_note"].(string); ok {
        data.InternalNote = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
