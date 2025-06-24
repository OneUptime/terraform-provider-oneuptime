package provider

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MonitorLogDataDataSource{}

func NewMonitorLogDataDataSource() datasource.DataSource {
    return &MonitorLogDataDataSource{}
}

// MonitorLogDataDataSource defines the data source implementation.
type MonitorLogDataDataSource struct {
    client *Client
}

// MonitorLogDataDataSourceModel describes the data source data model.
type MonitorLogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    MonitorId types.String `tfsdk:"monitor_id"`
    Time types.String `tfsdk:"time"`
    LogBody types.String `tfsdk:"log_body"`
}

func (d *MonitorLogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_monitor_log_data"
}

func (d *MonitorLogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "monitor_log_data data source",

        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "Identifier to filter by",
                Optional: true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "Name to filter by",
                Optional: true,
            },
            "project_id": schema.StringAttribute{
                MarkdownDescription: "Project ID",
                Computed: true,
            },
            "monitor_id": schema.StringAttribute{
                MarkdownDescription: "Monitor ID",
                Computed: true,
            },
            "time": schema.StringAttribute{
                MarkdownDescription: "Time",
                Computed: true,
            },
            "log_body": schema.StringAttribute{
                MarkdownDescription: "Log Body",
                Computed: true,
            },
        },
    }
}

func (d *MonitorLogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *MonitorLogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data MonitorLogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "monitor-log" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read monitor_log_data, got error: %s", err))
        return
    }

    var monitorLogDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &monitorLogDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse monitor_log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := monitorLogDataResponse["data"].(map[string]interface{}); ok {
        monitorLogDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := monitorLogDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := monitorLogDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := monitorLogDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := monitorLogDataResponse["monitor_id"].(string); ok {
        data.MonitorId = types.StringValue(val)
    }
    if val, ok := monitorLogDataResponse["time"].(string); ok {
        data.Time = types.StringValue(val)
    }
    if val, ok := monitorLogDataResponse["log_body"].(string); ok {
        data.LogBody = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
