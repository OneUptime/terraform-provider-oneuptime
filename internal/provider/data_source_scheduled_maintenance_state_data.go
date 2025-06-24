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
var _ datasource.DataSource = &ScheduledMaintenanceStateDataDataSource{}

func NewScheduledMaintenanceStateDataDataSource() datasource.DataSource {
    return &ScheduledMaintenanceStateDataDataSource{}
}

// ScheduledMaintenanceStateDataDataSource defines the data source implementation.
type ScheduledMaintenanceStateDataDataSource struct {
    client *Client
}

// ScheduledMaintenanceStateDataDataSourceModel describes the data source data model.
type ScheduledMaintenanceStateDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    Slug types.String `tfsdk:"slug"`
    Description types.String `tfsdk:"description"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    DeletedByUserId types.String `tfsdk:"deleted_by_user_id"`
    Color types.String `tfsdk:"color"`
    IsScheduledState types.Bool `tfsdk:"is_scheduled_state"`
    IsOngoingState types.Bool `tfsdk:"is_ongoing_state"`
    IsEndedState types.Bool `tfsdk:"is_ended_state"`
    IsResolvedState types.Bool `tfsdk:"is_resolved_state"`
    Order types.Number `tfsdk:"order"`
}

func (d *ScheduledMaintenanceStateDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance_state_data"
}

func (d *ScheduledMaintenanceStateDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "scheduled_maintenance_state_data data source",

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
            "slug": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance State], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance State], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance State], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance State]",
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
            "color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "is_scheduled_state": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance State], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance State], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance State]",
                Computed: true,
            },
            "is_ongoing_state": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance State], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance State], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance State]",
                Computed: true,
            },
            "is_ended_state": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance State], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance State], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance State]",
                Computed: true,
            },
            "is_resolved_state": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance State], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance State], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance State]",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Scheduled Maintenance State], Read: [Project Owner, Project Admin, Project Member, Read Scheduled Maintenance State], Update: [Project Owner, Project Admin, Project Member, Edit Scheduled Maintenance State]",
                Computed: true,
            },
        },
    }
}

func (d *ScheduledMaintenanceStateDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ScheduledMaintenanceStateDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ScheduledMaintenanceStateDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "scheduled-maintenance-state" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read scheduled_maintenance_state_data, got error: %s", err))
        return
    }

    var scheduledMaintenanceStateDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &scheduledMaintenanceStateDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse scheduled_maintenance_state_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := scheduledMaintenanceStateDataResponse["data"].(map[string]interface{}); ok {
        scheduledMaintenanceStateDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := scheduledMaintenanceStateDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := scheduledMaintenanceStateDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["slug"].(string); ok {
        data.Slug = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["description"].(string); ok {
        data.Description = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["deleted_by_user_id"].(string); ok {
        data.DeletedByUserId = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["color"].(string); ok {
        data.Color = types.StringValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["is_scheduled_state"].(bool); ok {
        data.IsScheduledState = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["is_ongoing_state"].(bool); ok {
        data.IsOngoingState = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["is_ended_state"].(bool); ok {
        data.IsEndedState = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["is_resolved_state"].(bool); ok {
        data.IsResolvedState = types.BoolValue(val)
    }
    if val, ok := scheduledMaintenanceStateDataResponse["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
