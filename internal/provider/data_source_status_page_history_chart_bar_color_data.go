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
var _ datasource.DataSource = &StatusPageHistoryChartBarColorDataDataSource{}

func NewStatusPageHistoryChartBarColorDataDataSource() datasource.DataSource {
    return &StatusPageHistoryChartBarColorDataDataSource{}
}

// StatusPageHistoryChartBarColorDataDataSource defines the data source implementation.
type StatusPageHistoryChartBarColorDataDataSource struct {
    client *Client
}

// StatusPageHistoryChartBarColorDataDataSourceModel describes the data source data model.
type StatusPageHistoryChartBarColorDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    CreatedAt types.String `tfsdk:"created_at"`
    UpdatedAt types.String `tfsdk:"updated_at"`
    DeletedAt types.String `tfsdk:"deleted_at"`
    Version types.Number `tfsdk:"version"`
    ProjectId types.String `tfsdk:"project_id"`
    StatusPageId types.String `tfsdk:"status_page_id"`
    UptimePercentGreaterThanOrEqualTo types.Number `tfsdk:"uptime_percent_greater_than_or_equal_to"`
    BarColor types.String `tfsdk:"bar_color"`
    CreatedByUserId types.String `tfsdk:"created_by_user_id"`
    Order types.Number `tfsdk:"order"`
}

func (d *StatusPageHistoryChartBarColorDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_status_page_history_chart_bar_color_data"
}

func (d *StatusPageHistoryChartBarColorDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "status_page_history_chart_bar_color_data data source",

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
            "uptime_percent_greater_than_or_equal_to": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page History Chart Bar Color Rule], Read: [Project Owner, Project Admin, Project Member, Read Status Page History Chart Bar Color Rule], Update: [Project Owner, Project Admin, Project Member, Edit Status Page History Chart Bar Color Rule]",
                Computed: true,
            },
            "bar_color": schema.StringAttribute{
                MarkdownDescription: "Color object",
                Computed: true,
            },
            "created_by_user_id": schema.StringAttribute{
                MarkdownDescription: "A unique identifier for an object, represented as a UUID.",
                Computed: true,
            },
            "order": schema.NumberAttribute{
                MarkdownDescription: "Permissions - Create: [Project Owner, Project Admin, Project Member, Create Status Page History Chart Bar Color Rule], Read: [Project Owner, Project Admin, Project Member, Read Status Page History Chart Bar Color Rule], Update: [Project Owner, Project Admin, Project Member, Edit Status Page History Chart Bar Color Rule]",
                Computed: true,
            },
        },
    }
}

func (d *StatusPageHistoryChartBarColorDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StatusPageHistoryChartBarColorDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data StatusPageHistoryChartBarColorDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "status-page-history-chart-bar-color-rule" + "/" + data.Id.ValueString() + "/" + "get-item"
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read status_page_history_chart_bar_color_data, got error: %s", err))
        return
    }

    var statusPageHistoryChartBarColorDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &statusPageHistoryChartBarColorDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse status_page_history_chart_bar_color_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := statusPageHistoryChartBarColorDataResponse["data"].(map[string]interface{}); ok {
        statusPageHistoryChartBarColorDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := statusPageHistoryChartBarColorDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["created_at"].(string); ok {
        data.CreatedAt = types.StringValue(val)
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["updated_at"].(string); ok {
        data.UpdatedAt = types.StringValue(val)
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["deleted_at"].(string); ok {
        data.DeletedAt = types.StringValue(val)
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["version"].(float64); ok {
        data.Version = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["status_page_id"].(string); ok {
        data.StatusPageId = types.StringValue(val)
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["uptime_percent_greater_than_or_equal_to"].(float64); ok {
        data.UptimePercentGreaterThanOrEqualTo = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["bar_color"].(string); ok {
        data.BarColor = types.StringValue(val)
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["created_by_user_id"].(string); ok {
        data.CreatedByUserId = types.StringValue(val)
    }
    if val, ok := statusPageHistoryChartBarColorDataResponse["order"].(float64); ok {
        data.Order = types.NumberValue(big.NewFloat(val))
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
