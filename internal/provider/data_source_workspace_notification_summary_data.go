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
var _ datasource.DataSource = &WorkspaceNotificationSummaryDataDataSource{}

func NewWorkspaceNotificationSummaryDataDataSource() datasource.DataSource {
    return &WorkspaceNotificationSummaryDataDataSource{}
}

// WorkspaceNotificationSummaryDataDataSource defines the data source implementation.
type WorkspaceNotificationSummaryDataDataSource struct {
    client *Client
}

// WorkspaceNotificationSummaryDataDataSourceModel describes the data source data model.
type WorkspaceNotificationSummaryDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Description types.String `tfsdk:"description"`
    WorkspaceType types.String `tfsdk:"workspace_type"`
    SummaryType types.String `tfsdk:"summary_type"`
    RecurringInterval types.String `tfsdk:"recurring_interval"`
    NumberOfDaysOfData types.Number `tfsdk:"number_of_days_of_data"`
    SendFirstReportAt types.String `tfsdk:"send_first_report_at"`
    ChannelNames types.String `tfsdk:"channel_names"`
    TeamName types.String `tfsdk:"team_name"`
    SummaryItems types.String `tfsdk:"summary_items"`
    Filters types.String `tfsdk:"filters"`
    FilterCondition types.String `tfsdk:"filter_condition"`
    NextSendAt types.String `tfsdk:"next_send_at"`
    LastSentAt types.String `tfsdk:"last_sent_at"`
    IsEnabled types.Bool `tfsdk:"is_enabled"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
}

func (d *WorkspaceNotificationSummaryDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_workspace_notification_summary_data"
}

func (d *WorkspaceNotificationSummaryDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "workspace_notification_summary_data data source",

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
            "project_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Description of the Summary Rule. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "workspace_type": schema.StringAttribute{
                MarkdownDescription: "Type of Workspace - Slack, Microsoft Teams, etc.. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "summary_type": schema.StringAttribute{
                MarkdownDescription: "Type of summary - Incident, Alert, Incident Episode, or Alert Episode. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "recurring_interval": schema.StringAttribute{
                MarkdownDescription: "How often should the summary be sent?. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "number_of_days_of_data": schema.NumberAttribute{
                MarkdownDescription: "How many days of data to include in the summary. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "send_first_report_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "channel_names": schema.StringAttribute{
                MarkdownDescription: "List of channel names to post the summary to. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "team_name": schema.StringAttribute{
                MarkdownDescription: "Microsoft Teams team name (only for Microsoft Teams). Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "summary_items": schema.StringAttribute{
                MarkdownDescription: "Checklist of items to include in the summary. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "filters": schema.StringAttribute{
                MarkdownDescription: "Filter conditions for which items to include in the summary. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "filter_condition": schema.StringAttribute{
                MarkdownDescription: "How to combine filters - Any or All. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
                Computed: true,
            },
            "next_send_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "last_sent_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "is_enabled": schema.BoolAttribute{
                MarkdownDescription: "Is this summary rule enabled?. Permissions - Create: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Create Workspace Notification Summary], Read: [Project Admin, Project Owner, Project Member, Viewer, Settings Admin, Settings Member, Settings Viewer, Read Workspace Notification Summary], Update: [Project Admin, Project Owner, Project Member, Settings Admin, Settings Member, Edit Workspace Notification Summary]",
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

func (d *WorkspaceNotificationSummaryDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceNotificationSummaryDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data WorkspaceNotificationSummaryDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "workspace-notification-summary" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read workspace_notification_summary_data, got error: %s", err))
        return
    }

    var workspaceNotificationSummaryDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &workspaceNotificationSummaryDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse workspace_notification_summary_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := workspaceNotificationSummaryDataResponse["data"].(map[string]interface{}); ok {
        workspaceNotificationSummaryDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := workspaceNotificationSummaryDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := workspaceNotificationSummaryDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["workspace_type"].(string); ok {
        data.WorkspaceType = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["summary_type"].(string); ok {
        data.SummaryType = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["recurring_interval"].(string); ok {
        data.RecurringInterval = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["number_of_days_of_data"].(float64); ok {
        data.NumberOfDaysOfData = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := workspaceNotificationSummaryDataResponse["send_first_report_at"].(string); ok {
        data.SendFirstReportAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["channel_names"].(string); ok {
        data.ChannelNames = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["team_name"].(string); ok {
        data.TeamName = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["summary_items"].(string); ok {
        data.SummaryItems = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["filters"].(string); ok {
        data.Filters = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["filter_condition"].(string); ok {
        data.FilterCondition = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["next_send_at"].(string); ok {
        data.NextSendAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["last_sent_at"].(string); ok {
        data.LastSentAt = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["is_enabled"].(bool); ok {
        data.IsEnabled = types.BoolValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := workspaceNotificationSummaryDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
