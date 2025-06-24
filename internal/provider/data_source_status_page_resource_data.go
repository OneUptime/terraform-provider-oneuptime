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
var _ datasource.DataSource = &StatusPageResourceDataDataSource{}

func NewStatusPageResourceDataDataSource() datasource.DataSource {
    return &StatusPageResourceDataDataSource{}
}

// StatusPageResourceDataDataSource defines the data source implementation.
type StatusPageResourceDataDataSource struct {
    client *Client
}

// StatusPageResourceDataDataSourceModel describes the data source data model.
type StatusPageResourceDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    MonitorGroupId types.String `tfsdk:"monitor_group_id"`
    StatusPageGroupId types.String `tfsdk:"status_page_group_id"`
    DisplayName types.String `tfsdk:"display_name"`
    DisplayDescription types.String `tfsdk:"display_description"`
    DisplayTooltip types.String `tfsdk:"display_tooltip"`
    ShowCurrentStatus types.Bool `tfsdk:"show_current_status"`
    ShowUptimePercent types.Bool `tfsdk:"show_uptime_percent"`
    UptimePercentPrecision types.String `tfsdk:"uptime_percent_precision"`
    ShowStatusHistoryChart types.Bool `tfsdk:"show_status_history_chart"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Order types.Number `tfsdk:"order"`
}

func (d *StatusPageResourceDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_resource_data"
}

func (d *StatusPageResourceDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_resource_data data source",

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
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "monitor_group_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "status_page_group_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "display_name": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Resource], Read: [Project Owner, Project Admin, Project Member, Read Status Page Resource], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Resource]",
                Computed: true,
            },
            "display_description": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Resource], Read: [Project Owner, Project Admin, Project Member, Read Status Page Resource], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Resource]",
                Computed: true,
            },
            "display_tooltip": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Resource], Read: [Project Owner, Project Admin, Project Member, Read Status Page Resource], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Resource]",
                Computed: true,
            },
            "show_current_status": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Resource], Read: [Project Owner, Project Admin, Project Member, Read Status Page Resource], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Resource]",
                Computed: true,
            },
            "show_uptime_percent": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Resource], Read: [Project Owner, Project Admin, Project Member, Read Status Page Resource], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Resource]",
                Computed: true,
            },
            "uptime_percent_precision": schema.StringAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Resource], Read: [Project Owner, Project Admin, Project Member, Read Status Page Resource], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Resource]",
                Computed: true,
            },
            "show_status_history_chart": schema.BoolAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Resource], Read: [Project Owner, Project Admin, Project Member, Read Status Page Resource], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Resource]",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page Resource], Read: [Project Owner, Project Admin, Project Member, Read Status Page Resource], Update: [Project Owner, Project Admin, Project Member, Edit Status Page Resource]",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageResourceDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageResourceDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageResourceDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-resource" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_resource_data, got error: %s", err))
        return
    }

    var statusPageResourceDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageResourceDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_resource_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageResourceDataResponse["data"].(map[string]interface{}); ok {
        statusPageResourceDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageResourceDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageResourceDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["monitor_id"].(string); ok {
        data.MonitorId = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["monitor_group_id"].(string); ok {
        data.MonitorGroupId = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["status_page_group_id"].(string); ok {
        data.StatusPageGroupId = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["display_name"].(string); ok {
        data.DisplayName = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["display_description"].(string); ok {
        data.DisplayDescription = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["display_tooltip"].(string); ok {
        data.DisplayTooltip = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["show_current_status"].(bool); ok {
        data.ShowCurrentStatus = types.BoolValue(val)
    }
    if val, ok := statusPageResourceDataResponse["show_uptime_percent"].(bool); ok {
        data.ShowUptimePercent = types.BoolValue(val)
    }
    if val, ok := statusPageResourceDataResponse["uptime_percent_precision"].(string); ok {
        data.UptimePercentPrecision = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["show_status_history_chart"].(bool); ok {
        data.ShowStatusHistoryChart = types.BoolValue(val)
    }
    if val, ok := statusPageResourceDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageResourceDataResponse["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
