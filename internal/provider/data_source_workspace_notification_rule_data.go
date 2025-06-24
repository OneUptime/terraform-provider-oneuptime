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
var _ datasource.DataSource = &WorkspaceNotificationRuleDataDataSource{}

func NewWorkspaceNotificationRuleDataDataSource() datasource.DataSource {
    return &WorkspaceNotificationRuleDataDataSource{}
}

// WorkspaceNotificationRuleDataDataSource defines the data source implementation.
type WorkspaceNotificationRuleDataDataSource struct {
    client *Client
}

// WorkspaceNotificationRuleDataDataSourceModel describes the data source data model.
type WorkspaceNotificationRuleDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    NotificationRule types.String `tfsdk:"notification_rule"`
    EventType types.String `tfsdk:"event_type"`
    WorkspaceType types.String `tfsdk:"workspace_type"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *WorkspaceNotificationRuleDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workspace_notification_rule_data"
}

func (d *WorkspaceNotificationRuleDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "workspace_notification_rule_data data source",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Admin, Project Owner, Project Member, Create Workspace Notification Rule], Read: [Project Admin, Project Owner, Project Member, Read Workspace Notification Rule], Update: [Project Admin, Project Owner, Project Member, Edit Workspace Notification Rule]",
                Computed: true,
            },
            "notification_rule": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Admin, Project Owner, Project Member, Create Workspace Notification Rule], Read: [Project Admin, Project Owner, Project Member, Read Workspace Notification Rule], Update: [Project Admin, Project Owner, Project Member, Edit Workspace Notification Rule]",
                Computed: true,
            },
            "event_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Admin, Project Owner, Project Member, Create Workspace Notification Rule], Read: [Project Admin, Project Owner, Project Member, Read Workspace Notification Rule], Update: [Project Admin, Project Owner, Project Member, Edit Workspace Notification Rule]",
                Computed: true,
            },
            "workspace_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Admin, Project Owner, Project Member, Create Workspace Notification Rule], Read: [Project Admin, Project Owner, Project Member, Read Workspace Notification Rule], Update: [Project Admin, Project Owner, Project Member, Edit Workspace Notification Rule]",
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
        },
    }
}

func (d *WorkspaceNotificationRuleDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceNotificationRuleDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkspaceNotificationRuleDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "workspace-notification-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace_notification_rule_data, got error: %s", err))
        return
    }

    var workspaceNotificationRuleDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &workspaceNotificationRuleDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse workspace_notification_rule_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := workspaceNotificationRuleDataResponse["data"].(map[string]interface{}); ok {
        workspaceNotificationRuleDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := workspaceNotificationRuleDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := workspaceNotificationRuleDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["notification_rule"].(string); ok {
        data.NotificationRule = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["event_type"].(string); ok {
        data.EventType = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["workspace_type"].(string); ok {
        data.WorkspaceType = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationRuleDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
