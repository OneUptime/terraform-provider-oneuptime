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
var _ datasource.DataSource = &LogDataDataSource{}

func NewLogDataDataSource() datasource.DataSource {
    return &LogDataDataSource{}
}

// LogDataDataSource defines the data source implementation.
type LogDataDataSource struct {
    client *Client
}

// LogDataDataSourceModel describes the data source data model.
type LogDataDataSourceModel struct {
    Id types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
    ProjectId types.String `tfsdk:"project_id"`
    ServiceId types.String `tfsdk:"service_id"`
    Time types.String `tfsdk:"time"`
    TimeUnixNano types.Number `tfsdk:"time_unix_nano"`
    SeverityText types.String `tfsdk:"severity_text"`
    SeverityNumber types.Number `tfsdk:"severity_number"`
    Attributes types.String `tfsdk:"attributes"`
    TraceId types.String `tfsdk:"trace_id"`
    SpanId types.String `tfsdk:"span_id"`
    Body types.String `tfsdk:"body"`
}

func (d *LogDataDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_log_data"
}

func (d *LogDataDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "log_data data source",

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
            "service_id": schema.StringAttribute{
                MarkdownDescription: "Service ID",
                Computed: true,
            },
            "time": schema.StringAttribute{
                MarkdownDescription: "Time",
                Computed: true,
            },
            "time_unix_nano": schema.NumberAttribute{
                MarkdownDescription: "Time (in Unix Nano)",
                Computed: true,
            },
            "severity_text": schema.StringAttribute{
                MarkdownDescription: "Severity Text",
                Computed: true,
            },
            "severity_number": schema.NumberAttribute{
                MarkdownDescription: "Severity Number",
                Computed: true,
            },
            "attributes": schema.StringAttribute{
                MarkdownDescription: "Attributes",
                Computed: true,
            },
            "trace_id": schema.StringAttribute{
                MarkdownDescription: "Trace ID",
                Computed: true,
            },
            "span_id": schema.StringAttribute{
                MarkdownDescription: "Span ID",
                Computed: true,
            },
            "body": schema.StringAttribute{
                MarkdownDescription: "Log Body",
                Computed: true,
            },
        },
    }
}

func (d *LogDataDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LogDataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data LogDataDataSourceModel

    // Read Terraform configuration data into the model
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    
    // Build API path
    apiPath := "/" + "logs" + "/" + data.Id.ValueString()
    
    // Prepare request body with select fields (if needed)
    requestBody := map[string]interface{}{
        "select": map[string]interface{}{}, // Add specific fields to select if needed
    }
    
    // Make API call
    httpResp, err := d.client.Post(apiPath, requestBody)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read log_data, got error: %s", err))
        return
    }

    var logDataResponse map[string]interface{}
    err = d.client.ParseResponse(httpResp, &logDataResponse)
    if err != nil {
        resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse log_data response, got error: %s", err))
        return
    }

    // Extract data from response
    if dataMap, ok := logDataResponse["data"].(map[string]interface{}); ok {
        logDataResponse = dataMap
    }

    // Update the model with response data
    if val, ok := logDataResponse["id"].(string); ok {
        data.Id = types.StringValue(val)
    }
    if val, ok := logDataResponse["name"].(string); ok {
        data.Name = types.StringValue(val)
    }
    if val, ok := logDataResponse["project_id"].(string); ok {
        data.ProjectId = types.StringValue(val)
    }
    if val, ok := logDataResponse["service_id"].(string); ok {
        data.ServiceId = types.StringValue(val)
    }
    if val, ok := logDataResponse["time"].(string); ok {
        data.Time = types.StringValue(val)
    }
    if val, ok := logDataResponse["time_unix_nano"].(float64); ok {
        data.TimeUnixNano = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := logDataResponse["severity_text"].(string); ok {
        data.SeverityText = types.StringValue(val)
    }
    if val, ok := logDataResponse["severity_number"].(float64); ok {
        data.SeverityNumber = types.NumberValue(big.NewFloat(val))
    }
    if val, ok := logDataResponse["attributes"].(string); ok {
        data.Attributes = types.StringValue(val)
    }
    if val, ok := logDataResponse["trace_id"].(string); ok {
        data.TraceId = types.StringValue(val)
    }
    if val, ok := logDataResponse["span_id"].(string); ok {
        data.SpanId = types.StringValue(val)
    }
    if val, ok := logDataResponse["body"].(string); ok {
        data.Body = types.StringValue(val)
    }

    // Write logs using the tflog package
    tflog.Trace(ctx, "read a data source")

    // Save data into Terraform state
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
