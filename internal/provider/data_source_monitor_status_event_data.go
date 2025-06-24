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
var _ datasource.DataSource = &MonitorStatusEventDataDataSource{}

func NewMonitorStatusEventDataDataSource() datasource.DataSource {
    return &MonitorStatusEventDataDataSource{}
}

// MonitorStatusEventDataDataSource defines the data source implementation.
type MonitorStatusEventDataDataSource struct {
    client *Client
}

// MonitorStatusEventDataDataSourceModel describes the data source data model.
type MonitorStatusEventDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    MonitorStatusId types.String `tfsdk:"monitor_status_id"`
    IsOwnerNotified types.Bool `tfsdk:"is_owner_notified"`
    StatusChangeLog types.String `tfsdk:"status_change_log"`
    RootCause types.String `tfsdk:"root_cause"`
    EndsAt types.String `tfsdk:"ends_at"`
    StartsAt types.String `tfsdk:"starts_at"`
}

func (d *MonitorStatusEventDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_monitor_status_event_data"
}

func (d *MonitorStatusEventDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "monitor_status_event_data data source",

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
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitor_status_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "is_owner_notified": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor Status Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "status_change_log": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [No access - you don't have permission for this operation], Read: [Project Owner, Project Admin, Project Member, Read Monitor Status Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "root_cause": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Monitor Status Timeline], Read: [Project Owner, Project Admin, Project Member, Read Monitor Status Timeline], Update: [No access - you don't have permission for this operation]",
                Computed: true,
            },
            "ends_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
            "starts_at": schema.StringAttribute{
                MarkdownDescription: "A date time object.",
                Computed: true,
            },
        },
    }
}

func (d *MonitorStatusEventDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorStatusEventDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MonitorStatusEventDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "monitor-status-timeline" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor_status_event_data, got error: %s", err))
        return
    }

    var monitorStatusEventDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &monitorStatusEventDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor_status_event_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := monitorStatusEventDataResponse["data"].(map[string]interface{}); ok {
        monitorStatusEventDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := monitorStatusEventDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := monitorStatusEventDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["monitor_id"].(string); ok {
        data.MonitorId = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["monitor_status_id"].(string); ok {
        data.MonitorStatusId = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["is_owner_notified"].(bool); ok {
        data.IsOwnerNotified = types.BoolValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["status_change_log"].(string); ok {
        data.StatusChangeLog = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["root_cause"].(string); ok {
        data.RootCause = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["ends_at"].(string); ok {
        data.EndsAt = types.StringValue(val)
    }
    if val, ok := monitorStatusEventDataResponse["starts_at"].(string); ok {
        data.StartsAt = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
