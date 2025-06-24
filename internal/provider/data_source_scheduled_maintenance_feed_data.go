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
var _ datasource.DataSource = &ScheduledMaintenanceFeedDataDataSource{}

func NewScheduledMaintenanceFeedDataDataSource() datasource.DataSource {
    return &ScheduledMaintenanceFeedDataDataSource{}
}

// ScheduledMaintenanceFeedDataDataSource defines the data source implementation.
type ScheduledMaintenanceFeedDataDataSource struct {
    client *Client
}

// ScheduledMaintenanceFeedDataDataSourceModel describes the data source data model.
type ScheduledMaintenanceFeedDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    ScheduledMaintenanceId types.String `tfsdk:"scheduled_maintenance_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    FeedInfoInMarkdown types.String `tfsdk:"feed_info_in_markdown"`
    MoreInformationInMarkdown types.String `tfsdk:"more_information_in_markdown"`
    ScheduledMaintenanceFeedEventType types.String `tfsdk:"scheduled_maintenance_feed_event_type"`
    DisplayColor types.String `tfsdk:"display_color"`
    UserId types.String `tfsdk:"user_id"`
    PostedAt types.String `tfsdk:"posted_at"`
}

func (d *ScheduledMaintenanceFeedDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_feed_data"
}

func (d *ScheduledMaintenanceFeedDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "scheduled_maintenance_feed_data data source",

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
            "scheduled_maintenance_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "feed_info_in_markdown": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance Log], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "more_information_in_markdown": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance Log], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "scheduled_maintenance_feed_event_type": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance Log], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance Log], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "display_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "posted_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceFeedDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceFeedDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceFeedDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "scheduled-maintenance-feed" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_feed_data, got error: %s", err))
        return
    }

    var scheduledMaintenanceFeedDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &scheduledMaintenanceFeedDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_feed_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := scheduledMaintenanceFeedDataResponse["data"].(map[string]interface{}); ok {
        scheduledMaintenanceFeedDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := scheduledMaintenanceFeedDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["scheduled_maintenance_id"].(string); ok {
        data.ScheduledMaintenanceId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["feed_info_in_markdown"].(string); ok {
        data.FeedInfoInMarkdown = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["more_information_in_markdown"].(string); ok {
        data.MoreInformationInMarkdown = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["scheduled_maintenance_feed_event_type"].(string); ok {
        data.ScheduledMaintenanceFeedEventType = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["display_color"].(string); ok {
        data.DisplayColor = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["user_id"].(string); ok {
        data.UserId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceFeedDataResponse["posted_at"].(string); ok {
        data.PostedAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
